# `gitmap merge-both` / `merge-left` / `merge-right`

File-level merge between two endpoints (folders or URLs). Files
present on only one side are copied to the other; files that exist
on both sides with different content trigger a conflict prompt.

Spec: `spec/01-app/97-move-and-merge.md`

## Usage

    gitmap merge-both  LEFT RIGHT [flags]   # both sides gain missing files
    gitmap merge-left  LEFT RIGHT [flags]   # write into LEFT only
    gitmap merge-right LEFT RIGHT [flags]   # write into RIGHT only

## Conflict prompt

For each conflict you'll see file sizes + mtimes and:

    [L]eft  [R]ight  [S]kip  [A]ll-left  [B]all-right  [Q]uit

| Key | Action |
|-----|--------|
| L | Take LEFT's version |
| R | Take RIGHT's version |
| S | Skip this file |
| A | All-Left for the rest |
| B | All-Right for the rest |
| Q | Quit immediately (already-applied changes are kept) |

## Bypass flags

| Flag | Effect |
|------|--------|
| -y, --yes, -a, --accept-all | Bypass prompt; apply per-command default |
| --prefer-left | LEFT always wins on conflict |
| --prefer-right | RIGHT always wins on conflict |
| --prefer-newer | Newer mtime wins (default for merge-both with -y) |
| --prefer-skip | Skip every conflict; only missing files are copied |

Per-command `-y` defaults: `merge-right` -> LEFT wins,
`merge-left` -> RIGHT wins, `merge-both` -> newer mtime wins.

## Other flags

`--no-push`, `--no-commit`, `--force-folder`, `--pull`, `--dry-run`,
`--include-vcs`, `--include-node-modules` — see `gitmap mv help`.

## Examples

    gitmap merge-both ./gitmap-v3 ./gitmap-v4
    gitmap merge-right ./local https://github.com/owner/repo -y
    gitmap merge-both ./local https://github.com/owner/repo -y --prefer-newer
    gitmap merge-right ./local https://github.com/owner/repo:develop
