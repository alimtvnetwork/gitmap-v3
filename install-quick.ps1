<#
.SYNOPSIS
    Short interactive installer for gitmap on Windows.

.DESCRIPTION
    Prompts the user for an install drive/folder (with a sensible default),
    then delegates to the canonical gitmap/scripts/install.ps1 with that path.

    Run via one-liner:
      irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v3/main/install-quick.ps1 | iex

    Or locally:
      ./install-quick.ps1
      ./install-quick.ps1 -InstallDir "E:\Tools\gitmap"
#>

param(
    [string]$InstallDir = "",
    [string]$Version    = ""
)

$ErrorActionPreference = "Stop"
$ProgressPreference    = "SilentlyContinue"

$Repo          = "alimtvnetwork/gitmap-v3"
$InstallerUrl  = "https://raw.githubusercontent.com/$Repo/main/gitmap/scripts/install.ps1"
$DefaultDir    = "D:\gitmap"

function Read-InstallDir([string]$default) {
    Write-Host ""
    Write-Host "  gitmap quick installer" -ForegroundColor Cyan
    Write-Host "  ---------------------" -ForegroundColor DarkGray
    Write-Host "  Choose install folder. Press Enter to accept the default." -ForegroundColor Gray
    Write-Host "  Default: $default" -ForegroundColor DarkGray

    $answer = Read-Host "  Install path"
    if ([string]::IsNullOrWhiteSpace($answer)) { return $default }
    return $answer.Trim('"').Trim()
}

if ([string]::IsNullOrWhiteSpace($InstallDir)) {
    $InstallDir = Read-InstallDir $DefaultDir
}

Write-Host ""
Write-Host "  Installing gitmap to: $InstallDir" -ForegroundColor Green
Write-Host ""

$script = (Invoke-WebRequest -Uri $InstallerUrl -UseBasicParsing).Content
$block  = [ScriptBlock]::Create($script)

if ($Version -ne "") {
    & $block -InstallDir $InstallDir -Version $Version
} else {
    & $block -InstallDir $InstallDir
}
