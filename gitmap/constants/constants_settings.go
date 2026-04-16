package constants

// Settings table.
const TableSettings = "Settings"

// SQL: create Settings table.
const SQLCreateSettings = `CREATE TABLE IF NOT EXISTS Settings (
	Key   TEXT PRIMARY KEY,
	Value TEXT NOT NULL
)`

// SQL: settings operations.
const (
	SQLUpsertSetting = `INSERT INTO Settings (Key, Value) VALUES (?, ?)
		ON CONFLICT(Key) DO UPDATE SET Value=excluded.Value`

	SQLSelectSetting = "SELECT Value FROM Settings WHERE Key = ?"

	SQLDeleteSetting = "DELETE FROM Settings WHERE Key = ?"
)

// SQL: reset.
const SQLDropSettings = "DROP TABLE IF EXISTS Settings"

// Settings keys.
const (
	SettingActiveGroup      = "active_group"
	SettingActiveMultiGroup = "active_multi_group"
	SettingSourceRepoPath   = "source_repo_path"
)

// Settings error messages.
const (
	ErrDBSettingUpsert = "failed to save setting: %v"
	ErrDBSettingQuery  = "failed to read setting: %v"
)
