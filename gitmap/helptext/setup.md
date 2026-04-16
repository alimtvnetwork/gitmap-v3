# gitmap setup

Interactive first-time configuration wizard that applies global Git settings and installs shell tab-completion.

## Alias

None

## Usage

    gitmap setup [--config <path>] [--dry-run]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --config \<path\> | data/git-setup.json beside gitmap | Path to git-setup.json config file |
| --dry-run | false | Preview changes without applying them |

## Prerequisites

- Git must be installed

## Examples

### Example 1: Run the setup wizard

    gitmap setup

**Output:**

    ■ Applying global Git configuration...
      ✓ core.autocrlf = true
      ✓ push.default = current
      ✓ pull.rebase = false
    ✓ 3 Git settings applied

    ■ Shell Completion
      Detected shell: powershell
      Installing completion to $PROFILE...
    ✓ Shell completion installed for PowerShell

    ■ Setup complete!
    → Run 'gitmap scan <directory>' to start tracking repos

### Example 2: Dry-run mode (preview only)

    gitmap setup --dry-run

**Output:**

    [DRY RUN] No changes will be made
    [DRY RUN] Would set core.autocrlf = true
    [DRY RUN] Would set push.default = current
    [DRY RUN] Would set pull.rebase = false
    [DRY RUN] Would install powershell completion to $PROFILE
    No changes made.

### Example 3: Setup with custom config file

    gitmap setup --config ./my-config/git-setup.json

**Output:**

    ■ Loading config from ./my-config/git-setup.json...
    ■ Applying global Git configuration...
      ✓ core.autocrlf = true
      ✓ init.defaultBranch = main
    ✓ 2 Git settings applied
    ✓ Setup complete!

## See Also

- [completion](completion.md) — Generate completion scripts manually
- [scan](scan.md) — Scan directories after setup
- [doctor](doctor.md) — Diagnose installation issues
