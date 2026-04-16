package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/user/gitmap/cloner"
	"github.com/user/gitmap/clonenext"
	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/desktop"
	"github.com/user/gitmap/model"
	"github.com/user/gitmap/verbose"
)

// applySSHKey sets GIT_SSH_COMMAND if an SSH key name is provided.
func applySSHKey(name string) {
	if len(name) == 0 {
		return
	}

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHQuery, err)
		os.Exit(1)
	}
	defer db.Close()

	key, err := db.FindSSHKeyByName(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHNotFound, name)
		os.Exit(1)
	}

	sshCmd := fmt.Sprintf("ssh -i %s -o IdentitiesOnly=yes", key.PrivatePath)
	os.Setenv("GIT_SSH_COMMAND", sshCmd)
	fmt.Fprintf(os.Stdout, constants.MsgSSHCloneUsing, name, key.PrivatePath)
}

// runClone handles the "clone" subcommand.
func runClone(args []string) {
	checkHelp("clone", args)
	source, folderName, targetDir, sshKeyName, safePull, ghDesktop, verboseMode := parseCloneFlags(args)
	if len(source) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrSourceRequired)
		fmt.Fprintln(os.Stderr, constants.ErrCloneUsage)
		os.Exit(1)
	}
	initCloneVerbose(verboseMode)
	requireOnline()
	applySSHKey(sshKeyName)

	if isDirectURL(source) {
		executeDirectClone(source, folderName, ghDesktop)

		return
	}

	source = resolveCloneShorthand(source)
	executeClone(source, targetDir, safePull, ghDesktop)
}

// isDirectURL returns true when source is a git URL (not a file path).
func isDirectURL(source string) bool {
	lower := strings.ToLower(source)

	return strings.HasPrefix(lower, constants.PrefixHTTPS) ||
		strings.HasPrefix(lower, "http://") ||
		strings.HasPrefix(lower, constants.PrefixSSH)
}

// repoNameFromURL derives the repository name from a clone URL.
func repoNameFromURL(url string) string {
	name := strings.TrimSuffix(url, ".git")
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		name = name[idx+1:]
	}
	if idx := strings.LastIndex(name, ":"); idx >= 0 {
		name = name[idx+1:]
	}

	return name
}

// executeDirectClone clones a single repo from a direct URL.
// When no folder name is given, versioned URLs are auto-flattened
// (e.g., wp-onboarding-v13 clones into wp-onboarding/).
func executeDirectClone(url, folderName string, ghDesktopFlag bool) {
	repoName := repoNameFromURL(url)
	if len(folderName) == 0 {
		parsed := clonenext.ParseRepoName(repoName)
		if parsed.HasVersion {
			folderName = parsed.BaseName
		} else {
			folderName = repoName
		}
	}

	absPath, err := filepath.Abs(folderName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not resolve absolute path for %s: %v\n", folderName, err)
		absPath = folderName
	}

	// Check if target folder already exists.
	if _, statErr := os.Stat(absPath); statErr == nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneURLExists, absPath)
		os.Exit(1)
	}

	// Enqueue pending task.
	workDir, _ := os.Getwd()
	cmdArgs := buildCommandArgs(append([]string{"clone"}, os.Args[2:]...))
	taskID, taskDB := createPendingTask(constants.TaskTypeClone, absPath, workDir, "clone", cmdArgs)
	if taskDB != nil {
		defer taskDB.Close()
	}

	// Clone.
	fmt.Printf(constants.MsgCloneURLCloning, repoName, folderName)
	cmd := exec.Command(constants.GitBin, constants.GitClone, url, absPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cloneErr := cmd.Run()
	if cloneErr != nil {
		failPendingTask(taskDB, taskID, fmt.Sprintf(constants.ErrCloneURLFailed, url, cloneErr))
		fmt.Fprintf(os.Stderr, constants.ErrCloneURLFailed, url, cloneErr)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgCloneURLDone, repoName)

	// Upsert to database.
	upsertDirectClone(url, repoName, folderName, absPath)

	// GitHub Desktop registration (auto-register by default for direct URL).
	registerSingleDesktop(repoName, absPath)

	// Open in VS Code if available.
	openInVSCode(absPath)

	completePendingTask(taskDB, taskID)
}

// upsertDirectClone persists the cloned repo in the database.
func upsertDirectClone(url, repoName, folderName, absPath string) {
	rec := model.ScanRecord{
		Slug:         strings.ToLower(repoName),
		RepoName:     repoName,
		RelativePath: folderName,
		AbsolutePath: absPath,
	}
	if strings.HasPrefix(url, constants.PrefixSSH) {
		rec.SSHUrl = url
	} else {
		rec.HTTPSUrl = url
	}

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not open database: %v\n", err)

		return
	}
	defer db.Close()

	if upsertErr := db.UpsertRepos([]model.ScanRecord{rec}); upsertErr != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not save repo to database: %v\n", upsertErr)
	}
}

// promptOrRegisterDesktop handles GitHub Desktop registration.
func promptOrRegisterDesktop(repoName, absPath string, autoRegister bool) {
	if autoRegister {
		registerSingleDesktop(repoName, absPath)

		return
	}

	fmt.Print(constants.MsgCloneDesktopPrompt)
	var answer string
	_, _ = fmt.Scanln(&answer)
	if strings.ToLower(strings.TrimSpace(answer)) == "y" {
		registerSingleDesktop(repoName, absPath)
	}
}

// registerSingleDesktop registers a single repo with GitHub Desktop.
func registerSingleDesktop(name, absPath string) {
	records := []model.ScanRecord{{
		RepoName:     name,
		AbsolutePath: absPath,
	}}
	result := desktop.AddRepos(records)
	if result.Added > 0 {
		fmt.Printf(constants.MsgDesktopSummary, result.Added, result.Failed)
	}
}

// initCloneVerbose sets up verbose logging if enabled.
func initCloneVerbose(enabled bool) {
	if enabled {
		log, err := verbose.Init()
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.WarnVerboseLogFailed, err)

			return
		}
		defer log.Close()
	}
}

// resolveCloneShorthand maps "json", "csv", and "text" to default output paths.
func resolveCloneShorthand(source string) string {
	shorthandMap := map[string]string{
		constants.ShorthandJSON: filepath.Join(constants.DefaultOutputFolder, constants.DefaultJSONFile),
		constants.ShorthandCSV:  filepath.Join(constants.DefaultOutputFolder, constants.DefaultCSVFile),
		constants.ShorthandText: filepath.Join(constants.DefaultOutputFolder, constants.DefaultTextFile),
	}
	resolved, ok := shorthandMap[strings.ToLower(source)]
	if ok {
		return validateShorthandPath(resolved)
	}

	return source
}

// validateShorthandPath checks that the resolved shorthand file exists.
func validateShorthandPath(resolved string) string {
	_, err := os.Stat(resolved)
	if err == nil {
		return resolved
	}
	fmt.Fprintf(os.Stderr, constants.ErrShorthandNotFound, resolved)
	os.Exit(1)

	return ""
}

// executeClone runs the clone operation and prints the summary.
func executeClone(source, targetDir string, safePull, ghDesktop bool) {
	// Enqueue clone as a pending task before execution.
	absTarget, absErr := filepath.Abs(targetDir)
	if absErr != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not resolve absolute path for %s: %v\n", targetDir, absErr)
		absTarget = targetDir
	}
	workDir, wdErr := os.Getwd()
	if wdErr != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not determine working directory: %v\n", wdErr)
	}
	cmdArgs := buildCommandArgs(append([]string{"clone"}, os.Args[2:]...))
	taskID, taskDB := createPendingTask(constants.TaskTypeClone, absTarget, workDir, "clone", cmdArgs)
	if taskDB != nil {
		defer taskDB.Close()
	}

	summary, err := cloner.CloneFromFile(source, targetDir, safePull)
	if err != nil {
		failPendingTask(taskDB, taskID, fmt.Sprintf(constants.ErrCloneFailed, source, err))
		fmt.Fprintf(os.Stderr, constants.ErrCloneFailed, source, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgCloneComplete, summary.Succeeded, summary.Failed)
	printCloneFailures(summary)
	registerCloned(summary, targetDir, ghDesktop)

	// Mark clone task as completed after all steps succeed.
	completePendingTask(taskDB, taskID)
}

// printCloneFailures lists any repos that failed to clone.
func printCloneFailures(s model.CloneSummary) {
	if s.Failed == 0 {
		return
	}

	fmt.Println(constants.MsgFailedClones)
	for _, e := range s.Errors {
		fmt.Printf(constants.MsgFailedEntry,
			e.Record.RepoName, e.Record.RelativePath, e.Error)
	}
}

// registerCloned adds successfully cloned repos to GitHub Desktop.
func registerCloned(s model.CloneSummary, targetDir string, enabled bool) {
	if enabled {
		absTarget, absErr := filepath.Abs(targetDir)
		if absErr != nil {
			fmt.Fprintf(os.Stderr, "  Warning: could not resolve absolute path for %s: %v\n", targetDir, absErr)
			absTarget = targetDir
		}
		records := make([]model.ScanRecord, 0, s.Succeeded)
		for _, r := range s.Cloned {
			r.Record.AbsolutePath = filepath.Join(absTarget, r.Record.RelativePath)
			records = append(records, r.Record)
		}
		result := desktop.AddRepos(records)
		fmt.Printf(constants.MsgDesktopSummary, result.Added, result.Failed)
	}
}
