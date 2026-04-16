package constants

// Pending task table names.
const (
	TableTaskType      = "TaskType"
	TablePendingTask   = "PendingTask"
	TableCompletedTask = "CompletedTask"
)

// Pending task type seed values.
const (
	TaskTypeDelete = "Delete"
	TaskTypeRemove = "Remove"
	TaskTypeScan   = "Scan"
	TaskTypeClone  = "Clone"
	TaskTypePull   = "Pull"
	TaskTypeExec   = "Exec"
)

// SQL: create TaskType table.
const SQLCreateTaskType = `CREATE TABLE IF NOT EXISTS TaskType (
	Id   INTEGER PRIMARY KEY AUTOINCREMENT,
	Name TEXT NOT NULL UNIQUE
)`

// SQL: create PendingTask table.
const SQLCreatePendingTask = `CREATE TABLE IF NOT EXISTS PendingTask (
	Id               INTEGER PRIMARY KEY AUTOINCREMENT,
	TaskTypeId       INTEGER NOT NULL REFERENCES TaskType(Id),
	TargetPath       TEXT    NOT NULL,
	WorkingDirectory TEXT    DEFAULT '',
	SourceCommand    TEXT    NOT NULL,
	CommandArgs      TEXT    DEFAULT '',
	FailureReason    TEXT    DEFAULT '',
	CreatedAt        TEXT    DEFAULT CURRENT_TIMESTAMP,
	UpdatedAt        TEXT    DEFAULT CURRENT_TIMESTAMP
)`

// SQL: create CompletedTask table.
const SQLCreateCompletedTask = `CREATE TABLE IF NOT EXISTS CompletedTask (
	Id               INTEGER PRIMARY KEY AUTOINCREMENT,
	OriginalTaskId   INTEGER NOT NULL,
	TaskTypeId       INTEGER NOT NULL REFERENCES TaskType(Id),
	TargetPath       TEXT    NOT NULL,
	WorkingDirectory TEXT    DEFAULT '',
	SourceCommand    TEXT    NOT NULL,
	CommandArgs      TEXT    DEFAULT '',
	CompletedAt      TEXT    DEFAULT CURRENT_TIMESTAMP,
	CreatedAt        TEXT    NOT NULL
)`

// SQL: seed TaskType values.
const SQLSeedTaskTypes = `INSERT OR IGNORE INTO TaskType (Name)
	VALUES ('Delete'), ('Remove'), ('Scan'), ('Clone'), ('Pull'), ('Exec')`

// SQL: drop pending task tables.
const (
	SQLDropCompletedTask = "DROP TABLE IF EXISTS CompletedTask"
	SQLDropPendingTask   = "DROP TABLE IF EXISTS PendingTask"
	SQLDropTaskType      = "DROP TABLE IF EXISTS TaskType"
)

// SQL: migrate existing tables to add new columns.
const (
	SQLMigratePendingWorkDir   = "ALTER TABLE PendingTask ADD COLUMN WorkingDirectory TEXT DEFAULT ''"
	SQLMigratePendingCmdArgs   = "ALTER TABLE PendingTask ADD COLUMN CommandArgs TEXT DEFAULT ''"
	SQLMigrateCompletedWorkDir = "ALTER TABLE CompletedTask ADD COLUMN WorkingDirectory TEXT DEFAULT ''"
	SQLMigrateCompletedCmdArgs = "ALTER TABLE CompletedTask ADD COLUMN CommandArgs TEXT DEFAULT ''"
)
