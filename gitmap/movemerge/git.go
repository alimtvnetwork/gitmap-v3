package movemerge

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/user/gitmap/constants"
)

// runGitClone clones url into dir, optionally on a specific branch.
// When branch is set and missing on the remote, callers handle the
// "create new branch" path post-clone.
func runGitClone(url, branch, dir string) error {
	args := []string{constants.GitClone}
	if len(branch) > 0 {
		args = append(args, constants.GitBranchFlag, branch)
	}
	args = append(args, url, dir)
	cmd := exec.Command(constants.GitBin, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git %s: %w\n%s", strings.Join(args, " "), err, string(out))
	}

	return nil
}

// runGitPullFFOnly runs `git pull --ff-only` inside dir.
func runGitPullFFOnly(dir string) error {
	cmd := exec.Command(constants.GitBin, constants.GitDirFlag, dir,
		constants.GitPull, constants.GitFFOnlyFlag)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git pull: %w\n%s", err, string(out))
	}

	return nil
}

// runGitInit runs `git init` inside dir.
func runGitInit(dir string) error {
	cmd := exec.Command(constants.GitBin, constants.GitDirFlag, dir, "init")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git init: %w\n%s", err, string(out))
	}

	return nil
}

// gitAddAll stages every change in dir.
func gitAddAll(dir string) error {
	cmd := exec.Command(constants.GitBin, constants.GitDirFlag, dir, "add", "-A")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git add: %w\n%s", err, string(out))
	}

	return nil
}

// gitCommit creates a commit with the given message. Returns the
// short SHA of the new commit. Empty commits are reported as a
// no-op (returns "" with no error).
func gitCommit(dir, message string) (string, error) {
	if hasNoStaged(dir) {
		return "", nil
	}
	cmd := exec.Command(constants.GitBin, constants.GitDirFlag, dir, "commit", "-m", message)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git commit: %w\n%s", err, string(out))
	}

	return readHeadShort(dir)
}

// hasNoStaged returns true when there's nothing to commit.
func hasNoStaged(dir string) bool {
	cmd := exec.Command(constants.GitBin, constants.GitDirFlag, dir,
		"diff", "--cached", "--quiet")

	return cmd.Run() == nil
}

// readHeadShort returns the short SHA of HEAD.
func readHeadShort(dir string) (string, error) {
	cmd := exec.Command(constants.GitBin, constants.GitDirFlag, dir,
		constants.GitRevParse, "--short", constants.GitHEAD)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// gitPush pushes the current branch. When branch is non-empty and
// the upstream is missing, --set-upstream is used.
func gitPush(dir, branch string) error {
	args := []string{constants.GitDirFlag, dir, constants.GitPush}
	if len(branch) > 0 {
		args = append(args, "--set-upstream", constants.GitOrigin, branch)
	}
	cmd := exec.Command(constants.GitBin, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git push: %w\n%s", err, string(out))
	}

	return nil
}
