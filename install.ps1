$ErrorActionPreference = 'Stop'

Write-Host "Installing Summer CLI for Windows..." -ForegroundColor Cyan

# Detect Architecture
$Arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
    $Arch = "arm64"
}

$BinaryUrl = "https://github.com/DropGuard/summer-cli/releases/latest/download/summer-windows-$Arch.exe"
$InstallDir = "$env:LOCALAPPDATA\Programs\summer\bin"
$BinaryPath = "$InstallDir\summer.exe"

# Create install directory
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}

Write-Host "Downloading $BinaryUrl..." -ForegroundColor Gray
try {
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
    Invoke-WebRequest -Uri $BinaryUrl -OutFile $BinaryPath -UseBasicParsing
} catch {
    Write-Error "Failed to download Summer CLI. Please check your network or visit https://github.com/DropGuard/summer-cli/releases"
    exit 1
}

# Add to User PATH if not present
$UserPath = [Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::User)
if ($UserPath -notlike "*$InstallDir*") {
    Write-Host "Adding $InstallDir to User PATH..." -ForegroundColor Gray
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", [EnvironmentVariableTarget]::User)
    $env:Path += ";$InstallDir"
}

Write-Host "Summer CLI installed successfully!" -ForegroundColor Green
Write-Host "Run 'summer --help' to get started." -ForegroundColor Cyan
