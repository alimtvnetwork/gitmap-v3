// Package movemerge implements `gitmap mv` and the `gitmap merge-*`
// command family.
//
// Spec: spec/01-app/97-move-and-merge.md
package movemerge

import (
	"path/filepath"
	"strings"

	"github.com/user/gitmap/constants"
)

// EndpointKind classifies an argument as a URL or local folder.
type EndpointKind int

const (
	// KindFolder = local file-system path.
	KindFolder EndpointKind = iota
	// KindURL = git remote URL (https/http/ssh/git@).
	KindURL
)

// Endpoint is a fully classified LEFT or RIGHT argument.
//
// After ResolveEndpoint runs, WorkingDir is the absolute path the
// caller should operate on (clone target for URLs, the path itself
// for folders). DisplayName is the original CLI argument (for log
// lines and commit messages).
type Endpoint struct {
	Kind        EndpointKind
	Raw         string // exactly what the user typed
	DisplayName string // Raw without trailing slash
	URL         string // populated for KindURL (without :branch)
	Branch      string // populated when URL had `:branch` suffix
	WorkingDir  string // absolute path on disk after resolution
	WasCloned   bool   // true when WorkingDir was created by us
}

// ClassifyEndpoint inspects a raw CLI argument and decides whether
// it points at a URL or a folder. It does NOT touch the filesystem
// or the network — that work happens in ResolveEndpoint.
func ClassifyEndpoint(raw string) Endpoint {
	trimmed := strings.TrimRight(raw, "/\\")
	ep := Endpoint{Raw: raw, DisplayName: trimmed}

	if isURL(trimmed) {
		ep.Kind = KindURL
		ep.URL, ep.Branch = splitURLAndBranch(trimmed)

		return ep
	}

	ep.Kind = KindFolder

	return ep
}

// isURL returns true when s starts with a known git URL scheme.
func isURL(s string) bool {
	lower := strings.ToLower(s)
	prefixes := []string{"https://", "http://", "ssh://", "git@"}
	for _, p := range prefixes {
		if strings.HasPrefix(lower, p) {
			return true
		}
	}

	return false
}

// splitURLAndBranch separates a trailing `:branch` suffix from a URL.
// SSH-shorthand (`git@host:owner/repo`) is detected so the colon
// inside it is NOT treated as a branch separator.
func splitURLAndBranch(raw string) (url, branch string) {
	if strings.HasPrefix(strings.ToLower(raw), "git@") {
		// Allow `:branch` only AFTER the first slash.
		slash := strings.Index(raw, "/")
		if slash < 0 {
			return raw, ""
		}
		head, tail := raw[:slash], raw[slash:]
		if idx := strings.LastIndex(tail, constants.BranchSuffixSep); idx > 0 {
			return head + tail[:idx], tail[idx+1:]
		}

		return raw, ""
	}
	// http/https/ssh schemes: a colon AFTER `://...` separates the branch.
	if scheme := strings.Index(raw, "://"); scheme >= 0 {
		body := raw[scheme+3:]
		if idx := strings.LastIndex(body, constants.BranchSuffixSep); idx > 0 {
			return raw[:scheme+3+idx], body[idx+1:]
		}
	}

	return raw, ""
}

// RepoNameFromURL derives the bare repository name (no `.git`,
// no version suffix, no path) from a URL — used as the candidate
// working folder for URL endpoints.
func RepoNameFromURL(url string) string {
	name := strings.TrimSuffix(url, ".git")
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		name = name[idx+1:]
	}
	if idx := strings.LastIndex(name, ":"); idx >= 0 {
		name = name[idx+1:]
	}

	return name
}

// CandidateWorkingDir returns the absolute path that a URL endpoint
// would map to in the current working directory.
func CandidateWorkingDir(url, cwd string) string {
	return filepath.Join(cwd, RepoNameFromURL(url))
}
