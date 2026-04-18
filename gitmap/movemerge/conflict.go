package movemerge

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/user/gitmap/constants"
)

// ConflictResolution describes what to do for one conflict.
type ConflictResolution int

const (
	// ResolveTakeLeft = overwrite RIGHT with LEFT.
	ResolveTakeLeft ConflictResolution = iota
	// ResolveTakeRight = overwrite LEFT with RIGHT.
	ResolveTakeRight
	// ResolveSkip = leave both sides as-is.
	ResolveSkip
	// ResolveQuit = abort; caller stops processing further conflicts.
	ResolveQuit
)

// PromptConfig controls the conflict prompt loop.
type PromptConfig struct {
	// AutoMode = "" (interactive) | "left" | "right" | "newer" | "skip".
	AutoMode string
	// In/Out are wired for testing; default to stdin/stderr.
	In  io.Reader
	Out io.Writer
}

// ConflictResolver returns a resolution per DiffEntry. It applies
// `--prefer-*` / `-y` bypass when AutoMode is set; otherwise it
// prompts the user once per conflict (with sticky All-Left/All-Right
// when the user picks A or B).
type ConflictResolver struct {
	cfg     PromptConfig
	scanner *bufio.Scanner
	sticky  ConflictResolution
	hasSticky bool
}

// NewConflictResolver constructs a resolver with sane defaults.
func NewConflictResolver(cfg PromptConfig) *ConflictResolver {
	if cfg.In == nil {
		cfg.In = os.Stdin
	}
	if cfg.Out == nil {
		cfg.Out = os.Stderr
	}

	return &ConflictResolver{
		cfg:     cfg,
		scanner: bufio.NewScanner(cfg.In),
	}
}

// Resolve returns the resolution for one conflict entry.
func (r *ConflictResolver) Resolve(entry DiffEntry) ConflictResolution {
	if r.hasSticky {
		return r.sticky
	}
	if len(r.cfg.AutoMode) > 0 {
		return resolveByPolicy(r.cfg.AutoMode, entry)
	}

	return r.promptOnce(entry)
}

// resolveByPolicy applies a non-interactive --prefer-* policy.
func resolveByPolicy(mode string, entry DiffEntry) ConflictResolution {
	switch mode {
	case constants.PreferLeft:
		return ResolveTakeLeft
	case constants.PreferRight:
		return ResolveTakeRight
	case constants.PreferSkip:
		return ResolveSkip
	case constants.PreferNewer:
		if entry.RightMTime > entry.LeftMTime {
			return ResolveTakeRight
		}

		return ResolveTakeLeft
	}

	return ResolveSkip
}

// promptOnce shows the conflict prompt and reads a single keystroke.
func (r *ConflictResolver) promptOnce(entry DiffEntry) ConflictResolution {
	fmt.Fprintf(r.cfg.Out, constants.ConflictPromptHeaderFmt, entry.RelPath)
	fmt.Fprintf(r.cfg.Out, constants.ConflictPromptMetaFmt,
		"LEFT", humanSize(entry.LeftSize), humanTime(entry.LeftMTime))
	fmt.Fprintf(r.cfg.Out, constants.ConflictPromptMetaFmt,
		"RIGHT", humanSize(entry.RightSize), humanTime(entry.RightMTime))
	fmt.Fprint(r.cfg.Out, constants.ConflictPromptKeysLine)

	if !r.scanner.Scan() {
		return ResolveQuit
	}

	return r.applyKey(strings.ToUpper(strings.TrimSpace(r.scanner.Text())))
}

// applyKey maps a keystroke to a ConflictResolution, handling sticky
// All-Left / All-Right modes.
func (r *ConflictResolver) applyKey(key string) ConflictResolution {
	switch key {
	case constants.ConflictKeyLeft:
		return ResolveTakeLeft
	case constants.ConflictKeyRight:
		return ResolveTakeRight
	case constants.ConflictKeySkip:
		return ResolveSkip
	case constants.ConflictKeyAllLeft:
		r.sticky, r.hasSticky = ResolveTakeLeft, true

		return ResolveTakeLeft
	case constants.ConflictKeyAllRight:
		r.sticky, r.hasSticky = ResolveTakeRight, true

		return ResolveTakeRight
	case constants.ConflictKeyQuit:
		return ResolveQuit
	}

	return ResolveSkip
}

// humanSize formats bytes as a short string.
func humanSize(b int64) string {
	const k = 1024
	if b < k {
		return fmt.Sprintf("%d B", b)
	}
	if b < k*k {
		return fmt.Sprintf("%.1f KB", float64(b)/float64(k))
	}

	return fmt.Sprintf("%.1f MB", float64(b)/float64(k*k))
}

// humanTime renders a Unix second as "YYYY-MM-DD HH:MM".
func humanTime(unix int64) string {
	if unix == 0 {
		return "-"
	}

	return time.Unix(unix, 0).Format("2006-01-02 15:04")
}
