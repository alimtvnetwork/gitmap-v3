package movemerge

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/gitmap/constants"
)

// CopyOptions controls file-tree copy behaviour.
type CopyOptions struct {
	IncludeVCS         bool // include .git/
	IncludeNodeModules bool // include node_modules/
	DryRun             bool
}

// CopyResult summarises a copy operation.
type CopyResult struct {
	FilesCopied int
	BytesCopied int64
}

// CopyTree copies every entry under srcDir into dstDir, honouring
// the default ignore list (.git, node_modules, .gitmap/release-assets/).
// File modes and symlinks are preserved.
func CopyTree(srcDir, dstDir string, opts CopyOptions) (CopyResult, error) {
	var result CopyResult
	walkErr := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(srcDir, path)
		if relErr != nil {
			return relErr
		}
		if rel == "." {
			return nil
		}
		if shouldIgnore(rel, info, opts) {
			if info.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		return copyEntry(path, filepath.Join(dstDir, rel), info, opts, &result)
	})

	return result, walkErr
}

// shouldIgnore returns true when rel matches the default ignore list.
func shouldIgnore(rel string, info os.FileInfo, opts CopyOptions) bool {
	base := filepath.Base(rel)
	if !opts.IncludeVCS && base == ".git" {
		return true
	}
	if !opts.IncludeNodeModules && base == "node_modules" {
		return true
	}
	relSlash := filepath.ToSlash(rel)
	for _, prefix := range constants.DefaultMoveMergeIgnorePrefixes {
		if strings.HasPrefix(relSlash, prefix) || relSlash == strings.TrimSuffix(prefix, "/") {
			return true
		}
	}

	return false
}

// copyEntry dispatches to the right copier for dir/symlink/regular.
func copyEntry(src, dst string, info os.FileInfo, opts CopyOptions, result *CopyResult) error {
	if opts.DryRun {
		result.FilesCopied++

		return nil
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return copySymlink(src, dst)
	}
	if info.IsDir() {
		return os.MkdirAll(dst, info.Mode().Perm())
	}

	n, err := copyFileWithMode(src, dst, info.Mode().Perm())
	if err != nil {
		return err
	}
	result.FilesCopied++
	result.BytesCopied += n

	return nil
}

// copySymlink recreates a symlink at dst pointing at the same target.
func copySymlink(src, dst string) error {
	target, err := os.Readlink(src)
	if err != nil {
		return fmt.Errorf("readlink %s: %w", src, err)
	}
	_ = os.Remove(dst)

	return os.Symlink(target, dst)
}

// copyFileWithMode copies src to dst, creating parent dirs and
// preserving mode bits. Returns bytes written.
func copyFileWithMode(src, dst string, mode os.FileMode) (int64, error) {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return 0, err
	}
	in, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	n, copyErr := io.Copy(out, in)
	if copyErr != nil {
		return n, copyErr
	}

	return n, out.Sync()
}
