package movemerge

import (
	"fmt"
	"path/filepath"

	"github.com/user/gitmap/constants"
)

// mergeFlowConfig parameterises runMergeFlow per command.
type mergeFlowConfig struct {
	WriteLeft   bool
	WriteRight  bool
	DefaultAuto string // applied when AutoMode is "" but bypass requested
}

// runMergeFlow drives the diff -> resolve -> apply loop shared by
// merge-left / merge-right / merge-both.
func runMergeFlow(left, right Endpoint, opts Options, logger *Logger, cfg mergeFlowConfig) error {
	copyOpts := CopyOptions{
		IncludeVCS: opts.IncludeVCS, IncludeNodeModules: opts.IncludeNodeMod,
		DryRun: opts.DryRun,
	}
	logger.Logf("diffing trees ...")
	entries, err := DiffTrees(left.WorkingDir, right.WorkingDir, copyOpts)
	if err != nil {
		return fmt.Errorf("diff: %w", err)
	}

	resolver := NewConflictResolver(PromptConfig{AutoMode: opts.AutoMode})

	return applyDiffEntries(entries, left, right, cfg, copyOpts, resolver, logger)
}

// applyDiffEntries iterates DiffEntries and applies them.
func applyDiffEntries(entries []DiffEntry, left, right Endpoint, cfg mergeFlowConfig,
	copyOpts CopyOptions, resolver *ConflictResolver, logger *Logger) error {
	for _, entry := range entries {
		if entry.Kind == DiffIdentical {
			continue
		}
		if err := applyOneEntry(entry, left, right, cfg, copyOpts, resolver, logger); err != nil {
			return err
		}
	}

	return nil
}

// applyOneEntry handles a single DiffEntry.
func applyOneEntry(entry DiffEntry, left, right Endpoint, cfg mergeFlowConfig,
	copyOpts CopyOptions, resolver *ConflictResolver, logger *Logger) error {
	switch entry.Kind {
	case DiffMissingRight:
		if !cfg.WriteRight {
			return nil
		}

		return copyOneFile(left.WorkingDir, right.WorkingDir, entry.RelPath, copyOpts, logger, "RIGHT")
	case DiffMissingLeft:
		if !cfg.WriteLeft {
			return nil
		}

		return copyOneFile(right.WorkingDir, left.WorkingDir, entry.RelPath, copyOpts, logger, "LEFT")
	case DiffConflict:
		return resolveConflict(entry, left, right, cfg, copyOpts, resolver, logger)
	}

	return nil
}

// resolveConflict consults the resolver and applies the choice.
func resolveConflict(entry DiffEntry, left, right Endpoint, cfg mergeFlowConfig,
	copyOpts CopyOptions, resolver *ConflictResolver, logger *Logger) error {
	choice := resolver.Resolve(entry)
	switch choice {
	case ResolveQuit:
		logger.Logf("user quit; stopping")

		return fmt.Errorf("user aborted at %s", entry.RelPath)
	case ResolveSkip:
		logger.Logf(constants.ConflictSkippedFmt, entry.RelPath)

		return nil
	case ResolveTakeLeft:
		if !cfg.WriteRight {
			return nil
		}
		logger.Logf(constants.ConflictAppliedFmt, entry.RelPath, "LEFT")

		return copyOneFile(left.WorkingDir, right.WorkingDir, entry.RelPath, copyOpts, logger, "RIGHT")
	case ResolveTakeRight:
		if !cfg.WriteLeft {
			return nil
		}
		logger.Logf(constants.ConflictAppliedFmt, entry.RelPath, "RIGHT")

		return copyOneFile(right.WorkingDir, left.WorkingDir, entry.RelPath, copyOpts, logger, "LEFT")
	}

	return nil
}

// copyOneFile copies a single relative path from srcRoot to dstRoot.
func copyOneFile(srcRoot, dstRoot, rel string, opts CopyOptions, logger *Logger, sideLabel string) error {
	if opts.DryRun {
		logger.Logf("  [dry-run] copy %s -> %s", rel, sideLabel)

		return nil
	}
	src := filepath.Join(srcRoot, rel)
	dst := filepath.Join(dstRoot, rel)
	info, err := osStat(src)
	if err != nil {
		return fmt.Errorf("stat %s: %w", src, err)
	}
	if info.IsDir() {
		return mkAllDir(dst, info.Mode().Perm())
	}
	if _, err = copyFileWithMode(src, dst, info.Mode().Perm()); err != nil {
		return fmt.Errorf("copy %s: %w", rel, err)
	}

	return nil
}
