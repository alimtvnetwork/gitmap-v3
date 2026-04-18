---
name: Move & Merge Commands
description: gitmap mv / merge-both / merge-left / merge-right — file-level move/merge between any two endpoints (folder OR URL), with L/R/S/A/B/Q conflict prompt, --prefer-* bypass, and URL-side commit+push (v2.96.0)
type: feature
---
The move/merge command family (gitmap mv, merge-both, merge-left,
merge-right) operates on two endpoints (LEFT, RIGHT) where each can
be a local folder OR a remote git URL with optional `:branch`
suffix. URL endpoints are auto-cloned (or re-pulled if the working
folder origin matches), and a commit + push happens on the URL side
after the file operation.

Implementation lives in `gitmap/movemerge/`:
- `endpoint.go` / `endpoint_resolve.go` — classify and resolve each
  endpoint (URL vs folder, :branch parsing, --force-folder replace,
  --pull, --init).
- `copy.go` — copy file tree excluding .git/ and node_modules/,
  preserving modes and symlinks.
- `diff.go` — pairwise SHA-256 diff producing
  DiffMissingLeft/Right, DiffIdentical, DiffConflict.
- `conflict.go` — L/R/S/A/B/Q interactive prompt + sticky All-Left/
  All-Right + --prefer-* / -y bypass policy.
- `operations.go` + `mergeflow.go` — Run() dispatcher and the
  per-command orchestrators.
- `git.go` + `finalize.go` — git clone/pull/add/commit/push helpers
  and same-folder/nested-folder safety guards.

CLI surface in `gitmap/cmd/move.go`, `gitmap/cmd/merge.go`,
`gitmap/cmd/movemergeflags.go`. Dispatch is wired through
`gitmap/cmd/dispatchmovemerge.go`.

Per-command `-y` defaults: merge-right -> LEFT wins, merge-left ->
RIGHT wins, merge-both -> newer mtime wins, mv -> LEFT wins (no
prompt anyway since mv is documented destructive).

Spec: `spec/01-app/97-move-and-merge.md` (all 14 acceptance items
ticked).
