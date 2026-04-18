package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/movemerge"
)

// parseMoveMergeArgs parses the shared flag set for mv/merge-* and
// returns LEFT, RIGHT, and a populated movemerge.Options.
//
// Spec: spec/01-app/97-move-and-merge.md
func parseMoveMergeArgs(command string, args []string) (left, right string, opts movemerge.Options) {
	fs := flag.NewFlagSet(command, flag.ExitOnError)
	noPush := fs.Bool(constants.FlagMMNoPush, false, "skip git push on URL endpoints")
	noCommit := fs.Bool(constants.FlagMMNoCommit, false, "skip commit & push on URL endpoints")
	forceFolder := fs.Bool(constants.FlagMMForceFolder, false, "replace mismatched URL-mapped folder")
	pull := fs.Bool(constants.FlagMMPull, false, "force git pull --ff-only on folder endpoints")
	initFlag := fs.Bool(constants.FlagMMInit, false, "git init when RIGHT is auto-created")
	dryRun := fs.Bool(constants.FlagMMDryRun, false, "print actions; perform none")
	yesLong := fs.Bool(constants.FlagMMYes, false, "bypass conflict prompt (source side wins)")
	yesShort := fs.Bool(constants.FlagMMYesShort, false, "shorthand for --yes")
	acceptLong := fs.Bool(constants.FlagMMAcceptAll, false, "alias for --yes")
	acceptShort := fs.Bool(constants.FlagMMAcceptAllShort, false, "shorthand for --accept-all")
	preferLeft := fs.Bool(constants.FlagMMPreferLeft, false, "LEFT wins on conflict (with -y)")
	preferRight := fs.Bool(constants.FlagMMPreferRight, false, "RIGHT wins on conflict (with -y)")
	preferNewer := fs.Bool(constants.FlagMMPreferNewer, false, "newer mtime wins on conflict (with -y)")
	preferSkip := fs.Bool(constants.FlagMMPreferSkip, false, "skip every conflict (with -y)")
	includeVCS := fs.Bool(constants.FlagMMIncludeVCS, false, "include .git/ in copy/diff")
	includeNM := fs.Bool(constants.FlagMMIncludeNodeMods, false, "include node_modules/ in copy/diff")

	positional := reorderFlagsBeforeArgs(args)
	if err := fs.Parse(positional); err != nil {
		os.Exit(2)
	}
	left, right = extractTwoPositional(fs.Args(), command)

	opts = movemerge.Options{
		NoPush: *noPush, NoCommit: *noCommit, ForceFolder: *forceFolder,
		Pull: *pull, Init: *initFlag, DryRun: *dryRun,
		IncludeVCS: *includeVCS, IncludeNodeMod: *includeNM,
	}
	bypass := *yesLong || *yesShort || *acceptLong || *acceptShort
	opts.AutoMode = resolveAutoMode(command, bypass,
		*preferLeft, *preferRight, *preferNewer, *preferSkip)

	return left, right, opts
}

// extractTwoPositional pulls LEFT and RIGHT from the parsed args.
func extractTwoPositional(rest []string, command string) (string, string) {
	if len(rest) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: gitmap %s LEFT RIGHT [flags]\n", command)
		os.Exit(2)
	}

	return rest[0], rest[1]
}

// resolveAutoMode maps the bypass + --prefer-* flags to a string mode.
// Returns "" for interactive prompting.
func resolveAutoMode(command string, bypass, pl, pr, pn, ps bool) string {
	switch {
	case pl:
		return constants.PreferLeft
	case pr:
		return constants.PreferRight
	case pn:
		return constants.PreferNewer
	case ps:
		return constants.PreferSkip
	}
	if !bypass {
		return ""
	}

	return defaultBypassMode(command)
}

// defaultBypassMode picks the per-command default for plain `-y`.
func defaultBypassMode(command string) string {
	switch command {
	case constants.CmdMergeRight, constants.CmdMergeRightA:
		return constants.PreferLeft
	case constants.CmdMergeLeft, constants.CmdMergeLeftA:
		return constants.PreferRight
	case constants.CmdMergeBoth, constants.CmdMergeBothA:
		return constants.PreferNewer
	}

	return constants.PreferLeft // mv: source-side wins (LEFT)
}
