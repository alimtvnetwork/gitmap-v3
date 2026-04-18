package cmd

import (
	"os"

	"github.com/user/gitmap/constants"
)

// dispatchMoveMerge routes the move/merge command family.
//
// Spec: spec/01-app/97-move-and-merge.md
func dispatchMoveMerge(command string) bool {
	if command == constants.CmdMove || command == constants.CmdMoveAlias {
		runMove(os.Args[2:])

		return true
	}
	if isMergeCommand(command) {
		runMerge(command, os.Args[2:])

		return true
	}

	return false
}

// isMergeCommand returns true for the three merge variants + aliases.
func isMergeCommand(command string) bool {
	switch command {
	case constants.CmdMergeBoth, constants.CmdMergeBothA,
		constants.CmdMergeLeft, constants.CmdMergeLeftA,
		constants.CmdMergeRight, constants.CmdMergeRightA:
		return true
	}

	return false
}
