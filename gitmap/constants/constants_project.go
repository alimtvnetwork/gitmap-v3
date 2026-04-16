package constants

// Project type IDs matching seed order in SQLSeedProjectTypes.
const (
	ProjectTypeGoID     int64 = 1
	ProjectTypeNodeID   int64 = 2
	ProjectTypeReactID  int64 = 3
	ProjectTypeCppID    int64 = 4
	ProjectTypeCSharpID int64 = 5
)

// Project type keys.
const (
	ProjectKeyGo     = "go"
	ProjectKeyNode   = "node"
	ProjectKeyReact  = "react"
	ProjectKeyCpp    = "cpp"
	ProjectKeyCSharp = "csharp"
)

// Project detection table names.
const (
	TableProjectTypes       = "ProjectTypes"
	TableDetectedProjects   = "DetectedProjects"
	TableGoProjectMetadata  = "GoProjectMetadata"
	TableGoRunnableFiles    = "GoRunnableFiles"
	TableCSharpProjectMeta  = "CSharpProjectMetadata"
	TableCSharpProjectFiles = "CSharpProjectFiles"
	TableCSharpKeyFiles     = "CSharpKeyFiles"
)

// Project JSON output filenames.
const (
	JSONFileGoProjects     = "go-projects.json"
	JSONFileNodeProjects   = "node-projects.json"
	JSONFileReactProjects  = "react-projects.json"
	JSONFileCppProjects    = "cpp-projects.json"
	JSONFileCSharpProjects = "csharp-projects.json"
)

// Detection indicator files.
const (
	IndicatorGoMod       = "go.mod"
	IndicatorPackageJSON = "package.json"
	IndicatorCMakeLists  = "CMakeLists.txt"
	IndicatorMesonBuild  = "meson.build"
)

// Detection file extensions.
const (
	ExtCsproj  = ".csproj"
	ExtFsproj  = ".fsproj"
	ExtVcxproj = ".vcxproj"
	ExtSln     = ".sln"
)

// Go structural indicators.
const (
	GoCmdDir       = "cmd"
	GoMainFile     = "main.go"
	GoSumFile      = "go.sum"
	CMakeBuildPfx  = "cmake-build-"
)

// Project query commands.
const (
	CmdGoRepos         = "go-repos"
	CmdGoReposAlias    = "gr"
	CmdNodeRepos       = "node-repos"
	CmdNodeReposAlias  = "nr"
	CmdReactRepos      = "react-repos"
	CmdReactReposAlias = "rr"
	CmdCppRepos        = "cpp-repos"
	CmdCppReposAlias   = "cr"
	CmdCSharpRepos     = "csharp-repos"
	CmdCSharpAlias     = "csr"
)

// Project query flags.
const (
	FlagProjectJSON  = "json"
	FlagProjectCount = "count"
)

// Project query help text.
const (
	HelpGoRepos     = "  go-repos (gr)       List repositories containing Go projects"
	HelpNodeRepos   = "  node-repos (nr)     List repositories containing Node.js projects"
	HelpReactRepos  = "  react-repos (rr)    List repositories containing React projects"
	HelpCppRepos    = "  cpp-repos (cr)      List repositories containing C++ projects"
	HelpCSharpRepos = "  csharp-repos (csr)  List repositories containing C# projects"
)

// Project detection messages.
const (
	MsgProjectDetectDone   = "Detected %d projects across %d repos\n"
	MsgProjectUpsertDone   = "Saved %d detected projects to database\n"
	MsgProjectJSONWritten  = "Wrote %s (%d records)\n"
	MsgProjectNoDB         = "No database found. Run 'gitmap scan' first.\n"
	MsgProjectNoneFound    = "No %s projects found.\n"
	MsgProjectCount        = "%d\n"
	MsgProjectCleanedStale = "Cleaned %d stale project records\n"
	MsgProjectListCount    = "\n%d projects found.\n"
)

// Project detection error messages.
const (
	ErrProjectDetect       = "failed to detect projects in %s: %v\n"
	ErrProjectUpsert       = "failed to upsert detected project: %v"
	ErrProjectQuery        = "failed to query projects: %v"
	ErrProjectJSONWrite    = "failed to write %s: %v (operation: write)\n"
	ErrProjectParseMod     = "failed to parse go.mod in %s: %v\n"
	ErrProjectParsePkgJSON = "failed to parse package.json in %s: %v\n"
	ErrProjectParseCsproj  = "failed to parse .csproj in %s: %v\n"
	ErrProjectCleanup      = "failed to clean stale projects for repo %d: %v\n"
	ErrGoMetadataUpsert    = "failed to upsert Go metadata: %v"
	ErrGoRunnableUpsert    = "failed to upsert Go runnable: %v"
	ErrCSharpMetaUpsert    = "failed to upsert C# metadata: %v"
	ErrCSharpFileUpsert    = "failed to upsert C# project file: %v"
	ErrCSharpKeyUpsert     = "failed to upsert C# key file: %v"
)

// Legacy data recovery messages.
const (
	MsgLegacyProjectData = "Database contains legacy project data from a previous version.\n" +
		"To fix, run one of:\n\n" +
		"  gitmap rescan          Re-scan repos and rebuild project data\n" +
		"  gitmap db-reset --confirm   Reset the entire database\n"
)

// React indicator dependencies.
var ReactIndicatorDeps = []string{
	"react",
	"@types/react",
	"react-scripts",
	"next",
	"gatsby",
	"remix",
	"@remix-run/react",
}

// C# key file patterns.
var CSharpKeyFilePatterns = []string{
	"global.json",
	"nuget.config",
	"Directory.Build.props",
	"Directory.Build.targets",
	"launchSettings.json",
	"appsettings.json",
}

// Project detection exclusion directories.
var ProjectExcludeDirs = []string{
	"node_modules",
	"vendor",
	".git",
	"dist",
	"build",
	"target",
	"bin",
	"obj",
	"out",
	"testdata",
	"packages",
	".venv",
	".cache",
}
