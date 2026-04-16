package cmd

import (
	"testing"
)

func TestParseStatusFlags_NoFlags(t *testing.T) {
	group, all := parseStatusFlags([]string{})
	if len(group) > 0 {
		t.Errorf("expected empty group, got %q", group)
	}
	if all {
		t.Error("expected all=false")
	}
}

func TestParseStatusFlags_GroupLong(t *testing.T) {
	group, all := parseStatusFlags([]string{"--group", "backend"})
	if group != "backend" {
		t.Errorf("expected group=backend, got %q", group)
	}
	if all {
		t.Error("expected all=false")
	}
}

func TestParseStatusFlags_GroupShort(t *testing.T) {
	group, all := parseStatusFlags([]string{"-g", "frontend"})
	if group != "frontend" {
		t.Errorf("expected group=frontend, got %q", group)
	}
	if all {
		t.Error("expected all=false")
	}
}

func TestParseStatusFlags_All(t *testing.T) {
	group, all := parseStatusFlags([]string{"--all"})
	if len(group) > 0 {
		t.Errorf("expected empty group, got %q", group)
	}
	if all != true {
		t.Error("expected all=true")
	}
}

func TestParseStatusFlags_GroupAndAll(t *testing.T) {
	group, all := parseStatusFlags([]string{"--group", "ops", "--all"})
	if group != "ops" {
		t.Errorf("expected group=ops, got %q", group)
	}
	if all != true {
		t.Error("expected all=true")
	}
}

func TestParseExecFlags_NoFlags(t *testing.T) {
	group, all, _, gitArgs := parseExecFlags([]string{"fetch", "--prune"})
	if len(group) > 0 {
		t.Errorf("expected empty group, got %q", group)
	}
	if all {
		t.Error("expected all=false")
	}
	if len(gitArgs) != 2 || gitArgs[0] != "fetch" || gitArgs[1] != "--prune" {
		t.Errorf("expected [fetch --prune], got %v", gitArgs)
	}
}

func TestParseExecFlags_GroupLong(t *testing.T) {
	group, all, _, gitArgs := parseExecFlags([]string{"--group", "backend", "status"})
	if group != "backend" {
		t.Errorf("expected group=backend, got %q", group)
	}
	if all {
		t.Error("expected all=false")
	}
	if len(gitArgs) != 1 || gitArgs[0] != "status" {
		t.Errorf("expected [status], got %v", gitArgs)
	}
}

func TestParseExecFlags_GroupShort(t *testing.T) {
	group, _, _, gitArgs := parseExecFlags([]string{"-g", "infra", "pull"})
	if group != "infra" {
		t.Errorf("expected group=infra, got %q", group)
	}
	if len(gitArgs) != 1 || gitArgs[0] != "pull" {
		t.Errorf("expected [pull], got %v", gitArgs)
	}
}

func TestParseExecFlags_All(t *testing.T) {
	group, all, _, gitArgs := parseExecFlags([]string{"--all", "fetch"})
	if len(group) > 0 {
		t.Errorf("expected empty group, got %q", group)
	}
	if all != true {
		t.Error("expected all=true")
	}
	if len(gitArgs) != 1 || gitArgs[0] != "fetch" {
		t.Errorf("expected [fetch], got %v", gitArgs)
	}
}

func TestParseExecFlags_NoArgs(t *testing.T) {
	group, all, _, gitArgs := parseExecFlags([]string{"--all"})
	if len(group) > 0 {
		t.Errorf("expected empty group, got %q", group)
	}
	if all != true {
		t.Error("expected all=true")
	}
	if len(gitArgs) != 0 {
		t.Errorf("expected empty gitArgs, got %v", gitArgs)
	}
}

func TestParsePullFlags_NoFlags(t *testing.T) {
	slug, group, all, verbose, _ := parsePullFlags([]string{"my-repo"})
	if slug != "my-repo" {
		t.Errorf("expected slug=my-repo, got %q", slug)
	}
	if len(group) > 0 || all || verbose {
		t.Error("expected no group/all/verbose")
	}
}

func TestParsePullFlags_GroupLong(t *testing.T) {
	slug, group, all, _, _ := parsePullFlags([]string{"--group", "backend"})
	if len(slug) > 0 {
		t.Errorf("expected empty slug, got %q", slug)
	}
	if group != "backend" {
		t.Errorf("expected group=backend, got %q", group)
	}
	if all {
		t.Error("expected all=false")
	}
}

func TestParsePullFlags_GroupShort(t *testing.T) {
	_, group, _, _, _ := parsePullFlags([]string{"-g", "infra"})
	if group != "infra" {
		t.Errorf("expected group=infra, got %q", group)
	}
}

func TestParsePullFlags_All(t *testing.T) {
	slug, group, all, _, _ := parsePullFlags([]string{"--all"})
	if len(slug) > 0 || len(group) > 0 {
		t.Error("expected empty slug and group")
	}
	if all != true {
		t.Error("expected all=true")
	}
}

func TestParsePullFlags_AllWithVerbose(t *testing.T) {
	_, _, all, verbose, _ := parsePullFlags([]string{"--all", "--verbose"})
	if all != true {
		t.Error("expected all=true")
	}
	if verbose != true {
		t.Error("expected verbose=true")
	}
}
