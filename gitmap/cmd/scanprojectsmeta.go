// Package cmd — scanprojectsmeta.go handles Go and C# metadata persistence.
package cmd

import (
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/detector"
	"github.com/user/gitmap/model"
	"github.com/user/gitmap/store"
)

// upsertGoProjectMeta persists Go metadata and runnables.
func upsertGoProjectMeta(db *store.DB, r detector.DetectionResult) {
	r.GoMeta.DetectedProjectID = r.Project.ID
	if err := db.UpsertGoMetadata(*r.GoMeta); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrGoMetadataUpsert, err)

		return
	}
	saved, err := db.SelectGoMetadata(r.Project.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrGoMetadataUpsert, err)

		return
	}
	r.GoMeta.ID = saved.ID
	runnableIDs := upsertGoRunnables(db, r.GoMeta)
	if err := db.DeleteStaleGoRunnables(r.GoMeta.ID, runnableIDs); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not clean stale Go runnables: %v\n", err)
	}
}

// upsertGoRunnables persists all runnable files and returns their IDs.
func upsertGoRunnables(db *store.DB, meta *model.GoProjectMetadata) []int64 {
	var ids []int64
	for _, run := range meta.Runnables {
		run.GoMetadataID = meta.ID
		if err := db.UpsertGoRunnable(run); err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrGoRunnableUpsert, err)

			continue
		}
		ids = append(ids, run.ID)
	}

	return ids
}

// upsertCSharpProjectMeta persists C# metadata, project files, and key files.
func upsertCSharpProjectMeta(db *store.DB, r detector.DetectionResult) {
	r.CSharp.DetectedProjectID = r.Project.ID
	if err := db.UpsertCSharpMetadata(*r.CSharp); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCSharpMetaUpsert, err)

		return
	}
	saved, err := db.SelectCSharpMetadata(r.Project.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCSharpMetaUpsert, err)

		return
	}
	r.CSharp.ID = saved.ID
	fileIDs := upsertCSharpFiles(db, r.CSharp)
	keyIDs := upsertCSharpKeyFiles(db, r.CSharp)
	if err := db.DeleteStaleCSharpFiles(r.CSharp.ID, fileIDs); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not clean stale C# files: %v\n", err)
	}
	if err := db.DeleteStaleCSharpKeyFiles(r.CSharp.ID, keyIDs); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not clean stale C# key files: %v\n", err)
	}
}

// upsertCSharpFiles persists C# project files and returns their IDs.
func upsertCSharpFiles(db *store.DB, meta *model.CSharpProjectMetadata) []int64 {
	var ids []int64
	for _, f := range meta.ProjectFiles {
		f.CSharpMetadataID = meta.ID
		if err := db.UpsertCSharpProjectFile(f); err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrCSharpFileUpsert, err)

			continue
		}
		ids = append(ids, f.ID)
	}

	return ids
}

// upsertCSharpKeyFiles persists C# key files and returns their IDs.
func upsertCSharpKeyFiles(db *store.DB, meta *model.CSharpProjectMetadata) []int64 {
	var ids []int64
	for _, f := range meta.KeyFiles {
		f.CSharpMetadataID = meta.ID
		if err := db.UpsertCSharpKeyFile(f); err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrCSharpKeyUpsert, err)

			continue
		}
		ids = append(ids, f.ID)
	}

	return ids
}

// collectRepoIDs extracts unique repo IDs from detection results.
func collectRepoIDs(results []detector.DetectionResult) map[int64]bool {
	ids := make(map[int64]bool)
	for _, r := range results {
		ids[r.Project.RepoID] = true
	}

	return ids
}

// cleanStaleProjects removes projects no longer detected for each repo.
func cleanStaleProjects(db *store.DB, repoIDs map[int64]bool, results []detector.DetectionResult) {
	for repoID := range repoIDs {
		keepIDs := collectKeepIDs(repoID, results)
		cleaned, err := db.DeleteStaleProjects(repoID, keepIDs)
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrProjectCleanup, repoID, err)

			continue
		}
		if cleaned > 0 {
			fmt.Printf(constants.MsgProjectCleanedStale, cleaned)
		}
	}
}

// collectKeepIDs collects project IDs to keep for a given repo.
func collectKeepIDs(repoID int64, results []detector.DetectionResult) []int64 {
	var ids []int64
	for _, r := range results {
		if r.Project.RepoID == repoID {
			ids = append(ids, r.Project.ID)
		}
	}

	return ids
}
