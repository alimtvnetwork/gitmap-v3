<#
.SYNOPSIS
    One-liner installer for gitmap CLI.

.DESCRIPTION
    Downloads the latest gitmap release from GitHub, verifies checksums,
    extracts to a local directory, and adds it to PATH.

.PARAMETER Version
    Install a specific version (e.g. v2.48.0). Default: latest.

.PARAMETER InstallDir
    Target directory. Default: $env:LOCALAPPDATA\gitmap

.PARAMETER NoPath
    Skip adding to PATH.

.PARAMETER Arch
    Force architecture (amd64, arm64). Default: auto-detect.

.EXAMPLE
    irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v2/main/gitmap/scripts/install.ps1 | iex

.EXAMPLE
    & ./install.ps1 -Version v2.48.0

.NOTES
    Repository: https://github.com/alimtvnetwork/gitmap-v2
#>

param(
    [string]$Version = "",
    [string]$InstallDir = "",
    [string]$Arch = "",
    [switch]$NoPath,
    [switch]$Uninstall
)

$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

$Repo = "alimtvnetwork/gitmap-v2"
$BinaryName = "gitmap.exe"

# --- Logging helpers ---

function Write-Step([string]$msg) {
    Write-Host "  $msg" -ForegroundColor Cyan
}

function Write-OK([string]$msg) {
    Write-Host "  $msg" -ForegroundColor Green
}

function Write-Err([string]$msg) {
    Write-Host "  $msg" -ForegroundColor Red
}

# --- Resolve install directory ---

function Resolve-InstallDir([string]$dir) {
    if ($dir -ne "") { return $dir }
    return Join-Path $env:LOCALAPPDATA "gitmap"
}

# --- Detect architecture ---

function Resolve-Arch([string]$arch) {
    if ($arch -ne "") { return $arch }

    $cpu = $env:PROCESSOR_ARCHITECTURE
    switch ($cpu) {
        "AMD64"   { return "amd64" }
        "ARM64"   { return "arm64" }
        "x86"     { return "amd64" }
        default   { return "amd64" }
    }
}

# --- Resolve version (latest or pinned) ---

function Resolve-Version([string]$version) {
    if ($version -ne "") { return $version }

    Write-Step "Fetching latest release..."
    $url = "https://api.github.com/repos/$Repo/releases/latest"

    try {
        $release = Invoke-RestMethod -Uri $url -UseBasicParsing
        return $release.tag_name
    }
    catch {
        Write-Err "Failed to fetch latest release: $_"
        exit 1
    }
}

# --- Download asset ---

function Get-Asset([string]$version, [string]$arch) {
    $assetName = "gitmap-${version}-windows-${arch}.zip"
    $baseUrl = "https://github.com/$Repo/releases/download/$version"
    $assetUrl = "$baseUrl/$assetName"
    $checksumUrl = "$baseUrl/checksums.txt"

    $tmpDir = Join-Path $env:TEMP "gitmap-install-$(Get-Random)"
    New-Item -ItemType Directory -Path $tmpDir -Force | Out-Null

    $zipPath = Join-Path $tmpDir $assetName
    $checksumPath = Join-Path $tmpDir "checksums.txt"

    Write-Step "Downloading $assetName ($version)..."

    try {
        Invoke-WebRequest -Uri $assetUrl -OutFile $zipPath -UseBasicParsing
        Invoke-WebRequest -Uri $checksumUrl -OutFile $checksumPath -UseBasicParsing
    }
    catch {
        Write-Err "Download failed: $_"
        Remove-Item $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
        exit 1
    }

    # Verify checksum
    Write-Step "Verifying checksum..."
    $expectedLine = (Get-Content $checksumPath | Where-Object { $_ -match $assetName })
    if (-not $expectedLine) {
        Write-Err "Asset not found in checksums.txt"
        Remove-Item $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
        exit 1
    }

    $expectedHash = ($expectedLine -split '\s+')[0]
    $actualHash = (Get-FileHash -Path $zipPath -Algorithm SHA256).Hash.ToLower()

    if ($actualHash -ne $expectedHash) {
        Write-Err "Checksum mismatch!"
        Write-Err "  Expected: $expectedHash"
        Write-Err "  Got:      $actualHash"
        Remove-Item $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
        exit 1
    }

    Write-OK "Checksum verified."
    return @{ ZipPath = $zipPath; TmpDir = $tmpDir }
}

# --- Extract and install ---

function Install-Binary([string]$zipPath, [string]$installDir) {
    Write-Step "Installing to $installDir..."

    if (-not (Test-Path $installDir)) {
        New-Item -ItemType Directory -Path $installDir -Force | Out-Null
    }

    $targetPath = Join-Path $installDir $BinaryName
    $extractDir = Join-Path $installDir ".install-extract"

    # Rename-first strategy for running binary
    if (Test-Path $targetPath) {
        $oldPath = "$targetPath.old"
        if (Test-Path $oldPath) { Remove-Item $oldPath -Force }
        Rename-Item $targetPath $oldPath -Force
    }

    if (Test-Path $extractDir) {
        Remove-Item $extractDir -Recurse -Force -ErrorAction SilentlyContinue
    }

    New-Item -ItemType Directory -Path $extractDir -Force | Out-Null
    Expand-Archive -Path $zipPath -DestinationPath $extractDir -Force

    # Match exact names OR versioned patterns like gitmap-v2.54.6-windows-amd64.exe
    $candidateNames = @(
        $BinaryName,
        [System.IO.Path]::GetFileNameWithoutExtension($BinaryName),
        "gitmap-windows-amd64.exe",
        "gitmap-windows-arm64.exe"
    )

    $extracted = Get-ChildItem -Path $extractDir -File -Recurse |
        Where-Object {
            ($candidateNames -icontains $_.Name) -or
            ($_.Name -match "^gitmap-v[\d.]+-windows-(amd64|arm64)\.exe$")
        } |
        Select-Object -First 1

    if (-not $extracted) {
        $availableFiles = Get-ChildItem -Path $extractDir -File -Recurse | Select-Object -ExpandProperty FullName
        Remove-Item $extractDir -Recurse -Force -ErrorAction SilentlyContinue
        Write-Err "Installed archive did not contain $BinaryName"
        if ($availableFiles) {
            Write-Err "Archive files:"
            foreach ($file in $availableFiles) {
                Write-Err "  $file"
            }
        }
        exit 1
    }

    Move-Item $extracted.FullName $targetPath -Force

    Remove-Item $extractDir -Recurse -Force -ErrorAction SilentlyContinue

    if (-not (Test-Path $targetPath)) {
        Write-Err "Install failed: $BinaryName was not written to $installDir"
        exit 1
    }

    # Cleanup .old
    $oldPath = "$targetPath.old"
    if (Test-Path $oldPath) {
        Remove-Item $oldPath -Force -ErrorAction SilentlyContinue
    }

    Write-OK "Installed $BinaryName to $installDir"
}

# --- Add to PATH ---

function Test-PathEntry([string]$pathValue, [string]$dir) {
    if ([string]::IsNullOrWhiteSpace($pathValue)) {
        return $false
    }

    $parts = $pathValue -split ";"

    foreach ($part in $parts) {
        if ($part.Trim() -ieq $dir) {
            return $true
        }
    }

    return $false
}

function Rebuild-SessionPath([string]$dir) {
    # Rebuild session PATH from registry (Machine + User) to pick up any changes
    $machinePath = [Environment]::GetEnvironmentVariable("PATH", "Machine")
    $userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    $parts = @()
    if ($machinePath) { $parts += $machinePath.TrimEnd(";") }
    if ($userPath) { $parts += $userPath.TrimEnd(";") }
    $rebuilt = $parts -join ";"

    # Ensure install dir is present even if not yet persisted
    if (-not (Test-PathEntry $rebuilt $dir)) {
        $rebuilt = $rebuilt.TrimEnd(";") + ";" + $dir
    }

    return $rebuilt
}

function Broadcast-EnvironmentChange {
    Add-Type -TypeDefinition @"
using System;
using System.Runtime.InteropServices;

public static class GitMapEnvNative {
    [DllImport("user32.dll", SetLastError = true, CharSet = CharSet.Auto)]
    public static extern IntPtr SendMessageTimeout(
        IntPtr hWnd,
        uint Msg,
        UIntPtr wParam,
        string lParam,
        uint fuFlags,
        uint uTimeout,
        out UIntPtr lpdwResult
    );
}
"@ -ErrorAction SilentlyContinue | Out-Null

    $HWND_BROADCAST = [IntPtr]0xffff
    $WM_SETTINGCHANGE = 0x001A
    $SMTO_ABORTIFHUNG = 0x0002
    [UIntPtr]$result = [UIntPtr]::Zero

    [void][GitMapEnvNative]::SendMessageTimeout(
        $HWND_BROADCAST,
        $WM_SETTINGCHANGE,
        [UIntPtr]::Zero,
        "Environment",
        $SMTO_ABORTIFHUNG,
        5000,
        [ref]$result
    )
}

function Add-ToPath([string]$dir) {
    $modified = @()

    # --- 1. Windows User PATH (registry) — covers CMD + new PowerShell windows ---
    $currentUserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    $userHasDir = Test-PathEntry $currentUserPath $dir

    if (-not $userHasDir) {
        if ([string]::IsNullOrWhiteSpace($currentUserPath)) {
            $newPath = $dir
        }
        else {
            $newPath = $currentUserPath.TrimEnd(";") + ";" + $dir
        }

        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        Broadcast-EnvironmentChange
        Write-OK "Added to User PATH (registry)."
        $modified += "User PATH (registry)"
    }
    else {
        Write-Step "Already in User PATH (registry)."
    }

    # --- 2. PowerShell $PROFILE — ensures PATH in all PowerShell sessions ---
    $psProfilePath = $PROFILE
    if ($psProfilePath) {
        $exportLine = "`$env:PATH = `"$dir;`$env:PATH`""
        $marker = "# gitmap-path"
        $markerLine = "$exportLine $marker"

        $profileExists = Test-Path $psProfilePath
        $alreadyPresent = $false

        if ($profileExists) {
            $content = Get-Content $psProfilePath -Raw -ErrorAction SilentlyContinue
            if ($content -and ($content -match [regex]::Escape($marker))) {
                $alreadyPresent = $true
            }
        }

        if (-not $alreadyPresent) {
            # Ensure parent directory exists
            $profileDir = Split-Path $psProfilePath -Parent
            if ($profileDir -and -not (Test-Path $profileDir)) {
                New-Item -ItemType Directory -Path $profileDir -Force | Out-Null
            }
            Add-Content -Path $psProfilePath -Value "`n$markerLine" -Encoding UTF8
            Write-OK "Added to PowerShell profile: $psProfilePath"
            $modified += "PowerShell `$PROFILE"
        }
        else {
            Write-Step "Already in PowerShell profile."
        }
    }

    # --- 3. Git Bash profiles (~/.bashrc, ~/.bash_profile) ---
    $homeDir = $env:USERPROFILE
    if ($homeDir) {
        $bashExportLine = "export PATH=`"$($dir -replace '\\','/'):`$PATH`""
        $bashMarker = "# gitmap-path"
        $bashProfiles = @(
            (Join-Path $homeDir ".bashrc"),
            (Join-Path $homeDir ".bash_profile")
        )

        foreach ($bashProfile in $bashProfiles) {
            $bashAlreadyPresent = $false
            $bashProfileName = Split-Path $bashProfile -Leaf

            if (Test-Path $bashProfile) {
                $bashContent = Get-Content $bashProfile -Raw -ErrorAction SilentlyContinue
                if ($bashContent -and ($bashContent -match [regex]::Escape($bashMarker))) {
                    $bashAlreadyPresent = $true
                }
            }

            if (-not $bashAlreadyPresent) {
                Add-Content -Path $bashProfile -Value "`n$bashExportLine $bashMarker" -Encoding UTF8
                Write-OK "Added to Git Bash profile: ~/$bashProfileName"
                $modified += "~/$bashProfileName"
            }
            else {
                Write-Step "Already in ~/$bashProfileName."
            }
        }
    }

    if ($modified.Count -gt 0) {
        return @{ Target = ($modified -join ", "); Status = "added" }
    }

    return @{ Target = "All profiles"; Status = "already present" }
}

# --- Remove from PATH (uninstall) ---

function Remove-FromPath([string]$dir) {
    $removed = @()
    $marker = "# gitmap-path"

    # --- 1. Windows User PATH (registry) ---
    $currentUserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($currentUserPath -and (Test-PathEntry $currentUserPath $dir)) {
        $parts = ($currentUserPath -split ";") | Where-Object { $_.Trim() -ine $dir -and $_.Trim() -ne "" }
        $newPath = $parts -join ";"
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        Broadcast-EnvironmentChange
        Write-OK "Removed from User PATH (registry)."
        $removed += "User PATH (registry)"
    }

    # --- 2. PowerShell $PROFILE ---
    $psProfilePath = $PROFILE
    if ($psProfilePath -and (Test-Path $psProfilePath)) {
        $lines = Get-Content $psProfilePath
        $filtered = $lines | Where-Object { $_ -notmatch [regex]::Escape($marker) }
        if ($filtered.Count -lt $lines.Count) {
            $filtered | Set-Content $psProfilePath -Encoding UTF8
            Write-OK "Removed marker lines from PowerShell profile: $psProfilePath"
            $removed += "PowerShell `$PROFILE"
        }
    }

    # --- 3. Git Bash profiles ---
    $homeDir = $env:USERPROFILE
    if ($homeDir) {
        $bashProfiles = @(
            (Join-Path $homeDir ".bashrc"),
            (Join-Path $homeDir ".bash_profile")
        )

        foreach ($bashProfile in $bashProfiles) {
            if (Test-Path $bashProfile) {
                $lines = Get-Content $bashProfile
                $filtered = $lines | Where-Object { $_ -notmatch [regex]::Escape($marker) }
                if ($filtered.Count -lt $lines.Count) {
                    $filtered | Set-Content $bashProfile -Encoding UTF8
                    $name = Split-Path $bashProfile -Leaf
                    Write-OK "Removed marker lines from ~/$name"
                    $removed += "~/$name"
                }
            }
        }
    }

    if ($removed.Count -gt 0) {
        Write-Host ""
        Write-OK "PATH entries removed from: $($removed -join ', ')"
    }
    else {
        Write-Step "No gitmap PATH entries found in any profile."
    }

    return $removed
}

function Write-InstallSummary([string]$version, [string]$binPath, [string]$installDir, [hashtable]$pathResult, [bool]$isNoPath) {
    Write-Host ""
    Write-Host "  -----------------------------------------------" -ForegroundColor DarkGray
    Write-Host "  gitmap install summary" -ForegroundColor White
    Write-Host "  -----------------------------------------------" -ForegroundColor DarkGray
    Write-Host "    Version    : $version"
    Write-Host "    Binary     : $binPath"
    Write-Host "    Install Dir: $installDir"

    if ($isNoPath) {
        Write-Host "    PATH       : skipped (-NoPath)"
        return
    }

    Write-Host "    PATH target: $($pathResult.Target) ($($pathResult.Status))"
    Write-Host "    Session    : refreshed for current PowerShell session"

    Write-Host ""
    Write-Host "  Profiles modified:" -ForegroundColor White
    Write-Host "    - User PATH (registry)  : CMD, new PowerShell windows"
    Write-Host "    - PowerShell `$PROFILE    : all PowerShell sessions"
    Write-Host "    - ~/.bashrc             : Git Bash interactive shells"
    Write-Host "    - ~/.bash_profile       : Git Bash login shells"

    Write-Host ""
    Write-Host "  If gitmap is not found in a new terminal, run:" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "    PowerShell:  `$env:PATH = `"$installDir;`$env:PATH`"" -ForegroundColor Cyan
    Write-Host "    CMD:         set PATH=$installDir;%PATH%" -ForegroundColor Cyan
    Write-Host "    Git Bash:    export PATH=`"$($installDir -replace '\\','/'):`$PATH`"" -ForegroundColor Cyan
    Write-Host ""
}

# --- Main ---

function Main {
    Write-Host ""
    Write-Host "  gitmap installer" -ForegroundColor White
    Write-Host "  github.com/$Repo" -ForegroundColor DarkGray
    Write-Host ""

    try {
        $resolvedArch = Resolve-Arch $Arch
        $resolvedVersion = Resolve-Version $Version
        $resolvedDir = Resolve-InstallDir $InstallDir

        $result = Get-Asset $resolvedVersion $resolvedArch

        try {
            Install-Binary $result.ZipPath $resolvedDir
        }
        finally {
            Remove-Item $result.TmpDir -Recurse -Force -ErrorAction SilentlyContinue
        }

        $pathResult = @{ Target = "-NoPath"; Status = "skipped" }
        if (-not $NoPath) {
            $pathResult = Add-ToPath $resolvedDir

            # Also try Chocolatey refreshenv if available
            $refreshCmd = Get-Command refreshenv -ErrorAction SilentlyContinue
            if ($refreshCmd) {
                try { refreshenv | Out-Null } catch {}
            }

            # Force-rebuild $env:PATH in this scope so gitmap is usable immediately
            $script:NewPath = Rebuild-SessionPath $resolvedDir
        }
        else {
            $script:NewPath = $env:PATH
        }

        return @{ InstallDir = $resolvedDir; NewPath = $script:NewPath; Version = $resolvedVersion; PathResult = $pathResult }
    }
    catch {
        Write-Err "Installation failed: $_"
        Write-Host ""
        Write-Err "If this persists, download manually from:"
        Write-Err "  https://github.com/$Repo/releases/latest"
        Write-Host ""
        return $null
    }
}

# --- Uninstall mode ---
if ($Uninstall) {
    Write-Host ""
    Write-Host "  gitmap uninstaller" -ForegroundColor White
    Write-Host ""

    $resolvedDir = Resolve-InstallDir $InstallDir
    $binPath = Join-Path $resolvedDir $BinaryName

    # Remove PATH entries from all profiles
    $removedProfiles = Remove-FromPath $resolvedDir

    # Remove binary and install directory
    if (Test-Path $binPath) {
        Remove-Item $binPath -Force -ErrorAction SilentlyContinue
        Write-OK "Removed binary: $binPath"
    }

    $oldPath = "$binPath.old"
    if (Test-Path $oldPath) {
        Remove-Item $oldPath -Force -ErrorAction SilentlyContinue
    }

    # Remove install dir if empty
    if ((Test-Path $resolvedDir) -and @(Get-ChildItem $resolvedDir).Count -eq 0) {
        Remove-Item $resolvedDir -Force -ErrorAction SilentlyContinue
        Write-OK "Removed empty directory: $resolvedDir"
    }

    # Rebuild session PATH without the dir
    $machinePath = [Environment]::GetEnvironmentVariable("PATH", "Machine")
    $userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    $parts = @()
    if ($machinePath) { $parts += ($machinePath -split ";") }
    if ($userPath) { $parts += ($userPath -split ";") }
    $env:PATH = ($parts | Where-Object { $_.Trim() -ine $resolvedDir -and $_.Trim() -ne "" }) -join ";"

    Write-Host ""
    Write-OK "gitmap has been uninstalled."
    Write-Host ""
    return
}

$installResult = Main

if (-not $installResult) {
    # Main failed gracefully — error already printed
    return
}

# Set $env:PATH at the TOP-LEVEL script scope (not inside a function)
# This ensures the change persists in the caller's session when run via iex
$env:PATH = $installResult.NewPath

# Verify the binary works
$binPath = Join-Path $installResult.InstallDir $BinaryName
$installedVersion = $installResult.Version
if (Test-Path $binPath) {
    Write-Host ""
    try {
        $versionOutput = & $binPath version 2>&1
        $installedVersion = ($versionOutput | Out-String).Trim()
        Write-OK "gitmap $installedVersion"
    }
    catch {
        Write-Err "Binary found but failed to run: $_"
    }
}
else {
    Write-Err "Binary not found at $binPath"
}

Write-InstallSummary $installedVersion $binPath $installResult.InstallDir $installResult.PathResult $NoPath.IsPresent

Write-Host ""
Write-OK "Done! Run 'gitmap --help' to get started."
Write-Host ""
