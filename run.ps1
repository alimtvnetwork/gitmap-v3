<#
.SYNOPSIS
    Build, deploy, and run gitmap CLI from the repo root.
.DESCRIPTION
    Pulls latest code, resolves Go dependencies, builds the binary
    into ./bin, copies data folder, deploys to a target directory,
    and optionally runs gitmap with any arguments.
.EXAMPLES
    .\run.ps1                                    # pull, build, deploy
    .\run.ps1 -NoPull                            # skip git pull
    .\run.ps1 -ForcePull                         # discard local changes + pull (no prompt)
    .\run.ps1 -NoDeploy                          # skip deploy step
    .\run.ps1 -R scan                            # build + scan parent folder
    .\run.ps1 -R scan D:\repos                   # build + scan specific path
    .\run.ps1 -R scan D:\repos --mode ssh        # build + scan with flags
    .\run.ps1 -R clone .\gitmap-output\gitmap.json --target-dir .\restored
    .\run.ps1 -R help                            # build + show help
    .\run.ps1 -NoPull -NoDeploy -R scan          # just build and scan
    .\run.ps1 -t                                 # run all unit tests with reports
.NOTES
    Configuration is read from gitmap/powershell.json.
    -R accepts ALL gitmap CLI arguments after it (scan, clone, help, flags, paths).
    If -R is used with no arguments, it defaults to: scan <parent folder>
    -t runs all Go unit tests and writes reports to gitmap/data/unit-test-reports/.
    -ForcePull automatically discards local changes and removes untracked files
    before pulling. Useful for CI or unattended builds.
#>

[CmdletBinding(PositionalBinding=$false)]
param(
    [switch]$NoPull,
    [switch]$NoDeploy,
    [switch]$ForcePull,
    [string]$DeployPath = "",
    [Alias("d")]
    [switch]$Deploy,
    [switch]$Update,
    [switch]$R,
    [Alias("t")]
    [switch]$Test,
    [Parameter(ValueFromRemainingArguments=$true)]
    [string[]]$RunArgs
)

$ErrorActionPreference = "Stop"
$RepoRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$GitMapDir = Join-Path $RepoRoot "gitmap"

# -- Logging helpers -------------------------------------------
function Write-Step {
    param([string]$Step, [string]$Message)
    Write-Host ""
    Write-Host "  [$Step] " -ForegroundColor Magenta -NoNewline
    Write-Host $Message -ForegroundColor White
    Write-Host ("  " + ("-" * 50)) -ForegroundColor DarkGray
}

function Write-Success {
    param([string]$Message)
    Write-Host "  OK " -ForegroundColor Green -NoNewline
    Write-Host $Message -ForegroundColor Green
}

function Write-Info {
    param([string]$Message)
    Write-Host "  -> " -ForegroundColor Cyan -NoNewline
    Write-Host $Message -ForegroundColor Gray
}

function Write-Warn {
    param([string]$Message)
    Write-Host "  !! " -ForegroundColor Yellow -NoNewline
    Write-Host $Message -ForegroundColor Yellow
}

function Write-Fail {
    param([string]$Message)
    Write-Host "  XX " -ForegroundColor Red -NoNewline
    Write-Host $Message -ForegroundColor Red
}

# -- Banner ----------------------------------------------------
function Show-Banner {
    Write-Host ""
    Write-Host "  +======================================+" -ForegroundColor DarkCyan
    Write-Host "  |         " -ForegroundColor DarkCyan -NoNewline
    Write-Host "gitmap builder" -ForegroundColor Cyan -NoNewline
    Write-Host "              |" -ForegroundColor DarkCyan
    Write-Host "  +======================================+" -ForegroundColor DarkCyan
    Write-Host ""
}

# -- Load config -----------------------------------------------
function Load-Config {
    $configPath = Join-Path $GitMapDir "powershell.json"
    if (Test-Path $configPath) {
        Write-Info "Config loaded from powershell.json"

        return Get-Content $configPath | ConvertFrom-Json
    }
    Write-Warn "No powershell.json found, using defaults"

    return @{
        deployPath  = "E:\bin-run"
        buildOutput = "./bin"
        binaryName  = "gitmap.exe"
        copyData    = $true
    }
}

# -- Ensure main branch ----------------------------------------
function Ensure-MainBranch {
    Push-Location $RepoRoot
    try {
        $prevPref = $ErrorActionPreference
        $ErrorActionPreference = "Continue"
        $currentBranch = (git rev-parse --abbrev-ref HEAD 2>&1).Trim()
        $ErrorActionPreference = $prevPref

        if ($currentBranch -ne "main") {
            Write-Warn "Currently on branch '$currentBranch', switching to main..."
            $ErrorActionPreference = "Continue"
            $checkoutOutput = git checkout main 2>&1
            $checkoutExit = $LASTEXITCODE
            $ErrorActionPreference = $prevPref

            if ($checkoutExit -ne 0) {
                Write-Fail "Failed to switch to main branch"
                foreach ($line in $checkoutOutput) {
                    Write-Host "  $line" -ForegroundColor Red
                }
                exit 1
            }
            Write-Success "Switched to main branch"
        }
    } finally {
        Pop-Location
    }
}

# -- Git pull --------------------------------------------------
function Invoke-GitPull {
    Write-Step "1/4" "Pulling latest changes"

    Ensure-MainBranch

    Push-Location $RepoRoot
    try {
        # Temporarily allow stderr output from git without throwing NativeCommandError.
        $prevPref = $ErrorActionPreference
        $ErrorActionPreference = "Continue"
        $output = git pull 2>&1
        $pullExit = $LASTEXITCODE
        $ErrorActionPreference = $prevPref

        foreach ($line in $output) {
            $text = "$line".Trim()
            if ($text.Length -gt 0) {
                Write-Info $text
            }
        }

        if ($pullExit -ne 0) {
            $outputText = ($output | ForEach-Object { "$_" }) -join "`n"
            $hasConflict = $outputText -match "Your local changes" -or
                           $outputText -match "overwritten by merge" -or
                           $outputText -match "not possible because you have unmerged" -or
                           $outputText -match "Please commit your changes or stash them"

            if ($hasConflict) {
                if ($ForcePull) {
                    Write-Warn "Force-pull: discarding local changes and removing untracked files..."
                    $prevPref = $ErrorActionPreference
                    $ErrorActionPreference = "Continue"

                    $resetOutput = git checkout -- . 2>&1
                    $resetExit = $LASTEXITCODE
                    if ($resetExit -ne 0) {
                        Write-Fail "Git checkout failed"
                        $ErrorActionPreference = $prevPref
                        exit 1
                    }
                    Write-Success "Local changes discarded"

                    $cleanOutput = git clean -fd 2>&1
                    $cleanExit = $LASTEXITCODE
                    $ErrorActionPreference = $prevPref

                    if ($cleanExit -ne 0) {
                        Write-Fail "Git clean failed"
                        exit 1
                    }

                    $cleanedFiles = @($cleanOutput | ForEach-Object { "$_".Trim() } | Where-Object { $_.Length -gt 0 })
                    if ($cleanedFiles.Count -gt 0) {
                        Write-Success "Removed $($cleanedFiles.Count) untracked file(s)"
                    }

                    Retry-GitPull
                } else {
                    Resolve-PullConflict
                }
            } else {
                Write-Fail "Git pull failed (exit code $pullExit)"
                exit 1
            }
        } else {
            Write-Success "Pull complete"
        }
    } finally {
        Pop-Location
    }
}

# -- Resolve pull conflict with local changes ------------------
function Resolve-PullConflict {
    Write-Warn "Git pull failed due to local changes"
    Write-Host ""
    Write-Host "  Choose how to proceed:" -ForegroundColor Yellow
    Write-Host "    [S] Stash changes (save for later, then pull)" -ForegroundColor Cyan
    Write-Host "    [D] Discard changes (reset working tree, then pull)" -ForegroundColor Cyan
    Write-Host "    [C] Clean all (discard changes + remove untracked files, then pull)" -ForegroundColor Cyan
    Write-Host "    [Q] Quit (abort without changes)" -ForegroundColor Cyan
    Write-Host ""

    $choice = Read-Host "  Enter choice (S/D/C/Q)"

    switch ($choice.ToUpper()) {
        "S" {
            Write-Info "Stashing local changes..."
            $prevPref = $ErrorActionPreference
            $ErrorActionPreference = "Continue"
            $stashOutput = git stash push -m "auto-stash before run.ps1 pull" 2>&1
            $stashExit = $LASTEXITCODE
            $ErrorActionPreference = $prevPref

            if ($stashExit -ne 0) {
                Write-Fail "Git stash failed"
                foreach ($line in $stashOutput) {
                    Write-Host "  $line" -ForegroundColor Red
                }
                exit 1
            }
            Write-Success "Changes stashed"
            Write-Info "Run 'git stash pop' later to restore your changes"

            Retry-GitPull
        }
        "D" {
            Write-Warn "Discarding all local changes..."
            $prevPref = $ErrorActionPreference
            $ErrorActionPreference = "Continue"
            $resetOutput = git checkout -- . 2>&1
            $resetExit = $LASTEXITCODE
            $ErrorActionPreference = $prevPref

            if ($resetExit -ne 0) {
                Write-Fail "Git checkout failed"
                foreach ($line in $resetOutput) {
                    Write-Host "  $line" -ForegroundColor Red
                }
                exit 1
            }
            Write-Success "Local changes discarded"

            Retry-GitPull
        }
        "C" {
            Write-Warn "Discarding all local changes and removing untracked files..."
            $prevPref = $ErrorActionPreference
            $ErrorActionPreference = "Continue"

            $resetOutput = git checkout -- . 2>&1
            $resetExit = $LASTEXITCODE

            if ($resetExit -ne 0) {
                Write-Fail "Git checkout failed"
                foreach ($line in $resetOutput) {
                    Write-Host "  $line" -ForegroundColor Red
                }
                $ErrorActionPreference = $prevPref
                exit 1
            }
            Write-Success "Local changes discarded"

            $cleanOutput = git clean -fd 2>&1
            $cleanExit = $LASTEXITCODE
            $ErrorActionPreference = $prevPref

            if ($cleanExit -ne 0) {
                Write-Fail "Git clean failed"
                foreach ($line in $cleanOutput) {
                    Write-Host "  $line" -ForegroundColor Red
                }
                exit 1
            }

            $cleanedFiles = @($cleanOutput | ForEach-Object { "$_".Trim() } | Where-Object { $_.Length -gt 0 })
            if ($cleanedFiles.Count -gt 0) {
                foreach ($line in $cleanedFiles) {
                    Write-Info $line
                }
                Write-Success "Removed $($cleanedFiles.Count) untracked file(s)"
            } else {
                Write-Info "No untracked files to remove"
            }

            Retry-GitPull
        }
        default {
            Write-Info "Aborted by user"
            exit 0
        }
    }
}

# -- Retry git pull after stash/discard -----------------------
function Retry-GitPull {
    Write-Info "Retrying git pull..."
    $prevPref = $ErrorActionPreference
    $ErrorActionPreference = "Continue"
    $retryOutput = git pull 2>&1
    $retryExit = $LASTEXITCODE
    $ErrorActionPreference = $prevPref

    foreach ($line in $retryOutput) {
        $text = "$line".Trim()
        if ($text.Length -gt 0) {
            Write-Info $text
        }
    }

    if ($retryExit -ne 0) {
        Write-Fail "Git pull failed again (exit code $retryExit)"
        exit 1
    }

    Write-Success "Pull complete"
}

# -- Resolve dependencies -------------------------------------
function Resolve-Dependencies {
    Write-Step "2/4" "Resolving Go dependencies"
    Push-Location $GitMapDir
    try {
        $prevPref = $ErrorActionPreference
        $ErrorActionPreference = "Continue"
        $tidyOutput = go mod tidy 2>&1
        $tidyExit = $LASTEXITCODE
        $ErrorActionPreference = $prevPref

        if ($tidyExit -ne 0) {
            Write-Fail "go mod tidy failed"
            foreach ($line in $tidyOutput) {
                Write-Host "  $line" -ForegroundColor Red
            }
            exit 1
        }
        Write-Success "Dependencies resolved"
    } finally {
        Pop-Location
    }
}

# -- Pre-build validation --------------------------------------
function Test-SourceFiles {
    Write-Info "Validating source files..."

    $requiredFiles = @(
        "main.go",
        "go.mod",
        "cmd/root.go",
        "cmd/scan.go",
        "cmd/clone.go",
        "cmd/update.go",
        "cmd/pull.go",
        "cmd/rescan.go",
        "cmd/desktopsync.go",
        "constants/constants.go",
        "config/config.go",
        "scanner/scanner.go",
        "mapper/mapper.go",
        "model/record.go",
        "formatter/csv.go",
        "formatter/json.go",
        "formatter/terminal.go",
        "formatter/text.go",
        "formatter/structure.go",
        "formatter/clonescript.go",
        "formatter/directclone.go",
        "formatter/desktopscript.go",
        "cloner/cloner.go",
        "cloner/safe_pull.go",
        "gitutil/gitutil.go",
        "desktop/desktop.go",
        "verbose/verbose.go",
        "setup/setup.go",
        "cmd/setup.go",
        "cmd/status.go",
        "cmd/exec.go",
        "cmd/release.go",
        "cmd/releasebranch.go",
        "cmd/releasepending.go",
        "cmd/changelog.go",
        "cmd/doctor.go",
        "release/semver.go",
        "release/metadata.go",
        "release/gitops.go",
        "release/github.go",
        "release/changelog.go",
        "release/workflow.go"
    )

    $missing = @()
    foreach ($file in $requiredFiles) {
        $fullPath = Join-Path $GitMapDir $file
        if (-not (Test-Path $fullPath)) {
            $missing += $file
        }
    }

    if ($missing.Count -gt 0) {
        Write-Fail "Missing source files ($($missing.Count)):"
        foreach ($f in $missing) {
            Write-Host "  - $f" -ForegroundColor Red
        }
        exit 1
    }

    Write-Success "All $($requiredFiles.Count) source files present"
}

# -- Build binary ----------------------------------------------
function Build-Binary {
    param($Config)

    Write-Step "3/4" "Building $($Config.binaryName)"
    Test-SourceFiles

    $binDir  = Join-Path $RepoRoot $Config.buildOutput
    $outPath = Join-Path $binDir $Config.binaryName

    if (-not (Test-Path $binDir)) {
        New-Item -ItemType Directory -Path $binDir -Force | Out-Null
        Write-Info "Created bin directory"
    }

    Push-Location $GitMapDir
    try {
        $absRepoRoot = (Resolve-Path $RepoRoot).Path
        $ldflags = "-X 'github.com/user/gitmap/constants.RepoPath=$absRepoRoot'"

        $prevPref = $ErrorActionPreference
        $ErrorActionPreference = "Continue"
        $buildOutput = go build -ldflags $ldflags -o $outPath . 2>&1
        $buildExit = $LASTEXITCODE
        $ErrorActionPreference = $prevPref

        if ($buildExit -ne 0) {
            Write-Fail "Go build failed"
            foreach ($line in $buildOutput) {
                $text = "$line".Trim()
                if ($text.Length -gt 0) {
                    Write-Host "  $text" -ForegroundColor Red
                }
            }
            exit 1
        }
    } finally {
        Pop-Location
    }

    if ($Config.copyData) {
        Copy-DataFolder -BinDir $binDir
    }

    $size = (Get-Item $outPath).Length / 1MB
    Write-Success ("Binary built ({0:N2} MB) -> $outPath" -f $size)

    return $outPath
}

# -- Copy data folder -----------------------------------------
function Copy-DataFolder {
    param($BinDir)

    $dataSource = Join-Path $GitMapDir "data"
    $dataDest   = Join-Path $BinDir "data"

    if (Test-Path $dataSource) {
        if (Test-Path $dataDest) {
            Remove-Item $dataDest -Recurse -Force
        }
        Copy-Item $dataSource $dataDest -Recurse
        Write-Info "Copied data folder to bin"
    }
}

# -- Resolve deploy target -------------------------------------
# Priority: 1) -DeployPath flag  2) globally installed gitmap location  3) powershell.json default
function Resolve-DeployTarget {
    param($Config, $OverridePath)

    # 1) Explicit CLI override always wins
    if ($OverridePath.Length -gt 0) {
        Write-Info "Deploy target: CLI override -> $OverridePath"

        return $OverridePath
    }

    # 2) If gitmap is already on PATH, deploy to its parent directory
    $activeCmd = Get-Command gitmap -ErrorAction SilentlyContinue
    if ($activeCmd) {
        $activePath = $activeCmd.Source
        if (Test-Path $activePath) {
            $resolvedActive = (Resolve-Path $activePath).Path
            $activeDir = Split-Path $resolvedActive -Parent
            $activeDirName = Split-Path $activeDir -Leaf

            # The binary lives in <deploy-target>/gitmap/gitmap.exe
            # So the deploy target is the parent of the gitmap/ folder
            if ($activeDirName -eq "gitmap") {
                $deployTarget = Split-Path $activeDir -Parent
                Write-Info "Deploy target: detected from PATH -> $deployTarget"

                return $deployTarget
            }

            # Binary is directly in a folder (not nested under gitmap/)
            # Deploy target = that folder's parent so we create gitmap/ there
            $deployTarget = Split-Path $activeDir -Parent
            Write-Info "Deploy target: detected from PATH -> $deployTarget"

            return $deployTarget
        }
    }

    # 3) Fall back to powershell.json default
    Write-Info "Deploy target: powershell.json default -> $($Config.deployPath)"

    return $Config.deployPath
}

# -- Deploy to target directory --------------------------------
function Deploy-Binary {
    param($Config, $BinaryPath, $OverridePath)

    Write-Step "4/4" "Deploying"

    $target = Resolve-DeployTarget -Config $Config -OverridePath $OverridePath

    Write-Info "Target: $target"

    if (-not (Test-Path $target)) {
        New-Item -ItemType Directory -Path $target -Force | Out-Null
        Write-Info "Created deploy directory"
    }

    # Deploy into nested gitmap/ subfolder
    $appDir = Join-Path $target "gitmap"
    if (-not (Test-Path $appDir)) {
        New-Item -ItemType Directory -Path $appDir -Force | Out-Null
        Write-Info "Created gitmap app directory"
    }

    $destFile = Join-Path $appDir $Config.binaryName
    $backupFile = "$destFile.old"
    $hasBackup = $false
    $deploySuccess = $false

    if (Test-Path $destFile) {
        # Rename-first strategy: Windows allows renaming a running binary
        # but not overwriting it. Rename to .old, then copy the new one.
        try {
            if (Test-Path $backupFile) {
                Remove-Item $backupFile -Force -ErrorAction SilentlyContinue
            }
            Rename-Item $destFile $backupFile -Force -ErrorAction Stop
            $hasBackup = $true
            Write-Info "Renamed existing binary to $($Config.binaryName).old (rename-first)"
        } catch {
            Write-Warn "Rename-first failed: $_"
            # Fallback: try a backup copy instead
            try {
                Copy-Item $destFile $backupFile -Force -ErrorAction Stop
                $hasBackup = $true
                Write-Info "Backed up existing binary to $($Config.binaryName).old"
            } catch {
                Write-Warn "Could not create backup: $_"
            }
        }
    }

    # Copy new binary — after rename-first, the destination is free
    $maxAttempts = 5
    $attempt = 1
    while ($true) {
        try {
            Copy-Item $BinaryPath $destFile -Force -ErrorAction Stop
            $deploySuccess = $true
            break
        } catch {
            if ($attempt -ge $maxAttempts) {
                # Restore backup on failure
                if ($hasBackup -and (Test-Path $backupFile) -and (-not (Test-Path $destFile))) {
                    Write-Warn "Deploy failed - restoring previous binary from backup"
                    try {
                        Rename-Item $backupFile $destFile -Force -ErrorAction Stop
                        Write-Success "Rollback complete - previous version restored"
                    } catch {
                        Write-Fail "Rollback also failed: $_"
                    }
                }
                throw
            }
            Write-Warn "Target still locked; retrying ($attempt/$maxAttempts)..."
            Start-Sleep -Milliseconds 500
            $attempt++
        }
    }

    # Leave .old file in place - cleaned up by: gitmap update-cleanup
    if ($hasBackup -and $deploySuccess) {
        Write-Info "Previous binary kept as $($Config.binaryName).old (run 'gitmap update-cleanup' to remove)"
    }

    $binDir   = Split-Path $BinaryPath -Parent
    $dataDir  = Join-Path $binDir "data"
    $dataDest = Join-Path $appDir "data"
    if (Test-Path $dataDir) {
        if (Test-Path $dataDest) {
            Remove-Item $dataDest -Recurse -Force
        }
        Copy-Item $dataDir $dataDest -Recurse
        Write-Info "Copied data folder to gitmap app directory"
    }

    Write-Success "Deployed to $appDir"
    Write-Info "Ensure $appDir is on your PATH to run: gitmap"
}

# -- Run gitmap ------------------------------------------------
function Invoke-Run {
    param($Config, $BinaryPath, [string[]]$CliArgs)

    Write-Host ""
    Write-Step "RUN" "Executing gitmap"

    # Always run from the local bin build, never from the deploy target
    $binDir = Split-Path $BinaryPath -Parent
    $dataDir = Join-Path $binDir "data"

    $resolvedArgs = Resolve-RunArgs -CliArgs $CliArgs
    $argString = $resolvedArgs -join ' '
    $currentDir = (Get-Location).Path
    Write-Info "Binary: $BinaryPath"
    Write-Info "Runner CWD: $currentDir"
    Write-Info "Command: gitmap $argString"
    if ($resolvedArgs.Count -ge 2 -and $resolvedArgs[0] -eq "scan") {
        Write-Info "Scan target: $($resolvedArgs[1])"
    }
    Write-Host ("  " + ("-" * 50)) -ForegroundColor DarkGray
    Write-Host ""

    $proc = Start-Process -FilePath $BinaryPath -ArgumentList $resolvedArgs -WorkingDirectory $binDir -NoNewWindow -Wait -PassThru

    Write-Host ""
    if ($proc.ExitCode -eq 0) {
        Write-Success "Run complete"
    } else {
        Write-Fail "gitmap exited with code $($proc.ExitCode)"
    }
}

# -- Resolve run arguments -------------------------------------
function Resolve-RunArgs {
    param([string[]]$CliArgs)

    if ($CliArgs.Count -eq 0) {
        $parentDir = Split-Path $RepoRoot -Parent
        Write-Info "No args provided, defaulting to: scan $parentDir"

        return @("scan", $parentDir)
    }

    # Resolve relative paths to absolute so Start-Process always receives correct targets
    $baseDir = (Get-Location).Path
    $resolved = @()
    foreach ($arg in $CliArgs) {
        if ($arg -match '^(\.\.[\\/]|\.[\\/]|\.\.?$)' -and -not $arg.StartsWith('-')) {
            $path = Resolve-Path -LiteralPath $arg -ErrorAction SilentlyContinue
            if ($path) {
                $resolved += $path.Path
            } else {
                $resolved += [System.IO.Path]::GetFullPath((Join-Path $baseDir $arg))
            }
        } else {
            $resolved += $arg
        }
    }

    return $resolved
}

# -- Run tests -------------------------------------------------
function Invoke-Tests {
    Write-Step "TEST" "Running unit tests"

    $reportDir = Join-Path (Join-Path $GitMapDir "data") "unit-test-reports"
    if (-not (Test-Path $reportDir)) {
        New-Item -ItemType Directory -Path $reportDir -Force | Out-Null
        Write-Info "Created report directory: $reportDir"
    }

    $overallLog = Join-Path $reportDir "overall.log.txt"
    $failingLog = Join-Path $reportDir "failingTest.log.txt"

    Push-Location $GitMapDir
    try {
        $prevPref = $ErrorActionPreference
        $ErrorActionPreference = "Continue"

        Write-Info "Running: go test ./..."
        $testOutput = go test ./... -v -count=1 2>&1
        $testExit = $LASTEXITCODE
        $ErrorActionPreference = $prevPref

        # Write overall report
        $testOutput | Out-File -FilePath $overallLog -Encoding UTF8
        Write-Info "Overall report: $overallLog"

        # Extract failing tests
        $failLines = @()
        $currentTest = ""
        $inFail = $false
        foreach ($line in $testOutput) {
            $text = "$line"
            if ($text -match "^--- FAIL:") {
                $inFail = $true
                $currentTest = $text
                $failLines += ""
                $failLines += $text
            } elseif ($text -match "^--- PASS:" -or $text -match "^=== RUN") {
                $inFail = $false
            } elseif ($text -match "^FAIL\s") {
                $failLines += $text
            } elseif ($inFail) {
                $failLines += $text
            }
        }

        if ($failLines.Count -gt 0) {
            $failLines | Out-File -FilePath $failingLog -Encoding UTF8
            Write-Fail "Some tests failed. See: $failingLog"
        } else {
            "No failing tests." | Out-File -FilePath $failingLog -Encoding UTF8
            Write-Success "All tests passed"
        }

        # Print summary
        $passCount = ($testOutput | Where-Object { "$_" -match "^--- PASS:" }).Count
        $failCount = ($testOutput | Where-Object { "$_" -match "^--- FAIL:" }).Count
        $skipCount = ($testOutput | Where-Object { "$_" -match "^--- SKIP:" }).Count
        Write-Info "Results: $passCount passed, $failCount failed, $skipCount skipped"

        # Show test output in terminal
        foreach ($line in $testOutput) {
            $text = "$line".Trim()
            if ($text -match "^--- FAIL:") {
                Write-Host "  $text" -ForegroundColor Red
            } elseif ($text -match "^--- PASS:") {
                Write-Host "  $text" -ForegroundColor Green
            } elseif ($text -match "^FAIL") {
                Write-Host "  $text" -ForegroundColor Red
            } elseif ($text -match "^ok\s") {
                Write-Host "  $text" -ForegroundColor Green
            } elseif ($text.Length -gt 0) {
                Write-Host "  $text" -ForegroundColor Gray
            }
        }

        if ($testExit -ne 0) {
            Write-Fail "Tests failed (exit code $testExit)"
        }
    } finally {
        Pop-Location
    }
}

# -- Main ------------------------------------------------------
Show-Banner
$config = Load-Config

if ($Test) {
    Write-Info "Test mode enabled (-t)"
    Resolve-Dependencies
    Invoke-Tests
    Write-Host ""
    Write-Success "All done!"
    Write-Host ""
    exit 0
}

if ($Update) {
    Write-Info "Update mode enabled (-Update)"
}

if (-not $NoPull) {
    Invoke-GitPull
} else {
    Write-Info "Skipping git pull (-NoPull)"
}

Resolve-Dependencies
$binaryPath = Build-Binary -Config $config

# Show built version
$versionOutput = & $binaryPath version 2>&1
Write-Info "Version: $versionOutput"

$deployedBinaryPath = $null
if ($Deploy) { $NoDeploy = $false }
if (-not $NoDeploy) {
    Deploy-Binary -Config $config -BinaryPath $binaryPath -OverridePath $DeployPath

    $effectiveDeployPath = Resolve-DeployTarget -Config $config -OverridePath $DeployPath
    $deployedBinaryPath = Join-Path (Join-Path $effectiveDeployPath "gitmap") $config.binaryName

    $activeCmd = Get-Command gitmap -ErrorAction SilentlyContinue
    if ($activeCmd -and (Test-Path $deployedBinaryPath)) {
        $activeBinaryPath = $activeCmd.Source
        if (Test-Path $activeBinaryPath) {
            $activeResolved = (Resolve-Path $activeBinaryPath).Path
            $deployedResolved = (Resolve-Path $deployedBinaryPath).Path
            if ($activeResolved -ne $deployedResolved) {
                Write-Warn "PATH points to a different gitmap binary."
                Write-Info "Active:   $activeResolved"
                Write-Info "Deployed: $deployedResolved"

                $maxSyncAttempts = 20
                $syncSuccess = $false

                if ($Update) {
                    Write-Info "Update mode: using rename-first PATH sync"
                    $activeBackup = "$activeBinaryPath.old"
                    try {
                        if (Test-Path $activeBackup) {
                            Remove-Item $activeBackup -Force -ErrorAction SilentlyContinue
                        }
                        Rename-Item $activeBinaryPath $activeBackup -Force -ErrorAction Stop
                        Copy-Item $deployedBinaryPath $activeBinaryPath -Force -ErrorAction Stop
                        $syncedVersion = & $activeBinaryPath version 2>&1
                        Write-Success "Synced active PATH binary via rename-first -> $syncedVersion"
                        $syncSuccess = $true
                    } catch {
                        if ((Test-Path $activeBackup) -and (-not (Test-Path $activeBinaryPath))) {
                            try {
                                Copy-Item $activeBackup $activeBinaryPath -Force -ErrorAction Stop
                            } catch {
                            }
                        }
                        Write-Warn "Rename-first sync failed; retrying with copy loop"
                    }
                }

                if (-not $syncSuccess) {
                    for ($syncAttempt = 1; $syncAttempt -le $maxSyncAttempts; $syncAttempt++) {
                        try {
                            Copy-Item $deployedBinaryPath $activeBinaryPath -Force -ErrorAction Stop
                            $syncedVersion = & $activeBinaryPath version 2>&1
                            Write-Success "Synced active PATH binary -> $syncedVersion"
                            $syncSuccess = $true
                            break
                        } catch {
                            if ($syncAttempt -lt $maxSyncAttempts) {
                                Write-Warn "Active PATH binary is in use; retrying ($syncAttempt/$maxSyncAttempts)..."
                                Start-Sleep -Milliseconds 500
                            }
                        }
                    }
                }

                if (-not $syncSuccess) {
                    $activeBackup = "$activeBinaryPath.old"
                    try {
                        if (Test-Path $activeBackup) {
                            Remove-Item $activeBackup -Force -ErrorAction SilentlyContinue
                        }
                        Rename-Item $activeBinaryPath $activeBackup -Force -ErrorAction Stop
                        Copy-Item $deployedBinaryPath $activeBinaryPath -Force -ErrorAction Stop
                        $syncedVersion = & $activeBinaryPath version 2>&1
                        Write-Success "Synced active PATH binary via rename fallback -> $syncedVersion"
                        $syncSuccess = $true
                    } catch {
                        if ((Test-Path $activeBackup) -and (-not (Test-Path $activeBinaryPath))) {
                            try {
                                Copy-Item $activeBackup $activeBinaryPath -Force -ErrorAction Stop
                            } catch {
                            }
                        }
                    }
                }

                if (-not $syncSuccess) {
                    try {
                        $staleProcs = Get-CimInstance Win32_Process -Filter "Name='gitmap.exe'" -ErrorAction SilentlyContinue |
                            Where-Object { $_.ExecutablePath -and ((Resolve-Path $_.ExecutablePath -ErrorAction SilentlyContinue).Path -eq $activeResolved) -and ($_.ProcessId -ne $PID) }
                        foreach ($p in $staleProcs) {
                            Stop-Process -Id $p.ProcessId -Force -ErrorAction SilentlyContinue
                        }
                        if ($staleProcs) {
                            Start-Sleep -Milliseconds 500
                            Copy-Item $deployedBinaryPath $activeBinaryPath -Force -ErrorAction Stop
                            $syncedVersion = & $activeBinaryPath version 2>&1
                            Write-Success "Synced active PATH binary after stopping stale gitmap process(es) -> $syncedVersion"
                            $syncSuccess = $true
                        }
                    } catch {
                    }
                }

                if (-not $syncSuccess) {
                    Write-Warn "Could not sync active PATH binary after retries and fallback attempts."
                    Write-Info "Close terminals/apps using gitmap and run:"
                    Write-Info ('Copy-Item "' + $deployedBinaryPath + '" "' + $activeBinaryPath + '" -Force')
                    Write-Info ('Or run directly: "' + $deployedBinaryPath + '" <command>')
                }
            }
        }
    }
} else {
    Write-Info "Skipping deploy (-NoDeploy)"
}

$changelogBinaryPath = $binaryPath
$activeCmdForChangelog = Get-Command gitmap -ErrorAction SilentlyContinue
if ($activeCmdForChangelog -and (Test-Path $activeCmdForChangelog.Source)) {
    $changelogBinaryPath = $activeCmdForChangelog.Source
} elseif ($deployedBinaryPath -and (Test-Path $deployedBinaryPath)) {
    $changelogBinaryPath = $deployedBinaryPath
}

if (Test-Path $changelogBinaryPath) {
    Write-Host ""
    Write-Info "Latest changelog:"
    & $changelogBinaryPath changelog --latest

    if ($Update) {
        Write-Host ""
        Write-Info "Running update cleanup"
        & $changelogBinaryPath update-cleanup
    }
}

if ($R) {
    Invoke-Run -Config $config -BinaryPath $binaryPath -CliArgs $RunArgs
}

Write-Host ""
Write-Success "All done!"
Write-Host ""

# -- Last release info -----------------------------------------
$lastReleaseScript = Join-Path (Join-Path (Join-Path $RepoRoot "gitmap") "scripts") "Get-LastRelease.ps1"
if (Test-Path $lastReleaseScript) {
    $lrBinary = $changelogBinaryPath
    & $lastReleaseScript -BinaryPath $lrBinary -RepoRoot $RepoRoot
    Write-Host ""
}
