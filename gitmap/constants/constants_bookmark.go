package constants

// Bookmark table name.
const TableBookmarks = "Bookmarks"

// SQL: create Bookmarks table.
const SQLCreateBookmarks = `CREATE TABLE IF NOT EXISTS Bookmarks (
	Id        INTEGER PRIMARY KEY AUTOINCREMENT,
	Name      TEXT NOT NULL UNIQUE,
	Command   TEXT NOT NULL,
	Args      TEXT DEFAULT '',
	Flags     TEXT DEFAULT '',
	CreatedAt TEXT DEFAULT CURRENT_TIMESTAMP
)`

// SQL: bookmark operations.
const (
	SQLInsertBookmark = `INSERT INTO Bookmarks (Name, Command, Args, Flags)
		VALUES (?, ?, ?, ?)`

	SQLSelectAllBookmarks = `SELECT Id, Name, Command, Args, Flags, CreatedAt
		FROM Bookmarks ORDER BY Name`

	SQLSelectBookmarkByName = `SELECT Id, Name, Command, Args, Flags, CreatedAt
		FROM Bookmarks WHERE Name = ?`

	SQLDeleteBookmark = "DELETE FROM Bookmarks WHERE Name = ?"

	SQLDropBookmarks = "DROP TABLE IF EXISTS Bookmarks"
)

// Bookmark CLI commands.
const (
	CmdBookmark      = "bookmark"
	CmdBookmarkAlias = "bk"
)

// Bookmark subcommands.
const (
	CmdBookmarkSave   = "save"
	CmdBookmarkList   = "list"
	CmdBookmarkRun    = "run"
	CmdBookmarkDelete = "delete"
)

// Bookmark help text.
const (
	HelpBookmark = "  bookmark (bk) <sub> Save and replay command+flag combinations (save, list, run, delete)"
)

// Bookmark messages.
const (
	MsgBookmarkSaved      = "Bookmark saved: %s → gitmap %s %s %s\n"
	MsgBookmarkDeleted    = "Bookmark deleted: %s\n"
	MsgBookmarkEmpty      = "No bookmarks saved.\n"
	MsgBookmarkRunning    = "Running bookmark: %s → gitmap %s %s %s\n"
	MsgBookmarkColumns    = "NAME                 COMMAND         ARGS             FLAGS"
	MsgBookmarkRowFmt     = "%-20s %-15s %-16s %s\n"
	ErrBookmarkUsage      = "usage: gitmap bookmark <save|list|run|delete> [args]\n"
	ErrBookmarkSaveUsage  = "usage: gitmap bookmark save <name> <command> [args...] [--flags...]\n"
	ErrBookmarkRunUsage   = "usage: gitmap bookmark run <name>\n"
	ErrBookmarkDelUsage   = "usage: gitmap bookmark delete <name>\n"
	ErrBookmarkNotFound   = "bookmark not found: %s\n"
	ErrBookmarkExists     = "bookmark already exists: %s (delete it first)\n"
	ErrBookmarkQuery      = "failed to query bookmarks: %v"
	ErrBookmarkSave       = "failed to save bookmark: %v\n"
	ErrBookmarkDelete     = "failed to delete bookmark: %v\n"
)
