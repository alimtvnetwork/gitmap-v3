package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/gitmap/constants"
)

// runUpdateCleanup handles the "update-cleanup" subcommand.
// Removes leftover temp binaries and .old backup files.
func runUpdateCleanup() {
	fmt.Println(constants.MsgUpdateCleanStart)

	tempCleaned := cleanupTempCopies()
	oldCleaned := cleanupOldBackups()

	total := tempCleaned + oldCleaned
	if total > 0 {
		fmt.Printf(constants.MsgUpdateCleanDone, total)
	} else {
		fmt.Println(constants.MsgUpdateCleanNone)
	}
}

// cleanupTempCopies removes leftover handoff binaries from previous updates.
func cleanupTempCopies() int {
	selfPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not determine executable path: %v\n", err)

		return 0
	}
	patterns := buildTempPatterns(selfPath)

	return removeMatchingFiles(patterns, selfPath)
}

// buildTempPatterns returns glob patterns for temp update copies.
func buildTempPatterns(selfPath string) []string {
	patterns := []string{
		filepath.Join(os.TempDir(), constants.UpdateCopyGlob),
	}
	if len(selfPath) > 0 {
		patterns = append(patterns, filepath.Join(filepath.Dir(selfPath), constants.UpdateCopyGlob))
	}

	return patterns
}

// removeMatchingFiles removes files matching patterns, skipping selfPath.
func removeMatchingFiles(patterns []string, selfPath string) int {
	seen := map[string]bool{}
	cleaned := 0
	for _, pattern := range patterns {
		cleaned += removeGlobMatches(pattern, selfPath, seen)
	}

	return cleaned
}

// removeGlobMatches removes files matching a single glob pattern.
func removeGlobMatches(pattern, selfPath string, seen map[string]bool) int {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return 0
	}

	cleaned := 0
	for _, match := range matches {
		if removeTempMatch(match, selfPath, seen) {
			cleaned++
		}
	}

	return cleaned
}

// removeTempMatch removes a single temp file if not already seen or self.
func removeTempMatch(match, selfPath string, seen map[string]bool) bool {
	cleanPath := filepath.Clean(match)
	if seen[cleanPath] {
		return false
	}
	seen[cleanPath] = true

	if len(selfPath) > 0 && cleanPath == filepath.Clean(selfPath) {
		return false
	}
	if os.Remove(match) == nil {
		fmt.Printf(constants.MsgUpdateTempRemoved, filepath.Base(match))

		return true
	}

	return false
}

// cleanupOldBackups removes .old backup binaries from the deploy directory.
func cleanupOldBackups() int {
	repoPath := constants.RepoPath
	if len(repoPath) == 0 {
		return 0
	}

	deployPath := readDeployPath(repoPath)
	if len(deployPath) == 0 {
		return 0
	}

	return removeOldFiles(deployPath)
}

// readDeployPath reads the deploy path from powershell.json.
func readDeployPath(repoPath string) string {
	configPath := filepath.Join(repoPath, constants.GitMapSubdir, constants.PowershellConfigFile)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return ""
	}

	return extractJSONString(data, constants.JSONKeyDeployPath)
}

// removeOldFiles removes .old files from the deploy app directory.
func removeOldFiles(deployPath string) int {
	appDir := filepath.Join(deployPath, "gitmap")
	pattern := filepath.Join(appDir, constants.OldBackupGlob)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return 0
	}

	cleaned := 0
	for _, match := range matches {
		if os.Remove(match) == nil {
			fmt.Printf(constants.MsgUpdateOldRemoved, filepath.Base(match))
			cleaned++
		}
	}

	return cleaned
}

// extractJSONString extracts a string value from JSON bytes by key.
func extractJSONString(data []byte, key string) string {
	s := string(data)
	needle := `"` + key + `"`
	idx := findKeyValue(s, needle)
	if idx < 0 {
		return ""
	}

	return extractQuotedValue(s, idx)
}

// findKeyValue finds the position after a JSON key and colon.
func findKeyValue(s, needle string) int {
	idx := indexOf(s, needle)
	if idx < 0 {
		return -1
	}
	idx += len(needle)

	for idx < len(s) && (s[idx] == ' ' || s[idx] == ':' || s[idx] == '\t') {
		idx++
	}

	return idx
}

// extractQuotedValue extracts a quoted string starting at idx.
func extractQuotedValue(s string, idx int) string {
	if idx >= len(s) || s[idx] != '"' {
		return ""
	}
	end := indexOf(s[idx+1:], `"`)
	if end >= 0 {
		return s[idx+1 : idx+1+end]
	}

	return ""
}

// indexOf returns the index of substr in s, or -1.
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}

	return -1
}
