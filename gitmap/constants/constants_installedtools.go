package constants

// SQL: create InstalledTools table.
const SQLCreateInstalledTools = `CREATE TABLE IF NOT EXISTS InstalledTools (
	Id             INTEGER PRIMARY KEY AUTOINCREMENT,
	Tool           TEXT NOT NULL UNIQUE,
	VersionMajor   INTEGER NOT NULL DEFAULT 0,
	VersionMinor   INTEGER NOT NULL DEFAULT 0,
	VersionPatch   INTEGER NOT NULL DEFAULT 0,
	VersionBuild   INTEGER NOT NULL DEFAULT 0,
	VersionString  TEXT NOT NULL DEFAULT '',
	PackageManager TEXT NOT NULL DEFAULT '',
	InstallPath    TEXT NOT NULL DEFAULT '',
	InstalledAt    TEXT NOT NULL DEFAULT '',
	UpdatedAt      TEXT NOT NULL DEFAULT ''
)`

// SQL: InstalledTools queries.
const (
	SQLInsertInstalledTool = `INSERT OR REPLACE INTO InstalledTools
		(Tool, VersionMajor, VersionMinor, VersionPatch, VersionBuild, VersionString, PackageManager, InstallPath, InstalledAt, UpdatedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))`

	SQLSelectInstalledTool = `SELECT Id, Tool, VersionMajor, VersionMinor, VersionPatch, VersionBuild, VersionString, PackageManager, InstallPath, InstalledAt, UpdatedAt FROM InstalledTools WHERE Tool = ?`
	SQLSelectAllInstalled  = `SELECT Id, Tool, VersionMajor, VersionMinor, VersionPatch, VersionBuild, VersionString, PackageManager, InstallPath, InstalledAt, UpdatedAt FROM InstalledTools ORDER BY Tool`
	SQLDeleteInstalledTool = `DELETE FROM InstalledTools WHERE Tool = ?`
	SQLExistsInstalledTool = `SELECT COUNT(*) FROM InstalledTools WHERE Tool = ?`
	SQLDropInstalledTools  = `DROP TABLE IF EXISTS InstalledTools`
)
