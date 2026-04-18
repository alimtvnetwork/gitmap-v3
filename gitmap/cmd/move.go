package cmd

import (
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/movemerge"
)

// runMove handles `gitmap mv` / `gitmap move`.
//
// Spec: spec/01-app/97-move-and-merge.md
func runMove(args []string) {
	checkHelp("move", args)
	left, right, opts := parseMoveMergeArgs(constants.CmdMove, args)
	if err := movemerge.Run(constants.CmdMove, left, right, opts); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, constants.LogPrefixMv+" done")
}
