// Package completion generates shell tab-completion scripts for gitmap.
package completion

import (
	"fmt"

	"github.com/user/gitmap/constants"
)

// Generate returns the completion script for the given shell.
func Generate(shell string) (string, error) {
	switch shell {
	case constants.ShellPowerShell:
		return generatePowerShell(), nil
	case constants.ShellBash:
		return generateBash(), nil
	case constants.ShellZsh:
		return generateZsh(), nil
	default:
		return "", fmt.Errorf(constants.ErrCompUnknownShell, shell)
	}
}

// AllCommands returns every command name and alias for completion.
func AllCommands() []string {
	return []string{
		"scan", "s",
		"clone", "c",
		"pull", "p",
		"status", "st",
		"exec", "x",
		"release", "r",
		"release-branch", "rb",
		"release-pending", "rp",
		"changelog", "cl",
		"latest-branch", "lb",
		"list", "ls",
		"group", "g",
		"multi-group", "mg",
		"cd", "go",
		"update",
		"version", "v",
		"desktop-sync", "ds",
		"rescan", "rs",
		"setup",
		"doctor",
		"db-reset",
		"list-versions", "lv",
		"list-releases", "lr",
		"revert",
		"seo-write", "sw",
		"amend", "am",
		"amend-list", "al",
		"history", "hi",
		"history-reset", "hr",
		"stats", "ss",
		"bookmark", "bk",
		"export", "ep",
		"import", "im",
		"profile", "pf",
		"diff-profiles", "dp",
		"watch", "w",
		"gomod", "gm",
		"go-repos", "gor",
		"node-repos", "nr",
		"react-repos", "rr",
		"cpp-repos", "cr",
		"csharp-repos", "csr",
		"completion", "cmp",
		"interactive", "i",
		"clear-release-json", "crj",
		"alias", "a",
		"zip-group", "z",
		"dashboard", "db",
		"ssh",
		"prune", "pr",
		"temp-release", "tr",
		"clone-next", "cn",
		"uninstall", "un",
		"help",
		"version-history", "vh",
		"llm-docs", "ld",
	}
}
