# Spec — Table of Contents

This folder contains all specification documents, issue post-mortems, design guidelines, and the generic CLI blueprint for gitmap.

---

## 01-app/ — Application Specifications

Feature specs, command designs, and refactor documentation for the gitmap CLI.

| # | File | Topic |
|---|------|-------|
| 01 | [overview.md](01-app/01-overview.md) | Project overview and architecture |
| 02 | [cli-interface.md](01-app/02-cli-interface.md) | CLI interface design |
| 03 | [scanner.md](01-app/03-scanner.md) | Directory scanner |
| 04 | [formatter.md](01-app/04-formatter.md) | Output formatting (CSV/JSON/text) |
| 05 | [cloner.md](01-app/05-cloner.md) | Repository cloner |
| 06 | [config.md](01-app/06-config.md) | Configuration system |
| 07 | [data-model.md](01-app/07-data-model.md) | Data model |
| 08 | [acceptance-criteria.md](01-app/08-acceptance-criteria.md) | Acceptance criteria |
| 09 | [build-deploy.md](01-app/09-build-deploy.md) | Build and deploy pipeline |
| 10 | [github-desktop.md](01-app/10-github-desktop.md) | GitHub Desktop integration |
| 11 | [desktop-sync.md](01-app/11-desktop-sync.md) | Desktop sync command |
| 12 | [release-command.md](01-app/12-release-command.md) | Release command |
| 13 | [release-data-model.md](01-app/13-release-data-model.md) | Release data model |
| 14 | [latest-branch.md](01-app/14-latest-branch.md) | Latest branch detection |
| 15 | [date-display-format.md](01-app/15-date-display-format.md) | Date display formatting |
| 16 | [database.md](01-app/16-database.md) | SQLite database design |
| 17 | [repo-grouping.md](01-app/17-repo-grouping.md) | Repository grouping |
| 18 | [compliance-audit.md](01-app/18-compliance-audit.md) | Compliance audit |
| 19 | [list-versions.md](01-app/19-list-versions.md) | List versions command |
| 20 | [revert.md](01-app/20-revert.md) | Revert command |
| 21 | [list-releases.md](01-app/21-list-releases.md) | List releases command |
| 22 | [scan-release-import.md](01-app/22-scan-release-import.md) | Scan release import |
| 23 | [seo-write.md](01-app/23-seo-write.md) | SEO write command |
| 24 | [amend-author.md](01-app/24-amend-author.md) | Amend author command |
| 25 | [command-history.md](01-app/25-command-history.md) | Command history |
| 26 | [stats.md](01-app/26-stats.md) | Repository statistics |
| 27 | [bookmarks.md](01-app/27-bookmarks.md) | Bookmarked commands |
| 28 | [export.md](01-app/28-export.md) | Export command |
| 29 | [import.md](01-app/29-import.md) | Import command |
| 30 | [profiles.md](01-app/30-profiles.md) | Scan profiles |
| 31 | [cd.md](01-app/31-cd.md) | Navigate to repo |
| 32 | [watch.md](01-app/32-watch.md) | Watch for changes |
| 33 | [diff-profiles.md](01-app/33-diff-profiles.md) | Diff profiles |
| 34 | [clone-progress.md](01-app/34-clone-progress.md) | Clone progress bar |
| 35 | [docs-site.md](01-app/35-docs-site.md) | Documentation site |
| 36 | [gomod-rename.md](01-app/36-gomod-rename.md) | Go module rename |
| 37 | [project-detection/](01-app/37-project-detection/) | Project type detection |
| 38 | [command-help.md](01-app/38-command-help.md) | Command help system |
| 39 | [shell-completion.md](01-app/39-shell-completion.md) | Shell completion |
| 40 | [enhanced-groups-and-listing.md](01-app/40-enhanced-groups-and-listing.md) | Enhanced groups and listing |
| 41 | [go-release-assets.md](01-app/41-go-release-assets.md) | Go release assets pipeline |
| 42 | [cross-platform.md](01-app/42-cross-platform.md) | Cross-platform support |
| 43 | [interactive-tui.md](01-app/43-interactive-tui.md) | Interactive TUI |
| 44 | [list-db-diagnostic.md](01-app/44-list-db-diagnostic.md) | List DB diagnostic |
| 45 | [release-pending-metadata.md](01-app/45-release-pending-metadata.md) | Release pending metadata |
| 46 | [clear-release-json.md](01-app/46-clear-release-json.md) | Clear release JSON |
| 47 | [zip-groups.md](01-app/47-zip-groups.md) | Zip groups |
| 48 | [repo-aliases.md](01-app/48-repo-aliases.md) | Repository aliases |
| 49 | [changelog-generate.md](01-app/49-changelog-generate.md) | Changelog generation |
| 50 | [ssh-keys.md](01-app/50-ssh-keys.md) | SSH key management |
| 51 | [prune.md](01-app/51-prune.md) | Prune command |
| 52 | [upload-retry.md](01-app/52-upload-retry.md) | Upload retry logic |
| 53 | [offline-detection.md](01-app/53-offline-detection.md) | Offline detection |
| 54 | [process-locking.md](01-app/54-process-locking.md) | Process locking |
| 55 | [temp-release.md](01-app/55-temp-release.md) | Temp release |
| 56 | [unified-gitmap-dir.md](01-app/56-unified-gitmap-dir.md) | Unified .gitmap directory |
| 57 | [skipmeta-integration-test.md](01-app/57-skipmeta-integration-test.md) | Skip-meta integration test |
| 58–78 | refactor-*.md | Code refactors (workflow, dispatch, archive, autocommit, SEO, branches, assets, TUI, aliases, ops, status, exec, logs, compress) |
| 79 | [task-watch.md](01-app/79-task-watch.md) | Task watch file sync |
| 80 | [env.md](01-app/80-env.md) | Environment variable management |
| 81 | [install.md](01-app/81-install.md) | Developer tool installer |
| 82 | [install-script.md](01-app/82-install-script.md) | One-liner install scripts |
| 83 | [future-features.md](01-app/82-future-features.md) | Future features (pending) |
| 84 | [chocolatey-package.md](01-app/84-chocolatey-package.md) | Chocolatey package distribution (research) |
| 85 | [winget-package.md](01-app/85-winget-package.md) | Winget package distribution (research) |
| 86 | [86-author-section.md](86-author-section.md) | Author section specification |

---

## 02-app-issues/ — Issue Post-Mortems

Root-cause analyses and resolution records for production bugs.

| # | File | Issue |
|---|------|-------|
| 01 | [update-file-lock.md](02-app-issues/01-update-file-lock.md) | Update file lock contention |
| 02 | [update-flow-spec-alignment.md](02-app-issues/02-update-flow-spec-alignment.md) | Update flow spec alignment |
| 03 | [update-sync-lock-loop.md](02-app-issues/03-update-sync-lock-loop.md) | Update sync lock loop |
| 04 | [database-path-resolution.md](02-app-issues/04-database-path-resolution.md) | Database path resolution |
| 05 | [list-empty-db-path.md](02-app-issues/05-list-empty-db-path.md) | List empty DB path |
| 06 | [release-orphaned-meta.md](02-app-issues/06-release-orphaned-meta.md) | Release orphaned metadata |
| 07 | [zip-group-release-silent-failure.md](02-app-issues/07-zip-group-release-silent-failure.md) | Zip group release silent failure |
| 08 | [autocommit-push-rejection.md](02-app-issues/08-autocommit-push-rejection.md) | Autocommit push rejection |
| 09 | [list-releases-repo-source.md](02-app-issues/09-list-releases-repo-source.md) | List releases repo source filter |
| 10 | [legacy-uuid-detection.md](02-app-issues/10-legacy-uuid-detection.md) | Legacy UUID detection |
| 11 | [auto-legacy-dir-migration.md](02-app-issues/11-auto-legacy-dir-migration.md) | Auto legacy directory migration |
| 12 | [legacy-id-migration.md](02-app-issues/12-legacy-id-migration.md) | Legacy ID migration |
| 13 | [release-pipeline-dist-directory.md](02-app-issues/13-release-pipeline-dist-directory.md) | Release pipeline dist directory error |
| 14 | [security-hardening-gosec-fixes.md](02-app-issues/14-security-hardening-gosec-fixes.md) | Security hardening (GoSec fixes) |
| 15 | [installer-progress-bar-and-binary-detection.md](02-app-issues/15-installer-progress-bar-and-binary-detection.md) | Installer crashes — progress bar & binary detection |
| 16 | [ci-passthrough-gate-pattern.md](02-app-issues/16-ci-passthrough-gate-pattern.md) | CI passthrough gate pattern |

---

## 03-general/ — Design Guidelines

Reusable architectural patterns and coding standards (generic, shareable across projects).

| # | File | Topic |
|---|------|-------|
| 01 | [cli-design-patterns.md](03-general/01-cli-design-patterns.md) | CLI design patterns |
| 02 | [powershell-build-deploy.md](03-general/02-powershell-build-deploy.md) | PowerShell build & deploy |
| 03 | [self-update-mechanism.md](03-general/03-self-update-mechanism.md) | Self-update mechanism |
| 04 | [output-formatting.md](03-general/04-output-formatting.md) | Output formatting standards |
| 05 | [config-pattern.md](03-general/05-config-pattern.md) | Configuration pattern |
| 06 | [code-style-rules.md](03-general/06-code-style-rules.md) | Code style rules |
| 07 | [date-display-format.md](03-general/07-date-display-format.md) | Date display format |

---

## 04-generic-cli/ — Generic CLI Blueprint

A production-quality CLI implementation blueprint usable as a starting point for any Go CLI project.

| # | File | Topic |
|---|------|-------|
| 01 | [overview.md](04-generic-cli/01-overview.md) | Blueprint overview |
| 02 | [project-structure.md](04-generic-cli/02-project-structure.md) | Project structure |
| 03 | [subcommand-architecture.md](04-generic-cli/03-subcommand-architecture.md) | Subcommand architecture |
| 04 | [flag-parsing.md](04-generic-cli/04-flag-parsing.md) | Flag parsing |
| 05 | [configuration.md](04-generic-cli/05-configuration.md) | Configuration |
| 06 | [output-formatting.md](04-generic-cli/06-output-formatting.md) | Output formatting |
| 07 | [error-handling.md](04-generic-cli/07-error-handling.md) | Error handling |
| 08 | [code-style.md](04-generic-cli/08-code-style.md) | Code style |
| 09 | [help-system.md](04-generic-cli/09-help-system.md) | Help system |
| 10 | [database.md](04-generic-cli/10-database.md) | Database integration |
| 11 | [build-deploy.md](04-generic-cli/11-build-deploy.md) | Build & deploy |
| 12 | [testing.md](04-generic-cli/12-testing.md) | Testing strategy |
| 13 | [checklist.md](04-generic-cli/13-checklist.md) | Implementation checklist |
| 14 | [date-formatting.md](04-generic-cli/14-date-formatting.md) | Date formatting |
| 15 | [constants-reference.md](04-generic-cli/15-constants-reference.md) | Constants reference |
| 16 | [verbose-logging.md](04-generic-cli/16-verbose-logging.md) | Verbose logging |
| 17 | [progress-tracking.md](04-generic-cli/17-progress-tracking.md) | Progress tracking |
| 18 | [batch-execution.md](04-generic-cli/18-batch-execution.md) | Batch execution |
| 19 | [shell-completion.md](04-generic-cli/19-shell-completion.md) | Shell completion |

---

## 05-coding-guidelines/ — Coding Standards

Coding conventions and style rules.

---

## 06-design-system/ — Design System

UI and design guidelines.

---

## 07-07-generic-release/ — Generic Release Pipeline

A reusable blueprint for cross-compiled CLI binary releases via CI/CD.

| # | File | Topic |
|---|------|-------|
| 01 | [cross-compilation.md](07-generic-release/01-cross-compilation.md) | Static binaries for 6+ targets |
| 02 | [release-pipeline.md](07-generic-release/02-release-pipeline.md) | CI/CD workflow structure |
| 03 | [install-scripts.md](07-generic-release/03-install-scripts.md) | Version-pinned installers |
| 04 | [checksums-verification.md](07-generic-release/04-checksums-verification.md) | SHA-256 verification |
| 05 | [release-assets.md](07-generic-release/05-release-assets.md) | Asset naming & packaging |
| 06 | [release-metadata.md](07-generic-release/06-release-metadata.md) | Version resolution & tagging |

Diagram: [`07-generic-release/images/release-pipeline-flow.jpg`](07-generic-release/images/release-pipeline-flow.jpg)

---

## 08-generic-update/ — Generic Self-Update

A reusable blueprint for CLI self-update: deploy-to-running-location, rename-first, handoff, cleanup.

| # | File | Topic |
|---|------|-------|
| 01 | [self-update-overview.md](08-generic-update/01-self-update-overview.md) | Problem, approach, platform differences |
| 02 | [deploy-path-resolution.md](08-generic-update/02-deploy-path-resolution.md) | Deploy to running location, PATH registration, data co-location |
| 03 | [rename-first-deploy.md](08-generic-update/03-rename-first-deploy.md) | Rename-first to bypass file locks |
| 04 | [build-scripts.md](08-generic-update/04-build-scripts.md) | `run.ps1` / `run.sh` build + deploy |
| 05 | [handoff-mechanism.md](08-generic-update/05-handoff-mechanism.md) | Copy-and-handoff for Windows |
| 06 | [cleanup.md](08-generic-update/06-cleanup.md) | Post-update artifact removal |

Diagram: [`08-generic-update/images/self-update-flow.jpg`](08-generic-update/images/self-update-flow.jpg)

---

## 09-pipeline/ — Pipeline Specifications

CI/CD pipeline architecture: CI, release, vulnerability scanning, installation, changelog, help, env, output, branding.

Diagram: [`09-pipeline/images/ci-pipeline-flow.jpg`](09-pipeline/images/ci-pipeline-flow.jpg)

---

## 86-author-section.md — Author Section Specification

Precise spec for author attribution blocks across all spec documents.
