package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
)

// releaseResponse is a minimal GitHub release API response.
type releaseResponse struct {
	TagName string `json:"tag_name"`
}

// fetchLatestTag queries the GitHub releases API for the latest tag.
func fetchLatestTag() (string, error) {
	req, err := http.NewRequest("GET", GitHubAPILatest, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "gitmap-updater/"+Version)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned HTTP %d", resp.StatusCode)
	}

	var release releaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return release.TagName, nil
}

// getInstalledVersion runs `gitmap version` and returns the output.
func getInstalledVersion() (string, error) {
	cmd := exec.Command(GitMapBin, "version")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// normalizeVersion strips the "v" prefix for comparison.
func normalizeVersion(v string) string {
	return strings.TrimPrefix(strings.TrimSpace(v), "v")
}
