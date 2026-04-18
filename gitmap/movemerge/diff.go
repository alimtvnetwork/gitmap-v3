package movemerge

import (
	"crypto/sha256"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// DiffEntryKind classifies how a single relative path differs.
type DiffEntryKind int

const (
	// DiffMissingRight = present on LEFT, missing on RIGHT.
	DiffMissingRight DiffEntryKind = iota
	// DiffMissingLeft = present on RIGHT, missing on LEFT.
	DiffMissingLeft
	// DiffIdentical = present on both, byte-equal.
	DiffIdentical
	// DiffConflict = present on both, byte-different.
	DiffConflict
)

// DiffEntry is one path classified for merge.
type DiffEntry struct {
	RelPath   string
	Kind      DiffEntryKind
	LeftSize  int64
	RightSize int64
	LeftMTime int64 // Unix seconds
	RightMTime int64
}

// DiffTrees walks both sides and classifies every relative path.
// Ignored entries (per CopyOptions defaults) are excluded.
func DiffTrees(leftDir, rightDir string, opts CopyOptions) ([]DiffEntry, error) {
	leftSet, err := indexTree(leftDir, opts)
	if err != nil {
		return nil, err
	}
	rightSet, err := indexTree(rightDir, opts)
	if err != nil {
		return nil, err
	}

	return classifyAll(leftSet, rightSet, leftDir, rightDir), nil
}

// indexTree returns rel-path -> os.FileInfo for every non-ignored file.
func indexTree(root string, opts CopyOptions) (map[string]os.FileInfo, error) {
	out := make(map[string]os.FileInfo)
	walkErr := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil || rel == "." {
			return relErr
		}
		if shouldIgnore(rel, info, opts) {
			if info.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}
		if info.IsDir() {
			return nil
		}
		out[filepath.ToSlash(rel)] = info

		return nil
	})

	return out, walkErr
}

// classifyAll merges the two indexes into a sorted DiffEntry list.
func classifyAll(left, right map[string]os.FileInfo, leftDir, rightDir string) []DiffEntry {
	seen := make(map[string]struct{}, len(left)+len(right))
	for k := range left {
		seen[k] = struct{}{}
	}
	for k := range right {
		seen[k] = struct{}{}
	}
	out := make([]DiffEntry, 0, len(seen))
	for rel := range seen {
		out = append(out, classifyOne(rel, left[rel], right[rel], leftDir, rightDir))
	}
	// Stable order = sorted ascending by rel path.
	sortByRelPath(out)

	return out
}

// classifyOne produces a single DiffEntry.
func classifyOne(rel string, l, r os.FileInfo, leftDir, rightDir string) DiffEntry {
	entry := DiffEntry{RelPath: rel}
	if l != nil {
		entry.LeftSize, entry.LeftMTime = l.Size(), l.ModTime().Unix()
	}
	if r != nil {
		entry.RightSize, entry.RightMTime = r.Size(), r.ModTime().Unix()
	}
	switch {
	case l != nil && r == nil:
		entry.Kind = DiffMissingRight
	case l == nil && r != nil:
		entry.Kind = DiffMissingLeft
	case sameContent(filepath.Join(leftDir, rel), filepath.Join(rightDir, rel)):
		entry.Kind = DiffIdentical
	default:
		entry.Kind = DiffConflict
	}

	return entry
}

// sameContent returns true when both files exist and have equal SHA-256.
func sameContent(a, b string) bool {
	ha, errA := hashFile(a)
	hb, errB := hashFile(b)
	if errA != nil || errB != nil {
		return false
	}

	return ha == hb
}

// hashFile streams a file through SHA-256.
func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", err
	}

	return string(h.Sum(nil)), nil
}

// sortByRelPath sorts the entries in place by RelPath ascending.
func sortByRelPath(entries []DiffEntry) {
	// Simple insertion sort — entry counts are small (single trees).
	for i := 1; i < len(entries); i++ {
		j := i
		for j > 0 && strings.Compare(entries[j-1].RelPath, entries[j].RelPath) > 0 {
			entries[j-1], entries[j] = entries[j], entries[j-1]
			j--
		}
	}
}
