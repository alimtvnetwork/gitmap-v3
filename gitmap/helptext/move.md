# `gitmap mv` / `gitmap move`

Move every file from LEFT into RIGHT, then delete LEFT entirely.
Either side can be a local folder OR a remote git URL. URL endpoints
are auto-cloned (or re-pulled if the working folder already matches),
and a commit + push is made on the URL side after the file copy.

Spec: `spec/01-app/97-move-and-merge.md`

## Usage

    gitmap mv   LEFT RIGHT [flags]
    gitmap move LEFT RIGHT [flags]

LEFT and RIGHT can each be:
- a local folder path (relative or absolute)
- a remote git URL with optional `:branch` suffix
  (`https://github.com/owner/repo:develop`)

## Flags

| Flag | Description |
|------|-------------|
| --no-push | Skip git push on URL endpoints (still commits) |
| --no-commit | Skip both commit and push on URL endpoints |
| --force-folder | Replace a folder whose origin doesn't match the URL |
| --pull | Force `git pull --ff-only` on a folder endpoint |
| --init | When RIGHT is auto-created, also `git init` it |
| --dry-run | Print every action; perform none |
| --include-vcs | Include `.git/` in the copy (default: excluded) |
| --include-node-modules | Include `node_modules/` in the copy |

## Examples

    gitmap mv ./gitmap-v3 ./gitmap-v4
    gitmap mv ./gitmap-v3 https://github.com/owner/gitmap-v4
    gitmap mv https://github.com/owner/gitmap-v3 ./another-folder
    gitmap mv https://github.com/owner/gitmap-v3 https://github.com/owner/gitmap-v4

    # preview only
    gitmap mv ./gitmap-v3 ./gitmap-v4 --dry-run

## Notes

- The `.git/` folder is never copied; LEFT's `.git/` is removed
  along with the rest of LEFT after the copy.
- LEFT and RIGHT must not resolve to the same folder, and one must
  not be nested inside the other.
- On a URL endpoint, the commit message is `gitmap mv from <LEFT>`.
