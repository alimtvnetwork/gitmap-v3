package constants

// Database location.
const (
	DBDir  = "data"
	DBFile = "gitmap.db"
)

// Lock file.
const (
	LockFileName       = "gitmap.lock"
	LockFilePermission = 0o644
	ErrLockHeld        = "another gitmap process is running (PID %d).\n  If incorrect, delete: %s"
)

// Table names.
const (
	TableRepos     = "Repos"
	TableGroups    = "Groups"
	TableGroupRepo = "GroupRepos"
	TableReleases  = "Releases"
)

// SQL: create Repos table.
const SQLCreateRepos = `CREATE TABLE IF NOT EXISTS Repos (
	Id               INTEGER PRIMARY KEY AUTOINCREMENT,
	Slug             TEXT NOT NULL,
	RepoName         TEXT NOT NULL,
	HttpsUrl         TEXT NOT NULL,
	SshUrl           TEXT NOT NULL,
	Branch           TEXT NOT NULL,
	RelativePath     TEXT NOT NULL,
	AbsolutePath     TEXT NOT NULL,
	CloneInstruction TEXT NOT NULL,
	Notes            TEXT DEFAULT '',
	CreatedAt        TEXT DEFAULT CURRENT_TIMESTAMP,
	UpdatedAt        TEXT DEFAULT CURRENT_TIMESTAMP
)`

// SQL: create Groups table.
const SQLCreateGroups = `CREATE TABLE IF NOT EXISTS Groups (
	Id          INTEGER PRIMARY KEY AUTOINCREMENT,
	Name        TEXT NOT NULL UNIQUE,
	Description TEXT DEFAULT '',
	Color       TEXT DEFAULT '',
	CreatedAt   TEXT DEFAULT CURRENT_TIMESTAMP
)`

// SQL: create GroupRepos join table.
const SQLCreateGroupRepos = `CREATE TABLE IF NOT EXISTS GroupRepos (
	GroupId INTEGER NOT NULL REFERENCES Groups(Id) ON DELETE CASCADE,
	RepoId  INTEGER NOT NULL REFERENCES Repos(Id) ON DELETE CASCADE,
	PRIMARY KEY (GroupId, RepoId)
)`

// SQL: create Releases table.
const SQLCreateReleases = `CREATE TABLE IF NOT EXISTS Releases (
	Id           INTEGER PRIMARY KEY AUTOINCREMENT,
	Version      TEXT NOT NULL,
	Tag          TEXT NOT NULL UNIQUE,
	Branch       TEXT NOT NULL,
	SourceBranch TEXT NOT NULL,
	CommitSha    TEXT NOT NULL,
	Changelog    TEXT DEFAULT '',
	Notes        TEXT DEFAULT '',
	Draft        INTEGER DEFAULT 0,
	PreRelease   INTEGER DEFAULT 0,
	IsLatest     INTEGER DEFAULT 0,
	Source       TEXT DEFAULT 'release',
	CreatedAt    TEXT DEFAULT CURRENT_TIMESTAMP
)`

// SQL: add Source column to existing Releases table.
const SQLAddSourceColumn = "ALTER TABLE Releases ADD COLUMN Source TEXT DEFAULT 'release'"

// SQL: enable foreign keys.
const SQLEnableFK = "PRAGMA foreign_keys = ON"

// SQL: repo operations.
const (
	SQLUpsertRepo = `INSERT INTO Repos (Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(AbsolutePath) DO UPDATE SET
			Slug=excluded.Slug, RepoName=excluded.RepoName, HttpsUrl=excluded.HttpsUrl,
			SshUrl=excluded.SshUrl, Branch=excluded.Branch, RelativePath=excluded.RelativePath,
			CloneInstruction=excluded.CloneInstruction, Notes=excluded.Notes, UpdatedAt=CURRENT_TIMESTAMP`

	SQLSelectAllRepos = "SELECT Id, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes FROM Repos ORDER BY Slug"

	SQLSelectRepoBySlug = "SELECT Id, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes FROM Repos WHERE Slug = ?"

	SQLSelectRepoByPath = "SELECT Id, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes FROM Repos WHERE AbsolutePath = ?"
)

// SQL: upsert by AbsolutePath (spec requirement).
const SQLUpsertRepoByPath = `INSERT INTO Repos (Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(AbsolutePath) DO UPDATE SET
		Slug=excluded.Slug, RepoName=excluded.RepoName, HttpsUrl=excluded.HttpsUrl,
		SshUrl=excluded.SshUrl, Branch=excluded.Branch, RelativePath=excluded.RelativePath,
		CloneInstruction=excluded.CloneInstruction, Notes=excluded.Notes, UpdatedAt=CURRENT_TIMESTAMP`

// SQL: create unique index on AbsolutePath for upsert-by-path.
const SQLCreateAbsPathIndex = "CREATE UNIQUE INDEX IF NOT EXISTS idx_Repos_AbsolutePath ON Repos(AbsolutePath)"

// SQL: group operations.
const (
	SQLInsertGroup = "INSERT INTO Groups (Name, Description, Color) VALUES (?, ?, ?)"

	SQLSelectAllGroups = "SELECT Id, Name, Description, Color, CreatedAt FROM Groups ORDER BY Name"

	SQLSelectGroupByName = "SELECT Id, Name, Description, Color, CreatedAt FROM Groups WHERE Name = ?"

	SQLDeleteGroup = "DELETE FROM Groups WHERE Name = ?"

	SQLInsertGroupRepo = "INSERT OR IGNORE INTO GroupRepos (GroupId, RepoId) VALUES (?, ?)"

	SQLDeleteGroupRepo = "DELETE FROM GroupRepos WHERE GroupId = ? AND RepoId = ?"

	SQLSelectGroupRepos = `SELECT r.Id, r.Slug, r.RepoName, r.HttpsUrl, r.SshUrl, r.Branch,
		r.RelativePath, r.AbsolutePath, r.CloneInstruction, r.Notes
		FROM Repos r JOIN GroupRepos gr ON r.Id = gr.RepoId WHERE gr.GroupId = ? ORDER BY r.Slug`

	SQLCountGroupRepos = "SELECT COUNT(*) FROM GroupRepos WHERE GroupId = ?"
)

// SQL: release operations.
const (
	SQLUpsertRelease = `INSERT INTO Releases (Version, Tag, Branch, SourceBranch, CommitSha, Changelog, Notes, Draft, PreRelease, IsLatest, Source, CreatedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(Tag) DO UPDATE SET
			Version=excluded.Version, Branch=excluded.Branch, SourceBranch=excluded.SourceBranch,
			CommitSha=excluded.CommitSha, Changelog=excluded.Changelog, Notes=excluded.Notes, Draft=excluded.Draft,
			PreRelease=excluded.PreRelease, IsLatest=excluded.IsLatest, Source=excluded.Source`

	SQLSelectAllReleases = `SELECT Id, Version, Tag, Branch, SourceBranch, CommitSha, Changelog, Notes, Draft, PreRelease, IsLatest, Source, CreatedAt
		FROM Releases ORDER BY CreatedAt DESC`

	SQLSelectReleaseByTag = `SELECT Id, Version, Tag, Branch, SourceBranch, CommitSha, Changelog, Notes, Draft, PreRelease, IsLatest, Source, CreatedAt
		FROM Releases WHERE Tag = ?`

	SQLClearLatestRelease = "UPDATE Releases SET IsLatest = 0 WHERE IsLatest = 1"

	SQLAddNotesColumn = "ALTER TABLE Releases ADD COLUMN Notes TEXT DEFAULT ''"
)

// SQL: reset operations.
const (
	SQLDropGroupRepos = "DROP TABLE IF EXISTS GroupRepos"
	SQLDropGroups     = "DROP TABLE IF EXISTS Groups"
	SQLDropRepos      = "DROP TABLE IF EXISTS Repos"
	SQLDropReleases   = "DROP TABLE IF EXISTS Releases"
)

// Store error messages.
const (
	ErrDBOpen           = "failed to open database at %s: %v (operation: open)"
	ErrDBMigrate        = "failed to initialize tables: %v"
	ErrDBUpsert         = "failed to upsert repo: %v"
	ErrDBQuery          = "failed to query repos: %v"
	ErrDBNoMatch        = "no repo matches slug: %s\n"
	ErrDBCreateDir      = "failed to create database directory at %s: %v (operation: mkdir)"
	ErrDBGroupCreate    = "failed to create group: %v"
	ErrDBGroupQuery     = "failed to query groups: %v"
	ErrDBGroupAdd       = "failed to add repo to group: %v"
	ErrDBGroupRemove    = "failed to remove repo from group: %v"
	ErrDBGroupDelete    = "failed to delete group: %v"
	ErrDBGroupNone      = "no group found: %s"
	ErrDBGroupExists    = "group already exists: %s"
	ErrDBReleaseUpsert  = "failed to upsert release: %v"
	ErrDBReleaseQuery   = "failed to query releases: %v"
)
