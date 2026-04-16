// Package model defines the core data structures for gitmap.
package model

// ReleaseRecord holds release metadata stored in the database.
type ReleaseRecord struct {
	ID           int64  `json:"id"`
	Version      string `json:"version"`
	Tag          string `json:"tag"`
	Branch       string `json:"branch"`
	SourceBranch string `json:"sourceBranch"`
	CommitSha    string `json:"commitSha"`
	Changelog    string `json:"changelog"`
	Notes        string `json:"notes"`
	Draft        bool   `json:"draft"`
	PreRelease   bool   `json:"preRelease"`
	IsLatest     bool   `json:"isLatest"`
	Source       string `json:"source"`
	CreatedAt    string `json:"createdAt"`
}

// Release source values.
const (
	SourceRelease = "release"
	SourceImport  = "import"
	SourceRepo    = "repo"
	SourceTag     = "tag"
)
