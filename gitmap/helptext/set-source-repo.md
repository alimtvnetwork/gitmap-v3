# gitmap set-source-repo

Save the source repository path to the database. Used by `run.ps1` after deploy
to keep the repo path current for future `gitmap update` calls.

## Alias

None (hidden command)

## Usage

    gitmap set-source-repo <path>

## Arguments

| Argument | Description |
|----------|-------------|
| `<path>` | Absolute or relative path to the gitmap source repository root |

## Behavior

1. Validates the path is a valid gitmap source repo (has `.git/`, correct module marker).
2. Normalizes to an absolute path.
3. Persists to the Settings table (`source_repo_path` key).
4. Prints confirmation message.

## Examples

### Example 1: Set repo path after moving the source

    gitmap set-source-repo D:\Projects\gitmap-v3

**Output:**

    ✓ Source repo path saved: D:\Projects\gitmap-v3

### Example 2: Invalid path

    gitmap set-source-repo C:\nonexistent

**Output:**

    ✗ Invalid source repo path: C:\nonexistent

### Example 3: No path argument

    gitmap set-source-repo

**Output:**

    ✗ set-source-repo requires a path argument

### Example 4: Called by run.ps1 after deploy (automated)

During `run.ps1` deploy, the script runs this automatically:

    gitmap.exe set-source-repo D:\gitmap-v3

**Output (suppressed by run.ps1):**

    ✓ Source repo path saved: D:\gitmap-v3

The deploy log shows:

    -> Source repo path synced to DB: D:\gitmap-v3

## Why This Exists

When the source repo is moved or cloned to a new location, the database still
stores the old path. Without this command, `gitmap update` would fail or prompt
the user unnecessarily. The post-deploy sync ensures the path stays current.

## See Also

- [update](update.md) — Self-update using the saved repo path
- [doctor](doctor.md) — Diagnose installation issues
