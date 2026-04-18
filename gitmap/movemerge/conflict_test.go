package movemerge

import (
	"bytes"
	"strings"
	"testing"
)

func TestConflictResolver_PreferLeft(t *testing.T) {
	r := NewConflictResolver(PromptConfig{AutoMode: "left"})
	got := r.Resolve(DiffEntry{RelPath: "x"})
	if got != ResolveTakeLeft {
		t.Errorf("expected ResolveTakeLeft, got %v", got)
	}
}

func TestConflictResolver_PreferNewer_RightNewer(t *testing.T) {
	r := NewConflictResolver(PromptConfig{AutoMode: "newer"})
	got := r.Resolve(DiffEntry{LeftMTime: 100, RightMTime: 200})
	if got != ResolveTakeRight {
		t.Errorf("right newer -> take right; got %v", got)
	}
}

func TestConflictResolver_AllLeftIsSticky(t *testing.T) {
	in := strings.NewReader("A\n")
	out := &bytes.Buffer{}
	r := NewConflictResolver(PromptConfig{In: in, Out: out})

	first := r.Resolve(DiffEntry{RelPath: "a"})
	if first != ResolveTakeLeft {
		t.Errorf("first call after A: want ResolveTakeLeft, got %v", first)
	}
	second := r.Resolve(DiffEntry{RelPath: "b"})
	if second != ResolveTakeLeft {
		t.Errorf("sticky All-Left should apply to second conflict; got %v", second)
	}
}

func TestConflictResolver_QuitAborts(t *testing.T) {
	in := strings.NewReader("Q\n")
	r := NewConflictResolver(PromptConfig{In: in, Out: &bytes.Buffer{}})
	if got := r.Resolve(DiffEntry{}); got != ResolveQuit {
		t.Errorf("Q should ResolveQuit; got %v", got)
	}
}
