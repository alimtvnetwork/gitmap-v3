package constants

// Table name for version history.
const TableRepoVersionHistory = "RepoVersionHistory"

// SQL: create RepoVersionHistory table.
const SQLCreateRepoVersionHistory = `CREATE TABLE IF NOT EXISTS RepoVersionHistory (
	Id              INTEGER PRIMARY KEY AUTOINCREMENT,
	RepoId          INTEGER NOT NULL REFERENCES Repos(Id) ON DELETE CASCADE,
	FromVersionTag  TEXT NOT NULL,
	FromVersionNum  INTEGER NOT NULL,
	ToVersionTag    TEXT NOT NULL,
	ToVersionNum    INTEGER NOT NULL,
	FlattenedPath   TEXT DEFAULT '',
	CreatedAt       TEXT DEFAULT CURRENT_TIMESTAMP
)`

// SQL: add version columns to Repos.
const (
	SQLAddCurrentVersionTag = "ALTER TABLE Repos ADD COLUMN CurrentVersionTag TEXT DEFAULT ''"
	SQLAddCurrentVersionNum = "ALTER TABLE Repos ADD COLUMN CurrentVersionNum INTEGER DEFAULT 0"
)

// SQL: version history operations.
const (
	SQLInsertVersionHistory = `INSERT INTO RepoVersionHistory
		(RepoId, FromVersionTag, FromVersionNum, ToVersionTag, ToVersionNum, FlattenedPath)
		VALUES (?, ?, ?, ?, ?, ?)`

	SQLSelectVersionHistory = `SELECT Id, RepoId, FromVersionTag, FromVersionNum,
		ToVersionTag, ToVersionNum, FlattenedPath, CreatedAt
		FROM RepoVersionHistory WHERE RepoId = ? ORDER BY CreatedAt DESC`

	SQLUpdateRepoVersion = `UPDATE Repos SET CurrentVersionTag = ?, CurrentVersionNum = ?,
		UpdatedAt = CURRENT_TIMESTAMP WHERE Id = ?`

	SQLSelectRepoIDByPath = "SELECT Id FROM Repos WHERE AbsolutePath = ?"

	SQLDropRepoVersionHistory = "DROP TABLE IF EXISTS RepoVersionHistory"
)

// Version history error messages.
const ErrDBVersionHistory = "failed to query version history: %v"

// Flatten messages.
const (
	MsgFlattenRemoving  = "Removing existing %s for fresh clone...\n"
	MsgFlattenCloning   = "Cloning %s into %s (flattened)...\n"
	MsgFlattenDone      = "✓ Cloned %s into %s\n"
	MsgFlattenVersionDB = "✓ Recorded version transition v%d -> v%d\n"
)
