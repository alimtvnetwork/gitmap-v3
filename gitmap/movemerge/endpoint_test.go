package movemerge

import "testing"

func TestClassifyEndpoint_Folder(t *testing.T) {
	ep := ClassifyEndpoint("./local-folder")
	if ep.Kind != KindFolder {
		t.Errorf("expected KindFolder, got %v", ep.Kind)
	}
}

func TestClassifyEndpoint_HTTPSURL_NoBranch(t *testing.T) {
	ep := ClassifyEndpoint("https://github.com/owner/repo")
	if ep.Kind != KindURL {
		t.Fatalf("expected KindURL, got %v", ep.Kind)
	}
	if ep.URL != "https://github.com/owner/repo" {
		t.Errorf("URL: %q", ep.URL)
	}
	if ep.Branch != "" {
		t.Errorf("branch: %q", ep.Branch)
	}
}

func TestClassifyEndpoint_HTTPSURL_WithBranch(t *testing.T) {
	ep := ClassifyEndpoint("https://github.com/owner/repo:develop")
	if ep.URL != "https://github.com/owner/repo" {
		t.Errorf("URL: %q", ep.URL)
	}
	if ep.Branch != "develop" {
		t.Errorf("branch: %q", ep.Branch)
	}
}

func TestClassifyEndpoint_SSHShorthand_NoBranch(t *testing.T) {
	ep := ClassifyEndpoint("git@github.com:owner/repo.git")
	if ep.Kind != KindURL {
		t.Fatalf("expected KindURL, got %v", ep.Kind)
	}
	if ep.URL != "git@github.com:owner/repo.git" {
		t.Errorf("URL: %q", ep.URL)
	}
	if ep.Branch != "" {
		t.Errorf("expected no branch, got %q", ep.Branch)
	}
}

func TestClassifyEndpoint_SSHShorthand_WithBranch(t *testing.T) {
	ep := ClassifyEndpoint("git@github.com:owner/repo.git:feature/x")
	if ep.URL != "git@github.com:owner/repo.git" {
		t.Errorf("URL: %q", ep.URL)
	}
	if ep.Branch != "feature/x" {
		t.Errorf("branch: %q", ep.Branch)
	}
}

func TestRepoNameFromURL(t *testing.T) {
	cases := map[string]string{
		"https://github.com/owner/repo":     "repo",
		"https://github.com/owner/repo.git": "repo",
		"git@github.com:owner/repo.git":     "repo",
	}
	for in, want := range cases {
		if got := RepoNameFromURL(in); got != want {
			t.Errorf("RepoNameFromURL(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestNormalizeRemote_HTTPSAndSSHEqual(t *testing.T) {
	a := normalizeRemote("https://github.com/owner/repo.git")
	b := normalizeRemote("git@github.com:owner/repo")
	if a != b {
		t.Errorf("https/ssh of same repo should normalize equal: %q vs %q", a, b)
	}
}
