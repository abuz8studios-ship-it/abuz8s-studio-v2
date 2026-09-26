# ABUZ8s Studio v2 — one-command Windows build.
# Usage:  powershell -ExecutionPolicy Bypass -File scripts\build.ps1
# Output: backend\ABUZ8sStudio.exe (frontend embedded, no console in prod use)
param([string]$Out = "ABUZ8sStudio.exe")
$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot

Write-Host "==> 1/4 frontend build"
Push-Location (Join-Path $root "frontend")
if (-not (Test-Path "node_modules")) { npm install }
npm run build
Pop-Location

Write-Host "==> 2/4 sync embedded assets (backend\frontend\dist)"
$embed = Join-Path $root "backend\frontend\dist"
if (Test-Path $embed) { Remove-Item $embed -Recurse -Force }
Copy-Item (Join-Path $root "frontend\dist") $embed -Recurse -Force

Write-Host "==> 3/4 go vet + tests"
Push-Location (Join-Path $root "backend")
go vet ./...
go test ./...

Write-Host "==> 4/4 go build -> $Out"
# NOTE: Wails requires the desktop,production tags; without them the exe
# boots to an "Error" dialog instead of the app.
go build -tags desktop,production -trimpath -ldflags "-s -w" -o $Out .
Pop-Location

$exe = Join-Path $root "backend\$Out"
Write-Host "OK: $exe ($([math]::Round((Get-Item $exe).Length / 1MB, 1)) MB)"
Write-Host "Run:  & '$exe'"
