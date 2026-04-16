package constants

// SQL: create ProjectTypes table.
const SQLCreateProjectTypes = `CREATE TABLE IF NOT EXISTS ProjectTypes (
	Id          INTEGER PRIMARY KEY AUTOINCREMENT,
	Key         TEXT NOT NULL UNIQUE,
	Name        TEXT NOT NULL,
	Description TEXT DEFAULT ''
)`

// SQL: create DetectedProjects table.
const SQLCreateDetectedProjects = `CREATE TABLE IF NOT EXISTS DetectedProjects (
	Id               INTEGER PRIMARY KEY AUTOINCREMENT,
	RepoId           INTEGER NOT NULL REFERENCES Repos(Id) ON DELETE CASCADE,
	ProjectTypeId    INTEGER NOT NULL REFERENCES ProjectTypes(Id),
	ProjectName      TEXT NOT NULL,
	AbsolutePath     TEXT NOT NULL,
	RepoPath         TEXT NOT NULL,
	RelativePath     TEXT NOT NULL,
	PrimaryIndicator TEXT NOT NULL,
	DetectedAt       TEXT DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(RepoId, ProjectTypeId, RelativePath)
)`

// SQL: create GoProjectMetadata table.
const SQLCreateGoProjectMetadata = `CREATE TABLE IF NOT EXISTS GoProjectMetadata (
	Id                INTEGER PRIMARY KEY AUTOINCREMENT,
	DetectedProjectId INTEGER NOT NULL UNIQUE
		REFERENCES DetectedProjects(Id) ON DELETE CASCADE,
	GoModPath         TEXT NOT NULL,
	GoSumPath         TEXT DEFAULT '',
	ModuleName        TEXT NOT NULL,
	GoVersion         TEXT DEFAULT ''
)`

// SQL: create GoRunnableFiles table.
const SQLCreateGoRunnableFiles = `CREATE TABLE IF NOT EXISTS GoRunnableFiles (
	Id           INTEGER PRIMARY KEY AUTOINCREMENT,
	GoMetadataId INTEGER NOT NULL
		REFERENCES GoProjectMetadata(Id) ON DELETE CASCADE,
	RunnableName TEXT NOT NULL,
	FilePath     TEXT NOT NULL,
	RelativePath TEXT NOT NULL,
	UNIQUE(GoMetadataId, RelativePath)
)`

// SQL: create CSharpProjectMetadata table.
const SQLCreateCSharpProjectMeta = `CREATE TABLE IF NOT EXISTS CSharpProjectMetadata (
	Id                INTEGER PRIMARY KEY AUTOINCREMENT,
	DetectedProjectId INTEGER NOT NULL UNIQUE
		REFERENCES DetectedProjects(Id) ON DELETE CASCADE,
	SlnPath           TEXT DEFAULT '',
	SlnName           TEXT DEFAULT '',
	GlobalJsonPath    TEXT DEFAULT '',
	SdkVersion        TEXT DEFAULT ''
)`

// SQL: create CSharpProjectFiles table.
const SQLCreateCSharpProjectFiles = `CREATE TABLE IF NOT EXISTS CSharpProjectFiles (
	Id               INTEGER PRIMARY KEY AUTOINCREMENT,
	CSharpMetadataId INTEGER NOT NULL
		REFERENCES CSharpProjectMetadata(Id) ON DELETE CASCADE,
	FilePath         TEXT NOT NULL,
	RelativePath     TEXT NOT NULL,
	FileName         TEXT NOT NULL,
	ProjectName      TEXT NOT NULL,
	TargetFramework  TEXT DEFAULT '',
	OutputType       TEXT DEFAULT '',
	Sdk              TEXT DEFAULT '',
	UNIQUE(CSharpMetadataId, RelativePath)
)`

// SQL: create CSharpKeyFiles table.
const SQLCreateCSharpKeyFiles = `CREATE TABLE IF NOT EXISTS CSharpKeyFiles (
	Id               INTEGER PRIMARY KEY AUTOINCREMENT,
	CSharpMetadataId INTEGER NOT NULL
		REFERENCES CSharpProjectMetadata(Id) ON DELETE CASCADE,
	FileType         TEXT NOT NULL,
	FilePath         TEXT NOT NULL,
	RelativePath     TEXT NOT NULL,
	UNIQUE(CSharpMetadataId, RelativePath)
)`

// SQL: seed project types.
const SQLSeedProjectTypes = `INSERT OR IGNORE INTO ProjectTypes (Key, Name, Description) VALUES
	('go',     'Go',      'Go modules and packages'),
	('node',   'Node.js', 'Node.js projects'),
	('react',  'React',   'React applications'),
	('cpp',    'C++',     'C and C++ projects'),
	('csharp', 'C#',      '.NET and C# projects')`

// SQL: upsert detected project.
const SQLUpsertDetectedProject = `INSERT INTO DetectedProjects
	(RepoId, ProjectTypeId, ProjectName, AbsolutePath, RepoPath, RelativePath, PrimaryIndicator)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(RepoId, ProjectTypeId, RelativePath) DO UPDATE SET
		ProjectName=excluded.ProjectName,
		AbsolutePath=excluded.AbsolutePath,
		RepoPath=excluded.RepoPath,
		PrimaryIndicator=excluded.PrimaryIndicator,
		DetectedAt=CURRENT_TIMESTAMP`

// SQL: upsert Go metadata.
const SQLUpsertGoMetadata = `INSERT INTO GoProjectMetadata
	(DetectedProjectId, GoModPath, GoSumPath, ModuleName, GoVersion)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(DetectedProjectId) DO UPDATE SET
		GoModPath=excluded.GoModPath,
		GoSumPath=excluded.GoSumPath,
		ModuleName=excluded.ModuleName,
		GoVersion=excluded.GoVersion`

// SQL: upsert Go runnable file.
const SQLUpsertGoRunnable = `INSERT INTO GoRunnableFiles
	(GoMetadataId, RunnableName, FilePath, RelativePath)
	VALUES (?, ?, ?, ?)
	ON CONFLICT(GoMetadataId, RelativePath) DO UPDATE SET
		RunnableName=excluded.RunnableName,
		FilePath=excluded.FilePath`

// SQL: upsert C# metadata.
const SQLUpsertCSharpMetadata = `INSERT INTO CSharpProjectMetadata
	(DetectedProjectId, SlnPath, SlnName, GlobalJsonPath, SdkVersion)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(DetectedProjectId) DO UPDATE SET
		SlnPath=excluded.SlnPath,
		SlnName=excluded.SlnName,
		GlobalJsonPath=excluded.GlobalJsonPath,
		SdkVersion=excluded.SdkVersion`

// SQL: upsert C# project file.
const SQLUpsertCSharpProjectFile = `INSERT INTO CSharpProjectFiles
	(CSharpMetadataId, FilePath, RelativePath, FileName, ProjectName, TargetFramework, OutputType, Sdk)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(CSharpMetadataId, RelativePath) DO UPDATE SET
		FilePath=excluded.FilePath,
		FileName=excluded.FileName,
		ProjectName=excluded.ProjectName,
		TargetFramework=excluded.TargetFramework,
		OutputType=excluded.OutputType,
		Sdk=excluded.Sdk`

// SQL: upsert C# key file.
const SQLUpsertCSharpKeyFile = `INSERT INTO CSharpKeyFiles
	(CSharpMetadataId, FileType, FilePath, RelativePath)
	VALUES (?, ?, ?, ?)
	ON CONFLICT(CSharpMetadataId, RelativePath) DO UPDATE SET
		FileType=excluded.FileType,
		FilePath=excluded.FilePath`

// SQL: query detected project ID by identity tuple.
const SQLSelectDetectedProjectID = `SELECT Id
	FROM DetectedProjects
	WHERE RepoId = ? AND ProjectTypeId = ? AND RelativePath = ?`

// SQL: query projects by type key.
const SQLSelectProjectsByTypeKey = `SELECT dp.Id, dp.RepoId, pt.Key, dp.ProjectName,
	dp.AbsolutePath, dp.RepoPath, dp.RelativePath,
	dp.PrimaryIndicator, dp.DetectedAt, r.RepoName
	FROM DetectedProjects dp
	JOIN ProjectTypes pt ON dp.ProjectTypeId = pt.Id
	JOIN Repos r ON dp.RepoId = r.Id
	WHERE pt.Key = ?
	ORDER BY r.RepoName, dp.RelativePath`

// SQL: count projects by type key.
const SQLCountProjectsByTypeKey = `SELECT COUNT(*)
	FROM DetectedProjects dp
	JOIN ProjectTypes pt ON dp.ProjectTypeId = pt.Id
	WHERE pt.Key = ?`

// SQL: query Go metadata.
const SQLSelectGoMetadata = `SELECT Id, DetectedProjectId, GoModPath, GoSumPath,
	ModuleName, GoVersion
	FROM GoProjectMetadata WHERE DetectedProjectId = ?`

// SQL: query Go runnables.
const SQLSelectGoRunnables = `SELECT Id, GoMetadataId, RunnableName, FilePath,
	RelativePath
	FROM GoRunnableFiles WHERE GoMetadataId = ?
	ORDER BY RunnableName`

// SQL: query C# metadata.
const SQLSelectCSharpMetadata = `SELECT Id, DetectedProjectId, SlnPath, SlnName,
	GlobalJsonPath, SdkVersion
	FROM CSharpProjectMetadata WHERE DetectedProjectId = ?`

// SQL: query C# project files.
const SQLSelectCSharpProjectFiles = `SELECT Id, CSharpMetadataId, FilePath,
	RelativePath, FileName, ProjectName, TargetFramework, OutputType, Sdk
	FROM CSharpProjectFiles WHERE CSharpMetadataId = ?
	ORDER BY RelativePath`

// SQL: query C# key files.
const SQLSelectCSharpKeyFiles = `SELECT Id, CSharpMetadataId, FileType, FilePath,
	RelativePath
	FROM CSharpKeyFiles WHERE CSharpMetadataId = ?
	ORDER BY RelativePath`

// SQL: stale cleanup.
const (
	SQLDeleteStaleProjects       = "DELETE FROM DetectedProjects WHERE RepoId = ? AND Id NOT IN (%s)"
	SQLDeleteStaleGoRunnables    = "DELETE FROM GoRunnableFiles WHERE GoMetadataId = ? AND Id NOT IN (%s)"
	SQLDeleteStaleCSharpFiles    = "DELETE FROM CSharpProjectFiles WHERE CSharpMetadataId = ? AND Id NOT IN (%s)"
	SQLDeleteStaleCSharpKeyFiles = "DELETE FROM CSharpKeyFiles WHERE CSharpMetadataId = ? AND Id NOT IN (%s)"
)

// SQL: drop project detection tables.
const (
	SQLDropGoRunnableFiles    = "DROP TABLE IF EXISTS GoRunnableFiles"
	SQLDropGoProjectMetadata  = "DROP TABLE IF EXISTS GoProjectMetadata"
	SQLDropCSharpKeyFiles     = "DROP TABLE IF EXISTS CSharpKeyFiles"
	SQLDropCSharpProjectFiles = "DROP TABLE IF EXISTS CSharpProjectFiles"
	SQLDropCSharpProjectMeta  = "DROP TABLE IF EXISTS CSharpProjectMetadata"
	SQLDropDetectedProjects   = "DROP TABLE IF EXISTS DetectedProjects"
	SQLDropProjectTypes       = "DROP TABLE IF EXISTS ProjectTypes"
)
