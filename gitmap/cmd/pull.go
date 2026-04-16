package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/gitmap/cloner"
	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
	"github.com/user/gitmap/verbose"
)

// runPull handles the "pull" subcommand.
func runPull(args []string) {
	checkHelp("pull", args)
	requireOnline()
	slug, groupName, all, verboseMode, stopOnFail := parsePullFlags(args)
	if verboseMode {
		initVerboseLog()
	}
	records := resolvePullTargets(slug, groupName, all)

	workDir, wdErr := os.Getwd()
	if wdErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not determine working directory: %v\n", wdErr)
	}
	cmdArgs := buildCommandArgs(append([]string{"pull"}, os.Args[2:]...))
	targetPath := workDir
	if len(records) == 1 {
		targetPath = records[0].AbsolutePath
	}

	taskID, taskDB := createPendingTask(constants.TaskTypePull, targetPath, workDir, "pull", cmdArgs)
	if taskDB != nil {
		defer taskDB.Close()
	}

	prog := cloner.NewBatchProgress(len(records), "Pull", false)
	prog.SetStopOnFail(stopOnFail)
	for _, rec := range records {
		if prog.Stopped() {
			break
		}
		prog.BeginItem(rec.RepoName)
		pullOneRepoTracked(rec, prog)
	}
	prog.PrintSummary()
	prog.PrintFailureReport()

	if code := prog.ExitCodeForBatch(); code != 0 {
		failPendingTask(taskDB, taskID, fmt.Sprintf("pull batch failed with exit code %d", code))
		os.Exit(code)
	}

	completePendingTask(taskDB, taskID)
}

// parsePullFlags parses flags for the pull command.
func parsePullFlags(args []string) (slug, group string, all, verboseFlag, stopOnFail bool) {
	fs := flag.NewFlagSet(constants.CmdPull, flag.ExitOnError)
	vFlag := fs.Bool("verbose", false, constants.FlagDescVerbose)
	gFlag := fs.String("group", "", constants.FlagDescGroup)
	fs.StringVar(gFlag, "g", "", constants.FlagDescGroup)
	aFlag := fs.Bool("all", false, constants.FlagDescAll)
	sFlag := fs.Bool(constants.FlagStopOnFail, false, constants.FlagDescStopOnFail)
	fs.Parse(args)

	if fs.NArg() > 0 {
		slug = fs.Arg(0)
	}

	return slug, *gFlag, *aFlag, *vFlag, *sFlag
}

// initVerboseLog sets up verbose logging, warning on failure.
func initVerboseLog() {
	log, err := verbose.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnVerboseLogFailed, err)

		return
	}
	log.Close()
}

// resolvePullTargets returns records based on alias, group, all, or slug lookup.
func resolvePullTargets(slug, groupName string, all bool) []model.ScanRecord {
	if HasAlias() {
		return []model.ScanRecord{{
			RepoName:     GetAliasSlug(),
			Slug:         GetAliasSlug(),
			AbsolutePath: GetAliasPath(),
		}}
	}
	if len(groupName) > 0 {
		return loadRecordsByGroup(groupName)
	}
	if all {
		return loadAllRecordsDB()
	}
	if len(slug) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrPullSlugRequired)
		fmt.Fprintln(os.Stderr, constants.ErrPullUsage)
		os.Exit(1)
	}

	return lookupBySlugDBFirst(slug)
}

// lookupBySlugDBFirst tries the database first, then falls back to JSON.
func lookupBySlugDBFirst(slug string) []model.ScanRecord {
	db, err := openDB()
	if err == nil {
		defer db.Close()
		repos, dbErr := db.FindBySlug(strings.ToLower(slug))
		if dbErr == nil && len(repos) > 0 {
			return repos
		}
	}

	return lookupBySlugJSON(slug)
}

// lookupBySlugJSON loads gitmap.json and matches by repo name.
func lookupBySlugJSON(slug string) []model.ScanRecord {
	jsonPath := filepath.Join(constants.DefaultOutputFolder, constants.DefaultJSONFile)
	records, err := loadJSONRecords(jsonPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrPullLoadFailed, jsonPath, err)
		os.Exit(1)
	}

	return findBySlug(records, slug)
}

// loadJSONRecords reads ScanRecords from a JSON file.
func loadJSONRecords(path string) ([]model.ScanRecord, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []model.ScanRecord
	err = json.NewDecoder(file).Decode(&records)

	return records, err
}

// findBySlug finds records matching the slug (case-insensitive, partial match).
func findBySlug(records []model.ScanRecord, slug string) []model.ScanRecord {
	slugLower := strings.ToLower(slug)
	exact, partial := partitionBySlug(records, slugLower)

	if len(exact) > 0 {
		return exact
	}

	return partial
}

// partitionBySlug separates records into exact and partial matches.
func partitionBySlug(records []model.ScanRecord, slugLower string) ([]model.ScanRecord, []model.ScanRecord) {
	var exact, partial []model.ScanRecord

	for _, r := range records {
		nameLower := strings.ToLower(r.RepoName)
		if nameLower == slugLower {
			exact = append(exact, r)
		} else if strings.Contains(nameLower, slugLower) {
			partial = append(partial, r)
		}
	}

	return exact, partial
}

// pullOneRepo runs safe-pull on a single repo using its absolute path.
func pullOneRepo(rec model.ScanRecord) {
	fmt.Printf(constants.MsgPullStarting, rec.RepoName, rec.AbsolutePath)

	if cloner.IsMissingRepo(rec.AbsolutePath) {
		fmt.Fprintf(os.Stderr, constants.ErrPullNotRepo, rec.AbsolutePath)

		return
	}

	result := cloner.SafePullOne(rec, rec.AbsolutePath)
	if result.Success {
		fmt.Printf(constants.MsgPullSuccess, rec.RepoName)
	} else {
		fmt.Fprintf(os.Stderr, constants.MsgPullFailed, rec.RepoName, result.Error)
	}
}

// pullOneRepoTracked runs safe-pull with progress tracking.
func pullOneRepoTracked(rec model.ScanRecord, prog *cloner.BatchProgress) {
	if cloner.IsMissingRepo(rec.AbsolutePath) {
		prog.Skip()

		return
	}

	result := cloner.SafePullOne(rec, rec.AbsolutePath)
	if result.Success {
		prog.Succeed()
	} else {
		prog.FailWithError(rec.RepoName, result.Error)
	}
}
