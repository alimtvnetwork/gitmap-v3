package movemerge

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/user/gitmap/constants"
)

// ResolveOptions controls side effects during endpoint resolution.
type ResolveOptions struct {
	Role        string // "LEFT" or "RIGHT" (constants.EndpointLeft/Right)
	IsRightOf   string // "mv" | "merge" — affects auto-create rules
	ForceFolder bool   // --force-folder (URL endpoint replace)
	Pull        bool   // --pull (folder endpoint forced pull)
	Init        bool   // --init (mv: git init the freshly-created RIGHT folder)
	DryRun      bool
	Logger      *Logger
}

// ResolveEndpoint converts a classified Endpoint into a concrete
// WorkingDir on disk: cloning URLs, validating folders, and applying
// the spec's per-role creation rules.
//
// Per spec section "Endpoint Resolution".
func ResolveEndpoint(ep Endpoint, opts ResolveOptions) (Endpoint, error) {
	if ep.Kind == KindURL {
		return resolveURLEndpoint(ep, opts)
	}

	return resolveFolderEndpoint(ep, opts)
}

// resolveURLEndpoint handles a URL LEFT/RIGHT.
func resolveURLEndpoint(ep Endpoint, opts ResolveOptions) (Endpoint, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return ep, fmt.Errorf("getwd: %w", err)
	}
	ep.WorkingDir = CandidateWorkingDir(ep.URL, cwd)
	opts.Logger.Logf("resolving %s : %s", opts.Role, ep.DisplayName)
	opts.Logger.Logf("  -> mapped to working folder: %s", ep.WorkingDir)

	if _, statErr := os.Stat(ep.WorkingDir); os.IsNotExist(statErr) {
		return cloneURLFresh(ep, opts)
	}

	return reuseOrReplaceURLFolder(ep, opts)
}

// cloneURLFresh clones the URL into a brand-new working folder.
func cloneURLFresh(ep Endpoint, opts ResolveOptions) (Endpoint, error) {
	opts.Logger.Logf("  -> folder does not exist; cloning")
	if opts.DryRun {
		ep.WasCloned = true

		return ep, nil
	}
	if err := runGitClone(ep.URL, ep.Branch, ep.WorkingDir); err != nil {
		return ep, fmt.Errorf("clone %s: %w", ep.URL, err)
	}
	ep.WasCloned = true
	opts.Logger.Logf("  -> clone OK")

	return ep, nil
}

// reuseOrReplaceURLFolder handles the case where the candidate
// working folder already exists.
func reuseOrReplaceURLFolder(ep Endpoint, opts ResolveOptions) (Endpoint, error) {
	origin, err := readOrigin(ep.WorkingDir)
	if err != nil {
		return ep, fmt.Errorf("read origin in %s: %w", ep.WorkingDir, err)
	}
	if remotesMatch(origin, ep.URL) {
		opts.Logger.Logf("  -> existing folder origin matches; pulling --ff-only")
		if !opts.DryRun {
			if pullErr := runGitPullFFOnly(ep.WorkingDir); pullErr != nil {
				return ep, fmt.Errorf(constants.ErrMMPullFailedFmt, ep.WorkingDir, pullErr)
			}
		}

		return ep, nil
	}
	if !opts.ForceFolder {
		return ep, fmt.Errorf(constants.ErrMMRemoteMismatchFmt,
			filepath.Base(ep.WorkingDir), origin, ep.URL)
	}
	opts.Logger.Logf("  -> --force-folder set; replacing existing folder")
	if !opts.DryRun {
		if rmErr := os.RemoveAll(ep.WorkingDir); rmErr != nil {
			return ep, fmt.Errorf("remove existing folder: %w", rmErr)
		}
		if cloneErr := runGitClone(ep.URL, ep.Branch, ep.WorkingDir); cloneErr != nil {
			return ep, fmt.Errorf("re-clone %s: %w", ep.URL, cloneErr)
		}
	}
	ep.WasCloned = true

	return ep, nil
}

// resolveFolderEndpoint handles a local-folder LEFT/RIGHT.
func resolveFolderEndpoint(ep Endpoint, opts ResolveOptions) (Endpoint, error) {
	abs, err := filepath.Abs(ep.Raw)
	if err != nil {
		return ep, fmt.Errorf("abs path %s: %w", ep.Raw, err)
	}
	ep.WorkingDir = abs
	opts.Logger.Logf("resolving %s : %s (folder)", opts.Role, ep.DisplayName)

	info, statErr := os.Stat(abs)
	if statErr == nil {
		return handleExistingFolder(ep, info, opts)
	}
	if !os.IsNotExist(statErr) {
		return ep, fmt.Errorf("stat %s: %w", abs, statErr)
	}

	return handleMissingFolder(ep, opts)
}

// handleExistingFolder applies --pull when requested.
func handleExistingFolder(ep Endpoint, info os.FileInfo, opts ResolveOptions) (Endpoint, error) {
	if !info.IsDir() {
		return ep, fmt.Errorf("path %s is not a directory", ep.WorkingDir)
	}
	if opts.Pull && isGitRepo(ep.WorkingDir) {
		opts.Logger.Logf("  -> --pull set; pulling --ff-only")
		if !opts.DryRun {
			if pullErr := runGitPullFFOnly(ep.WorkingDir); pullErr != nil {
				return ep, fmt.Errorf(constants.ErrMMPullFailedFmt, ep.WorkingDir, pullErr)
			}
		}
	}

	return ep, nil
}

// handleMissingFolder enforces the per-role creation policy.
func handleMissingFolder(ep Endpoint, opts ResolveOptions) (Endpoint, error) {
	if opts.Role == constants.EndpointLeft {
		return ep, fmt.Errorf(constants.ErrMMSrcMissingFmt, ep.DisplayName)
	}
	if opts.IsRightOf == "merge" {
		return ep, fmt.Errorf(constants.ErrMMRightMissingFmt, ep.DisplayName)
	}
	// `mv` RIGHT: auto-create.
	opts.Logger.Logf("  -> RIGHT does not exist; creating %s", ep.WorkingDir)
	if !opts.DryRun {
		if mkErr := os.MkdirAll(ep.WorkingDir, 0o755); mkErr != nil {
			return ep, fmt.Errorf("mkdir %s: %w", ep.WorkingDir, mkErr)
		}
		if opts.Init {
			if initErr := runGitInit(ep.WorkingDir); initErr != nil {
				return ep, fmt.Errorf("git init %s: %w", ep.WorkingDir, initErr)
			}
		}
	}

	return ep, nil
}

// remotesMatch compares two remote URLs after normalising scheme &
// trailing `.git` so https/ssh variants of the same repo match.
func remotesMatch(a, b string) bool {
	return normalizeRemote(a) == normalizeRemote(b)
}

// normalizeRemote lowercases, strips `.git`, scheme, and `git@` so
// equivalent URLs collapse to the same canonical string.
func normalizeRemote(raw string) string {
	s := strings.ToLower(strings.TrimSpace(raw))
	s = strings.TrimSuffix(s, ".git")
	for _, prefix := range []string{"https://", "http://", "ssh://", "git@"} {
		s = strings.TrimPrefix(s, prefix)
	}
	s = strings.ReplaceAll(s, ":", "/")

	return s
}

// isGitRepo returns true when dir contains a `.git` entry.
func isGitRepo(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, ".git"))

	return err == nil
}

// readOrigin returns the URL of `origin` in the given repo.
func readOrigin(dir string) (string, error) {
	out, err := exec.Command(constants.GitBin, constants.GitDirFlag, dir,
		constants.GitConfigCmd, constants.GitGetFlag, constants.GitRemoteOrigin).Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}
