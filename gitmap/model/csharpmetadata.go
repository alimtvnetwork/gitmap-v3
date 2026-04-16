// Package model — csharpmetadata.go defines C#-specific metadata structs.
package model

// CSharpProjectMetadata holds C#-specific metadata for a detected project.
type CSharpProjectMetadata struct {
	ID                int64               `json:"id"`
	DetectedProjectID int64               `json:"detectedProjectId"`
	SlnPath           string              `json:"slnPath"`
	SlnName           string              `json:"slnName"`
	GlobalJsonPath    string              `json:"globalJsonPath"`
	SdkVersion        string              `json:"sdkVersion"`
	ProjectFiles      []CSharpProjectFile `json:"projectFiles"`
	KeyFiles          []CSharpKeyFile     `json:"keyFiles"`
}

// CSharpProjectFile represents a .csproj or .fsproj discovered in a C# project.
type CSharpProjectFile struct {
	ID               int64  `json:"id"`
	CSharpMetadataID int64  `json:"csharpMetadataId"`
	FilePath         string `json:"filePath"`
	RelativePath     string `json:"relativePath"`
	FileName         string `json:"fileName"`
	ProjectName      string `json:"projectName"`
	TargetFramework  string `json:"targetFramework"`
	OutputType       string `json:"outputType"`
	Sdk              string `json:"sdk"`
}

// CSharpKeyFile represents a key configuration file in a C# project.
type CSharpKeyFile struct {
	ID               int64  `json:"id"`
	CSharpMetadataID int64  `json:"csharpMetadataId"`
	FileType         string `json:"fileType"`
	FilePath         string `json:"filePath"`
	RelativePath     string `json:"relativePath"`
}

// CSharpProjectRecord combines a DetectedProject with its C# metadata for JSON output.
type CSharpProjectRecord struct {
	DetectedProject
	CSharpMetadata *CSharpProjectMetadata `json:"csharpMetadata,omitempty"`
}
