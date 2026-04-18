package cmd

import (
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/movemerge"
)

// runMerge handles `gitmap merge-both | merge-left | merge-right`.
//
// Spec: spec/01-app/97-move-and-merge.md
func runMerge(command string, args []string) {
	checkHelp(command, args)
	left, right, opts := parseMoveMergeArgs(command, args)
	if err := movemerge.Run(command, left, right, opts); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, mergeLogPrefix(command)+" done")
}

// mergeLogPrefix maps a command to its log prefix.
func mergeLogPrefix(command string) string {
	switch command {
	case constants.CmdMergeBoth, constants.CmdMergeBothA:
		return constants.LogPrefixMergeBoth
	case constants.CmdMergeLeft, constants.CmdMergeLeftA:
		return constants.LogPrefixMergeLeft
	case constants.CmdMergeRight, constants.CmdMergeRightA:
		return constants.LogPrefixMergeRight
	}

	return "[merge]"
}
