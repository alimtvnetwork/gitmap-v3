package constants

// Move & merge command identifiers.
//
// Spec: spec/01-app/97-move-and-merge.md
const (
	CmdMove        = "move"
	CmdMoveAlias   = "mv"
	CmdMergeBoth   = "merge-both"
	CmdMergeLeft   = "merge-left"
	CmdMergeRight  = "merge-right"
	CmdMergeBothA  = "mb"
	CmdMergeLeftA  = "ml"
	CmdMergeRightA = "mr"
)

// Move & merge flag names.
const (
	FlagMMNoPush          = "no-push"
	FlagMMNoCommit        = "no-commit"
	FlagMMForceFolder     = "force-folder"
	FlagMMPull            = "pull"
	FlagMMInit            = "init"
	FlagMMDryRun          = "dry-run"
	FlagMMYes             = "yes"
	FlagMMYesShort        = "y"
	FlagMMAcceptAll       = "accept-all"
	FlagMMAcceptAllShort  = "a"
	FlagMMPreferLeft      = "prefer-left"
	FlagMMPreferRight     = "prefer-right"
	FlagMMPreferNewer     = "prefer-newer"
	FlagMMPreferSkip      = "prefer-skip"
	FlagMMIncludeVCS      = "include-vcs"
	FlagMMIncludeNodeMods = "include-node-modules"
)

// Conflict resolution choices (single-letter prompt keys).
const (
	ConflictKeyLeft     = "L"
	ConflictKeyRight    = "R"
	ConflictKeySkip     = "S"
	ConflictKeyAllLeft  = "A"
	ConflictKeyAllRight = "B"
	ConflictKeyQuit     = "Q"
)

// Conflict resolution preferences (used with --prefer-* / -y bypass).
const (
	PreferLeft   = "left"
	PreferRight  = "right"
	PreferNewer  = "newer"
	PreferSkip   = "skip"
	PreferPrompt = "prompt"
)

// Endpoint role labels for logging.
const (
	EndpointLeft  = "LEFT"
	EndpointRight = "RIGHT"
)

// Log prefixes for move/merge operations.
const (
	LogPrefixMv         = "[mv]"
	LogPrefixMergeBoth  = "[merge-both]"
	LogPrefixMergeLeft  = "[merge-left]"
	LogPrefixMergeRight = "[merge-right]"
)

// Commit message templates (Sprintf with the "other" endpoint display).
const (
	CommitMsgMvFmt          = "gitmap mv from %s"
	CommitMsgMergeBothFmt   = "gitmap merge-both with %s"
	CommitMsgMergeLeftFmt   = "gitmap merge-left from %s"
	CommitMsgMergeRightFmt  = "gitmap merge-right from %s"
)

// Conflict prompt strings.
const (
	ConflictPromptHeaderFmt = "  conflict: %s\n"
	ConflictPromptMetaFmt   = "    %-5s : %s  modified %s\n"
	ConflictPromptKeysLine  = "  [L]eft  [R]ight  [S]kip  [A]ll-left  [B]all-right  [Q]uit\n  > "
	ConflictAppliedFmt      = "  conflict %s -> took %s"
	ConflictSkippedFmt      = "  conflict %s -> skipped"
)

// Error & abort messages.
const (
	ErrMMSameFolderFmt    = "error: LEFT and RIGHT resolve to the same folder: %s"
	ErrMMNestedFmt        = "error: RIGHT is nested inside LEFT (or vice versa): %s vs %s"
	ErrMMSrcMissingFmt    = "error: source %q does not exist"
	ErrMMRightMissingFmt  = "error: merge target %q does not exist"
	ErrMMRemoteMismatchFmt = "error: folder %q exists but its remote is %q, not %q. " +
		"Pass --force-folder to overwrite, or rename it."
	ErrMMPullFailedFmt = "error: pull --ff-only failed in %s: %v"
	ErrMMPushFailedFmt = "Push failed. Local commit is preserved at %s. " +
		"Resolve manually or re-run with --no-push to skip."
	ErrMMUnknownEndpointFmt = "error: cannot classify endpoint %q"
)

// Default ignore directories (relative path basenames).
var DefaultMoveMergeIgnoreDirs = []string{
	".git",
	"node_modules",
}

// Default ignore-path-prefixes (relative to endpoint root).
var DefaultMoveMergeIgnorePrefixes = []string{
	".gitmap/release-assets/",
}

// Branch suffix separator for URL endpoints (e.g. https://x/y:develop).
const BranchSuffixSep = ":"
