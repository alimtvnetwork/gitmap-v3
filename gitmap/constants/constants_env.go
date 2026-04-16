package constants

// Env CLI commands.
const (
	CmdEnv      = "env"
	CmdEnvAlias = "ev"
)

// Env subcommands.
const (
	CmdEnvSet        = "set"
	CmdEnvGet        = "get"
	CmdEnvDelete     = "delete"
	CmdEnvList       = "list"
	CmdEnvPathAdd    = "path"
	CmdEnvPathSub    = "add"
	CmdEnvPathRemove = "remove"
	CmdEnvPathList   = "list"
)

// Env help text.
const HelpEnv = "  env (ev) <sub>      Manage environment variables and PATH"

// Env registry file.
const (
	EnvRegistryFileName = "env-registry.json"
	EnvRegistryFilePath = GitMapDir + "/" + EnvRegistryFileName
)

// Env flag names.
const (
	FlagEnvSystem  = "system"
	FlagEnvShell   = "shell"
	FlagEnvVerbose = "verbose"
	FlagEnvDryRun  = "dry-run"
)

// Env flag descriptions.
const (
	FlagDescEnvSystem  = "Target system-level variables (Windows, requires admin)"
	FlagDescEnvShell   = "Target shell profile: bash, zsh (Unix only)"
	FlagDescEnvVerbose = "Show detailed operation output"
	FlagDescEnvDryRun  = "Preview changes without applying"
)

// Env shell profile paths.
const (
	EnvProfileBashRC  = ".bashrc"
	EnvProfileZshRC   = ".zshrc"
	EnvProfileBash    = ".bash_profile"
	EnvExportPrefix   = "export "
	EnvExportFmt      = "export %s=\"%s\""
	EnvPathExportFmt  = "export PATH=\"$PATH:%s\""
	EnvManagedComment = "# managed by gitmap"
)

// Env terminal messages.
const (
	MsgEnvSet         = "Set %s=%s\n"
	MsgEnvDeleted     = "Removed %s\n"
	MsgEnvPathAdded   = "Added to PATH: %s\n"
	MsgEnvPathRemoved = "Removed from PATH: %s\n"
	MsgEnvListHeader  = "Managed variables:\n"
	MsgEnvListRow     = "  %s = %s\n"
	MsgEnvListEmpty   = "No managed variables. Use 'gitmap env set' to add one.\n"
	MsgEnvPathHeader  = "Managed PATH entries:\n"
	MsgEnvPathRow     = "  %s\n"
	MsgEnvPathEmpty   = "No managed PATH entries.\n"
	MsgEnvDrySet      = "[dry-run] Would set %s=%s\n"
	MsgEnvDryPath     = "[dry-run] Would add to PATH: %s\n"
	MsgEnvDryDelete   = "[dry-run] Would remove %s\n"
	MsgEnvGetFmt      = "%s=%s\n"
)

// Env error messages.
const (
	ErrEnvNameRequired   = "Variable name is required."
	ErrEnvValueRequired  = "Variable value is required."
	ErrEnvPathRequired   = "PATH entry is required."
	ErrEnvNotFound       = "Variable '%s' is not managed by gitmap.\n"
	ErrEnvPathNotExist   = "Error: directory does not exist at %s (operation: resolve, reason: file does not exist)\n"
	ErrEnvInvalidName    = "Invalid variable name: %s (must be alphanumeric and underscore only)\n"
	ErrEnvProfileWrite   = "Error: failed to write shell profile at %s: %v (operation: write)\n"
	ErrEnvRegistryLoad   = "Error: failed to load env registry at %s: %v (operation: read)\n"
	ErrEnvRegistrySave   = "Error: failed to save env registry at %s: %v (operation: write)\n"
	ErrEnvSubcommand     = "Unknown env subcommand: %s\n"
	ErrEnvSystemWindows  = "System-level variables require administrator privileges."
	ErrEnvPathDuplicate  = "PATH entry already exists: %s\n"
)