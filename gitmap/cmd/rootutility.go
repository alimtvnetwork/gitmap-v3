package cmd

import (
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
)

// dispatchUtility routes setup, update, doctor, and other utility commands.
func dispatchUtility(command string) bool {
	if command == constants.CmdUpdate {
		checkHelp("update", os.Args[2:])
		runUpdate()

		return true
	}
	if command == constants.CmdUpdateRunner {
		runUpdateRunner()

		return true
	}
	if command == constants.CmdUpdateCleanup {
		runUpdateCleanup()

		return true
	}
	if command == constants.CmdInstalledDir || command == constants.CmdInstalledDirAlias {
		checkHelp("installed-dir", os.Args[2:])
		runInstalledDir()

		return true
	}
	if command == constants.CmdRevert {
		runRevert(os.Args[2:])

		return true
	}
	if command == constants.CmdRevertRunner {
		runRevertRunner()

		return true
	}
	if command == constants.CmdVersion || command == constants.CmdVersionAlias {
		checkHelp("version", os.Args[2:])
		fmt.Printf(constants.MsgVersionFmt, constants.Version)

		return true
	}
	if command == constants.CmdHelp {
		if hasFlag(constants.FlagGroups) {
			printHelpGroups()

			return true
		}
		if hasFlag(constants.FlagCompact) {
			printUsageCompact()

			return true
		}
		printUsage()

		return true
	}
	if command == constants.CmdDocs || command == constants.CmdDocsAlias {
		runDocs(os.Args[2:])

		return true
	}
	if command == constants.CmdHelpDashboard || command == constants.CmdHelpDashboardAlias {
		runHelpDashboard(os.Args[2:])

		return true
	}
	if command == constants.CmdLLMDocs || command == constants.CmdLLMDocsAlias {
		runLLMDocs(os.Args[2:])

		return true
	}

	return false
}
