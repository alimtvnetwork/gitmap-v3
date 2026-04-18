package movemerge

import (
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
)

// Options bundles every runtime knob shared by the four operations.
type Options struct {
	NoPush         bool
	NoCommit       bool
	ForceFolder    bool
	Pull           bool
	Init           bool
	DryRun         bool
	AutoMode       string // "" | left | right | newer | skip
	IncludeVCS     bool
	IncludeNodeMod bool
}

// Run dispatches one of the four operations by command keyword.
//
// command: constants.CmdMove | CmdMergeBoth | CmdMergeLeft | CmdMergeRight
func Run(command string, leftRaw, rightRaw string, opts Options) error {
	left := ClassifyEndpoint(leftRaw)
	right := ClassifyEndpoint(rightRaw)
	logger := NewLogger(logPrefixFor(command))

	rolePolicy := rolePolicyFor(command)
	leftRes, err := ResolveEndpoint(left, ResolveOptions{
		Role: constants.EndpointLeft, IsRightOf: rolePolicy,
		ForceFolder: opts.ForceFolder, Pull: opts.Pull,
		DryRun: opts.DryRun, Logger: logger,
	})
	if err != nil {
		return err
	}
	rightRes, err := ResolveEndpoint(right, ResolveOptions{
		Role: constants.EndpointRight, IsRightOf: rolePolicy,
		ForceFolder: opts.ForceFolder, Pull: opts.Pull, Init: opts.Init,
		DryRun: opts.DryRun, Logger: logger,
	})
	if err != nil {
		return err
	}
	if guardErr := guardSamePath(leftRes, rightRes); guardErr != nil {
		return guardErr
	}

	return dispatchOperation(command, leftRes, rightRes, opts, logger)
}

// dispatchOperation routes to the per-command implementation.
func dispatchOperation(command string, left, right Endpoint, opts Options, logger *Logger) error {
	switch command {
	case constants.CmdMove, constants.CmdMoveAlias:
		return doMove(left, right, opts, logger)
	case constants.CmdMergeBoth, constants.CmdMergeBothA:
		return doMergeBoth(left, right, opts, logger)
	case constants.CmdMergeLeft, constants.CmdMergeLeftA:
		return doMergeLeft(left, right, opts, logger)
	case constants.CmdMergeRight, constants.CmdMergeRightA:
		return doMergeRight(left, right, opts, logger)
	}

	return fmt.Errorf("unknown command: %s", command)
}

// doMove implements `gitmap mv`: copy LEFT -> RIGHT (overwriting),
// then delete LEFT entirely.
func doMove(left, right Endpoint, opts Options, logger *Logger) error {
	logger.Logf("copying files LEFT -> RIGHT (excluding .git/) ...")
	copyOpts := CopyOptions{
		IncludeVCS: opts.IncludeVCS, IncludeNodeModules: opts.IncludeNodeMod,
		DryRun: opts.DryRun,
	}
	res, err := CopyTree(left.WorkingDir, right.WorkingDir, copyOpts)
	if err != nil {
		return fmt.Errorf("copy: %w", err)
	}
	logger.Logf("  copied %d files", res.FilesCopied)

	logger.Logf("deleting LEFT (%s) ...", left.DisplayName)
	if !opts.DryRun {
		if rmErr := os.RemoveAll(left.WorkingDir); rmErr != nil {
			return fmt.Errorf("remove LEFT: %w", rmErr)
		}
	}
	logger.Logf("  deleted")

	return finalizeURLEndpoint(right, fmt.Sprintf(constants.CommitMsgMvFmt, left.DisplayName), opts, logger)
}

// doMergeBoth implements `gitmap merge-both`.
func doMergeBoth(left, right Endpoint, opts Options, logger *Logger) error {
	autoMode := opts.AutoMode
	if autoMode == "" && opts.NoCommit { // -y default for merge-both = newer
		autoMode = constants.PreferNewer
	}
	if err := runMergeFlow(left, right, opts, logger, mergeFlowConfig{
		WriteLeft: true, WriteRight: true, DefaultAuto: constants.PreferNewer,
	}); err != nil {
		return err
	}
	msg := fmt.Sprintf(constants.CommitMsgMergeBothFmt, left.DisplayName)
	if err := finalizeURLEndpoint(right, msg, opts, logger); err != nil {
		return err
	}
	leftMsg := fmt.Sprintf(constants.CommitMsgMergeBothFmt, right.DisplayName)

	return finalizeURLEndpoint(left, leftMsg, opts, logger)
}

// doMergeLeft implements `gitmap merge-left`: write into LEFT only.
func doMergeLeft(left, right Endpoint, opts Options, logger *Logger) error {
	if err := runMergeFlow(left, right, opts, logger, mergeFlowConfig{
		WriteLeft: true, WriteRight: false, DefaultAuto: constants.PreferRight,
	}); err != nil {
		return err
	}
	msg := fmt.Sprintf(constants.CommitMsgMergeLeftFmt, right.DisplayName)

	return finalizeURLEndpoint(left, msg, opts, logger)
}

// doMergeRight implements `gitmap merge-right`: write into RIGHT only.
func doMergeRight(left, right Endpoint, opts Options, logger *Logger) error {
	if err := runMergeFlow(left, right, opts, logger, mergeFlowConfig{
		WriteLeft: false, WriteRight: true, DefaultAuto: constants.PreferLeft,
	}); err != nil {
		return err
	}
	msg := fmt.Sprintf(constants.CommitMsgMergeRightFmt, left.DisplayName)

	return finalizeURLEndpoint(right, msg, opts, logger)
}

// rolePolicyFor maps a command to the IsRightOf policy ("mv"|"merge").
func rolePolicyFor(command string) string {
	if command == constants.CmdMove || command == constants.CmdMoveAlias {
		return "mv"
	}

	return "merge"
}

// logPrefixFor returns the structured log prefix for a command.
func logPrefixFor(command string) string {
	switch command {
	case constants.CmdMove, constants.CmdMoveAlias:
		return constants.LogPrefixMv
	case constants.CmdMergeBoth, constants.CmdMergeBothA:
		return constants.LogPrefixMergeBoth
	case constants.CmdMergeLeft, constants.CmdMergeLeftA:
		return constants.LogPrefixMergeLeft
	case constants.CmdMergeRight, constants.CmdMergeRightA:
		return constants.LogPrefixMergeRight
	}

	return "[gitmap]"
}
