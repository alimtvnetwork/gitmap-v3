// Package store — csharpmetadata.go manages CSharp metadata + files.
package store

import (
	"fmt"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
)

// UpsertCSharpMetadata inserts or updates C# metadata for a detected project.
func (db *DB) UpsertCSharpMetadata(m model.CSharpProjectMetadata) error {
	_, err := db.conn.Exec(constants.SQLUpsertCSharpMetadata,
		m.DetectedProjectID, m.SlnPath, m.SlnName,
		m.GlobalJsonPath, m.SdkVersion)

	return err
}

// UpsertCSharpProjectFile inserts or updates a C# project file record.
func (db *DB) UpsertCSharpProjectFile(f model.CSharpProjectFile) error {
	_, err := db.conn.Exec(constants.SQLUpsertCSharpProjectFile,
		f.CSharpMetadataID, f.FilePath, f.RelativePath,
		f.FileName, f.ProjectName, f.TargetFramework, f.OutputType, f.Sdk)

	return err
}

// UpsertCSharpKeyFile inserts or updates a C# key file record.
func (db *DB) UpsertCSharpKeyFile(f model.CSharpKeyFile) error {
	_, err := db.conn.Exec(constants.SQLUpsertCSharpKeyFile,
		f.CSharpMetadataID, f.FileType, f.FilePath, f.RelativePath)

	return err
}

// SelectCSharpMetadata returns C# metadata for a detected project.
func (db *DB) SelectCSharpMetadata(detectedProjectID int64) (*model.CSharpProjectMetadata, error) {
	var m model.CSharpProjectMetadata
	err := db.conn.QueryRow(constants.SQLSelectCSharpMetadata, detectedProjectID).Scan(
		&m.ID, &m.DetectedProjectID, &m.SlnPath, &m.SlnName,
		&m.GlobalJsonPath, &m.SdkVersion)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

// SelectCSharpProjectFiles returns all .csproj files for a metadata ID.
func (db *DB) SelectCSharpProjectFiles(metadataID int64) ([]model.CSharpProjectFile, error) {
	rows, err := db.conn.Query(constants.SQLSelectCSharpProjectFiles, metadataID)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrProjectQuery, err)
	}
	defer rows.Close()

	return scanCSharpFileRows(rows)
}

// SelectCSharpKeyFiles returns all key files for a metadata ID.
func (db *DB) SelectCSharpKeyFiles(metadataID int64) ([]model.CSharpKeyFile, error) {
	rows, err := db.conn.Query(constants.SQLSelectCSharpKeyFiles, metadataID)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrProjectQuery, err)
	}
	defer rows.Close()

	return scanCSharpKeyFileRows(rows)
}

// DeleteStaleCSharpFiles removes project files not in the keep list.
func (db *DB) DeleteStaleCSharpFiles(metadataID int64, keepIDs []int64) error {
	if len(keepIDs) == 0 {
		return nil
	}
	placeholders := buildPlaceholders(len(keepIDs))
	query := fmt.Sprintf(constants.SQLDeleteStaleCSharpFiles, placeholders)
	args := buildStaleArgsInt64(metadataID, keepIDs)
	_, err := db.conn.Exec(query, args...)

	return err
}

// DeleteStaleCSharpKeyFiles removes key files not in the keep list.
func (db *DB) DeleteStaleCSharpKeyFiles(metadataID int64, keepIDs []int64) error {
	if len(keepIDs) == 0 {
		return nil
	}
	placeholders := buildPlaceholders(len(keepIDs))
	query := fmt.Sprintf(constants.SQLDeleteStaleCSharpKeyFiles, placeholders)
	args := buildStaleArgsInt64(metadataID, keepIDs)
	_, err := db.conn.Exec(query, args...)

	return err
}

// scanCSharpFileRows scans rows into CSharpProjectFile slices.
func scanCSharpFileRows(rows interface{ Next() bool; Scan(...interface{}) error }) ([]model.CSharpProjectFile, error) {
	var files []model.CSharpProjectFile
	for rows.Next() {
		var f model.CSharpProjectFile
		err := rows.Scan(&f.ID, &f.CSharpMetadataID, &f.FilePath,
			&f.RelativePath, &f.FileName, &f.ProjectName,
			&f.TargetFramework, &f.OutputType, &f.Sdk)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	return files, nil
}

// scanCSharpKeyFileRows scans rows into CSharpKeyFile slices.
func scanCSharpKeyFileRows(rows interface{ Next() bool; Scan(...interface{}) error }) ([]model.CSharpKeyFile, error) {
	var files []model.CSharpKeyFile
	for rows.Next() {
		var f model.CSharpKeyFile
		err := rows.Scan(&f.ID, &f.CSharpMetadataID, &f.FileType,
			&f.FilePath, &f.RelativePath)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	return files, nil
}
