package movemerge

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/gitmap/constants"
)

// finalizeURLEndpoint stages, commits, and pushes when the endpoint
// originated from a URL. No-op for folder endpoints.
func finalizeURLEndpoint(ep Endpoint, commitMessage string, opts Options, logger *Logger) error {
	if ep.Kind != KindURL {
		return nil
	}
	if _, err := os.Stat(ep.WorkingDir); err != nil {
		// LEFT was deleted by mv — nothing to push.
		return nil
	}
	if opts.DryRun {
		logger.Logf("[dry-run] would commit & push %s with message %q", ep.DisplayName, commitMessage)

		return nil
	}

	logger.Logf("committing in %s ...", ep.DisplayName)
	if err := gitAddAll(ep.WorkingDir); err != nil {
		return err
	}
	sha, err := gitCommit(ep.WorkingDir, commitMessage)
	if err != nil {
		return err
	}
	if len(sha) == 0 {
		logger.Logf("  nothing to commit (no changes)")

		return nil
	}
	logger.Logf("  commit %s %q", sha, commitMessage)

	return pushIfRequested(ep, sha, opts, logger)
}

// pushIfRequested honours --no-commit / --no-push.
func pushIfRequested(ep Endpoint, sha string, opts Options, logger *Logger) error {
	if opts.NoCommit || opts.NoPush {
		logger.Logf("  push skipped (--no-push/--no-commit)")

		return nil
	}
	logger.Logf("pushing %s ...", ep.DisplayName)
	if err := gitPush(ep.WorkingDir, ep.Branch); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf(constants.ErrMMPushFailedFmt, sha))

		return err
	}
	logger.Logf("  push OK")

	return nil
}

// guardSamePath enforces the spec's same-folder and nested-folder
// safety checks before any file write happens.
func guardSamePath(left, right Endpoint) error {
	lAbs, _ := filepath.Abs(left.WorkingDir)
	rAbs, _ := filepath.Abs(right.WorkingDir)
	if lAbs == rAbs {
		return fmt.Errorf(constants.ErrMMSameFolderFmt, lAbs)
	}
	if isNested(lAbs, rAbs) || isNested(rAbs, lAbs) {
		return fmt.Errorf(constants.ErrMMNestedFmt, lAbs, rAbs)
	}

	return nil
}

// isNested returns true when child is strictly inside parent.
func isNested(parent, child string) bool {
	rel, err := filepath.Rel(parent, child)
	if err != nil || rel == "." || rel == "" {
		return false
	}

	return !startsWithDotDot(rel)
}

// startsWithDotDot returns true when rel starts with `..` segment.
func startsWithDotDot(rel string) bool {
	if len(rel) >= 2 && rel[0] == '.' && rel[1] == '.' {
		if len(rel) == 2 {
			return true
		}

		return rel[2] == '/' || rel[2] == '\\'
	}

	return false
}

// osStat / mkAllDir indirections keep mergeflow.go testable.
func osStat(p string) (os.FileInfo, error) { return os.Stat(p) }
func mkAllDir(p string, mode os.FileMode) error {
	return os.MkdirAll(p, mode)
}
