package movemerge

import (
	"fmt"
	"io"
	"os"
)

// Logger emits structured `[mv]` / `[merge-*]` prefixed lines to a
// writer (default os.Stderr). Quiet mode silences everything.
type Logger struct {
	prefix string
	out    io.Writer
	quiet  bool
}

// NewLogger constructs a Logger for one operation.
func NewLogger(prefix string) *Logger {
	return &Logger{prefix: prefix, out: os.Stderr}
}

// SetOutput replaces the writer (used in tests).
func (l *Logger) SetOutput(w io.Writer) {
	if w != nil {
		l.out = w
	}
}

// SetQuiet silences all output.
func (l *Logger) SetQuiet(quiet bool) {
	l.quiet = quiet
}

// Logf prints one prefixed line with Sprintf-style formatting.
func (l *Logger) Logf(format string, args ...interface{}) {
	if l.quiet {
		return
	}
	fmt.Fprintf(l.out, "%s %s\n", l.prefix, fmt.Sprintf(format, args...))
}
