// Package main is a tiny forwarding shim placed at a drive root
// (e.g. E:\gitmap.exe) that re-executes the real gitmap binary
// living inside the deploy folder (e.g. E:\bin-run\gitmap\gitmap.exe).
//
// The target path is baked in at build time via -ldflags:
//
//	go build -ldflags "-X main.target=E:\\bin-run\\gitmap\\gitmap.exe" \
//	    -o E:\gitmap.exe ./gitmap/scripts/shim
//
// stdin/stdout/stderr are inherited; exit code is propagated.
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// target is set at build time via -ldflags "-X main.target=<path>".
var target = ""

func main() {
	if target == "" {
		fmt.Fprintln(os.Stderr, "gitmap-shim: target binary path not set at build time")
		os.Exit(2)
	}

	if _, err := os.Stat(target); err != nil {
		fmt.Fprintf(os.Stderr, "gitmap-shim: target not found: %s (%v)\n", target, err)
		os.Exit(2)
	}

	cmd := exec.Command(target, os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "gitmap-shim: failed to execute %s: %v\n", target, err)
		os.Exit(1)
	}
}
