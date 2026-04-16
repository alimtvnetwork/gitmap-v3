package constants

// Alias command names.
const (
	CmdAlias        = "alias"
	CmdAliasShort   = "a"
	SubCmdAliasSet  = "set"
	SubCmdAliasRm   = "remove"
	SubCmdAliasList = "list"
	SubCmdAliasShow = "show"
	SubCmdAliasSug  = "suggest"
)

// Alias table name.
const TableAliases = "Aliases"

// SQL: create Aliases table.
const SQLCreateAliases = `CREATE TABLE IF NOT EXISTS Aliases (
	Id        INTEGER PRIMARY KEY AUTOINCREMENT,
	Alias     TEXT NOT NULL UNIQUE,
	RepoId    INTEGER NOT NULL REFERENCES Repos(Id) ON DELETE CASCADE,
	CreatedAt TEXT DEFAULT CURRENT_TIMESTAMP
)`

// SQL: alias operations.
const (
	SQLInsertAlias = `INSERT INTO Aliases (Alias, RepoId) VALUES (?, ?)`

	SQLUpdateAlias = `UPDATE Aliases SET RepoId = ? WHERE Alias = ?`

	SQLSelectAllAliases = `SELECT a.Id, a.Alias, a.RepoId, a.CreatedAt
		FROM Aliases a ORDER BY a.Alias`

	SQLSelectAliasByName = `SELECT a.Id, a.Alias, a.RepoId, a.CreatedAt
		FROM Aliases a WHERE a.Alias = ?`

	SQLSelectAliasByRepoID = `SELECT a.Id, a.Alias, a.RepoId, a.CreatedAt
		FROM Aliases a WHERE a.RepoId = ?`

	SQLDeleteAlias = `DELETE FROM Aliases WHERE Alias = ?`

	SQLSelectAliasWithRepo = `SELECT a.Id, a.Alias, a.RepoId, a.CreatedAt,
		r.AbsolutePath, r.Slug
		FROM Aliases a JOIN Repos r ON a.RepoId = r.Id
		WHERE a.Alias = ?`

	SQLSelectAllAliasesWithRepo = `SELECT a.Id, a.Alias, a.RepoId, a.CreatedAt,
		r.AbsolutePath, r.Slug
		FROM Aliases a JOIN Repos r ON a.RepoId = r.Id
		ORDER BY a.Alias`

	SQLSelectUnaliasedRepos = `SELECT r.Id, r.Slug, r.RepoName
		FROM Repos r LEFT JOIN Aliases a ON r.Id = a.RepoId
		WHERE a.Id IS NULL ORDER BY r.Slug`
)

// SQL: drop Aliases table.
const SQLDropAliases = "DROP TABLE IF EXISTS Aliases"

// Alias flag descriptions.
const (
	FlagDescAliasApply = "Auto-accept all alias suggestions"
	FlagDescAliasFlag  = "Target a repository by its alias"
)

// Alias messages.
const (
	MsgAliasCreated     = "  ✓ Alias %q → %s\n"
	MsgAliasUpdated     = "  ✓ Updated alias %q → %s\n"
	MsgAliasRemoved     = "  ✓ Removed alias %q\n"
	MsgAliasResolved    = "  → Resolved alias %q → %s (slug: %s)\n"
	MsgAliasSuggest     = "  %-20s → %-10s Accept? (y/N): "
	MsgAliasSuggestDone = "  ✓ Created %d alias(es).\n"
	MsgAliasSuggestNone = "  All repos already have aliases."
	MsgAliasListHeader  = "\n  Aliases (%d):\n\n"
	MsgAliasListRow     = "  %-15s → %s\n"
	MsgAliasConflict    = "  ⚠ Alias %q already points to %s.\n"
	MsgAliasReassign    = "  → Reassign to %s? (y/N): "
	MsgAliasBothWarn    = "  ⚠ Both alias and slug provided — using alias %q.\n"
)

// Alias error messages.
const (
	ErrAliasNotFound    = "no alias found: %s"
	ErrAliasEmpty       = "alias name cannot be empty"
	ErrAliasInvalid     = "alias must be alphanumeric with hyphens: %s"
	ErrAliasShadow      = "alias cannot shadow command: %s"
	ErrAliasCreate      = "failed to create alias: %v"
	ErrAliasQuery       = "failed to query aliases: %v"
	ErrAliasDelete      = "failed to delete alias: %v"
	ErrAliasRepoMissing = "repo not found for alias target: %s"
)
