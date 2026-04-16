package constants

// Setup section headers.
const (
	SetupSectionDiff  = "Diff Tool"
	SetupSectionMerge = "Merge Tool"
	SetupSectionAlias = "Aliases"
	SetupSectionCred  = "Credential Helper"
	SetupSectionCore  = "Core Settings"
	SetupSectionComp  = "■ Shell Completion —"
	SetupGlobalFlag   = "--global"
)

// Release messages.
const (
	MsgReleaseStart         = "\n  Creating release %s...\n"
	MsgReleaseBranch        = "  ✓ Created branch %s\n"
	MsgReleaseTag           = "  ✓ Created tag %s\n"
	MsgReleasePushed        = "  ✓ Pushed branch and tag to origin\n"
	MsgReleaseMeta          = "  ✓ Release metadata written to %s\n"
	MsgReleaseMetaCommitted = "  ✓ Committed release metadata on %s\n"
	MsgReleaseLatest        = "  ✓ Marked %s as latest release\n"
	MsgReleaseAttach        = "  ✓ Attached %s\n"
	MsgReleaseChangelog     = "  ✓ Using CHANGELOG.md as release body\n"
	MsgReleaseReadme        = "  ✓ Attached README.md\n"
	MsgReleaseDryRun        = "  [dry-run] %s\n"
	MsgReleaseComplete      = "\n  Release %s complete.\n"
	MsgReleaseBranchStart   = "\n  Completing release from %s...\n"
	MsgReleaseBranchPending = "\n  → On release branch %s with no tag — completing pending release...\n"
	MsgReleaseVersionRead   = "  → Version from %s: %s\n"
	MsgReleaseBumpResult    = "  → Bumped %s → %s\n"
	MsgReleaseNotes         = "  → Release notes: %s\n"
	MsgReleaseSwitchedBack  = "  ✓ Switched back to %s\n"
	MsgReleasePendingNone   = "  No pending release branches found."
	MsgReleasePendingFound  = "\n  Found %d pending release branch(es).\n"
	MsgReleasePendingFailed = "  ✗ Failed to release %s: %v\n"
	ReleaseBranchPrefix     = "release/"
	ChangelogFile           = "CHANGELOG.md"
	ReadmeFile              = "README.md"
	ReleaseTagPrefix        = "Release "
	FlagDescNotes           = "Release notes or title for the release"
)

// Release orphaned metadata messages.
const (
	MsgReleaseOrphanedMeta    = "  ⚠ Release metadata exists for %s but no tag or branch was found.\n"
	MsgReleaseOrphanedPrompt  = "  → Do you want to remove the release JSON and proceed? (y/N): "
	MsgReleaseOrphanedRemoved = "  ✓ Removed orphaned release metadata for %s\n"
	ErrReleaseOrphanedRemove  = "failed to remove release metadata at %s: %w (operation: delete)"
	ErrReleaseAborted         = "release aborted by user"
)

// Self-release messages.
const (
	MsgSelfReleaseSwitch      = "\n  → Self-release: switching to %s\n"
	MsgSelfReleaseReturn      = "  ✓ Returned to %s\n"
	MsgSelfReleaseSameDir     = "\n  → Self-release: already in source repo %s\n"
	MsgSelfReleasePromptPath  = "  → Enter gitmap source repo path: "
	MsgSelfReleaseSavedPath   = "  ✓ Saved gitmap source repo path: %s\n"
	MsgSelfReleaseInvalidPath = "  ✗ Invalid gitmap source repo path: %s\n"
	ErrSelfReleaseExec        = "could not resolve executable path at %s: %w (operation: resolve)"
	ErrSelfReleaseNoRepo      = "could not locate gitmap source repository"
)

// Install hint constants (printed after release for gitmap repos).
const (
	GitmapRepoPrefix     = "github.com/alimtvnetwork/gitmap-v2"
	MsgInstallHintHeader = `

  📦 Install gitmap %s
`
	MsgInstallHintWindows = `  🪟 Windows (PowerShell)
     irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v2/main/gitmap/scripts/install.ps1 | iex
`
	MsgInstallHintUnix = `
  🐧 Linux / macOS
     curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v2/main/gitmap/scripts/install.sh | sh
`
)

// Release rollback messages.
const (
	MsgRollbackStart  = "\n  ⚠ Push failed — rolling back local branch and tag...\n"
	MsgRollbackBranch = "  ✓ Deleted local branch %s\n"
	MsgRollbackTag    = "  ✓ Deleted local tag %s\n"
	MsgRollbackDone   = "  ✓ Rollback complete. No changes remain.\n"
	MsgRollbackWarn   = "  ⚠ Rollback warning (%s): %v\n"
)
