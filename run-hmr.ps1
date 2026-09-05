#Requires -Version 5.1
<#
.SYNOPSIS
  Run Knowledge Center API + WebUI on this machine with hot-reload.

.DESCRIPTION
  Frees TCP ports 4005 (API) and 4006 (WebUI), then starts Air and Vite watchers.
  Does not register a Windows service.
#>
Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$ApiPort = 4005
$WebUiPort = 4006
$Root = $PSScriptRoot
$ApiDir = Join-Path $Root 'knowledge-center-api'
$WebDir = Join-Path $Root 'knowledge-center-webui'
$StateDir = Join-Path $Root '.armin\deploy\local-windows-hmr'
$LogDir = Join-Path $StateDir 'logs'
$StatePath = Join-Path $StateDir 'state.json'
$ApiLog = Join-Path $LogDir 'api.log'
$WebLog = Join-Path $LogDir 'webui.log'

function Write-Step([string]$Message) {
    Write-Host ">> $Message" -ForegroundColor Cyan
}

function Write-Ok([string]$Message) {
    Write-Host "OK  $Message" -ForegroundColor Green
}

function Write-Fail([string]$Message) {
    Write-Host "ERR $Message" -ForegroundColor Red
}

function Get-ListenPids([int]$Port) {
    $ids = New-Object 'System.Collections.Generic.List[int]'
    try {
        $conns = Get-NetTCPConnection -LocalPort $Port -State Listen -ErrorAction SilentlyContinue
        foreach ($conn in @($conns)) {
            if ($null -ne $conn.OwningProcess -and [int]$conn.OwningProcess -gt 4) {
                if (-not $ids.Contains([int]$conn.OwningProcess)) {
                    [void]$ids.Add([int]$conn.OwningProcess)
                }
            }
        }
    } catch {
        # Fallback when Get-NetTCPConnection is unavailable.
    }
    if ($ids.Count -eq 0) {
        $pattern = '(:' + $Port + '\s+).+LISTENING\s+(\d+)'
        foreach ($line in (& netstat.exe -ano)) {
            if ($line -match $pattern) {
                $procId = [int]$Matches[2]
                if ($procId -gt 4 -and -not $ids.Contains($procId)) {
                    [void]$ids.Add($procId)
                }
            }
        }
    }
    return @($ids)
}

function Test-PortListening([int]$Port) {
    return @(Get-ListenPids -Port $Port).Count -gt 0
}

function Stop-PortOccupants([int]$Port) {
    Write-Step "Releasing port $Port"
    for ($attempt = 1; $attempt -le 8; $attempt++) {
        $ids = @(Get-ListenPids -Port $Port)
        if ($ids.Count -eq 0) {
            Write-Ok "Port $Port is free"
            return
        }
        foreach ($procId in $ids) {
            Write-Host "    taskkill /F /T PID $procId (port $Port)"
            & taskkill.exe /F /T /PID $procId 2>$null | Out-Null
        }
        Start-Sleep -Milliseconds 400
    }
    if (Test-PortListening -Port $Port) {
        throw "Port $Port is still in use after killing listeners."
    }
}

function Assert-Command([string]$Name) {
    $cmd = Get-Command $Name -ErrorAction SilentlyContinue
    if (-not $cmd) {
        throw "Missing command '$Name' on PATH. Install it, then re-run this script."
    }
    return $cmd
}

function Resolve-AirPath {
    $cmd = Get-Command air -ErrorAction SilentlyContinue
    if ($cmd) {
        return $cmd.Source
    }
    $go = Get-Command go -ErrorAction SilentlyContinue
    if ($go) {
        $gopath = (& go env GOPATH).Trim()
        $candidate = Join-Path $gopath 'bin\air.exe'
        if (Test-Path -LiteralPath $candidate) {
            return $candidate
        }
    }
    throw "Missing 'air'. Install with: go install github.com/air-verse/air@latest"
}

function Test-TcpOpen([string]$HostName, [int]$Port) {
    try {
        $client = New-Object System.Net.Sockets.TcpClient
        $iar = $client.BeginConnect($HostName, $Port, $null, $null)
        $ok = $iar.AsyncWaitHandle.WaitOne(800)
        if (-not $ok) {
            $client.Close()
            return $false
        }
        $client.EndConnect($iar)
        $client.Close()
        return $true
    } catch {
        return $false
    }
}

function Ensure-Postgres {
    if (Test-TcpOpen -HostName '127.0.0.1' -Port 5434) {
        Write-Ok "Postgres already accepting connections on 127.0.0.1:5434"
        return
    }
    $compose = Join-Path $ApiDir 'docker-compose.yml'
    $docker = Get-Command docker -ErrorAction SilentlyContinue
    if (-not $docker -or -not (Test-Path -LiteralPath $compose)) {
        throw "Postgres is not listening on 127.0.0.1:5434. Start it, then re-run this script."
    }
    Write-Step "Starting Postgres via docker compose (host 5434)"
    & docker network create knowledge-center-net 2>$null | Out-Null
    Push-Location $ApiDir
    try {
        & docker compose up -d postgres
        if ($LASTEXITCODE -ne 0) {
            throw "docker compose up -d postgres failed."
        }
    } finally {
        Pop-Location
    }
    for ($i = 1; $i -le 30; $i++) {
        if (Test-TcpOpen -HostName '127.0.0.1' -Port 5434) {
            Write-Ok "Postgres is ready on 127.0.0.1:5434"
            return
        }
        Start-Sleep -Seconds 1
    }
    throw "Postgres did not become ready on 127.0.0.1:5434."
}

function Start-LoggedWatcher {
    param(
        [string]$FilePath,
        [string]$Arguments,
        [string]$WorkingDirectory,
        [string]$LogPath,
        [string]$LauncherName
    )
    $launcher = Join-Path $StateDir $LauncherName
    @(
        '@echo off',
        'setlocal',
        "cd /d `"$WorkingDirectory`"",
        "`"$FilePath`" $Arguments > `"$LogPath`" 2>&1"
    ) | Set-Content -LiteralPath $launcher -Encoding ASCII
    $proc = Start-Process -FilePath 'cmd.exe' -ArgumentList "/c `"$launcher`"" -WindowStyle Hidden -PassThru
    if (-not $proc) {
        throw "Failed to start $FilePath"
    }
    return $proc
}

function Wait-Port([int]$Port, [string]$Name) {
    for ($i = 1; $i -le 45; $i++) {
        if (Test-PortListening -Port $Port) {
            Write-Ok "$Name is listening on port $Port"
            return
        }
        Start-Sleep -Seconds 1
    }
    throw "$Name did not start listening on port $Port. See logs under $LogDir"
}

if (-not (Test-Path -LiteralPath $ApiDir)) {
    throw "API folder not found: $ApiDir"
}
if (-not (Test-Path -LiteralPath $WebDir)) {
    throw "WebUI folder not found: $WebDir"
}

Assert-Command go | Out-Null
Assert-Command node | Out-Null
$npm = Assert-Command npm.cmd
$airPath = Resolve-AirPath

New-Item -ItemType Directory -Force -Path $LogDir | Out-Null

Stop-PortOccupants -Port $ApiPort
Stop-PortOccupants -Port $WebUiPort
Ensure-Postgres

$cors = @(
    "http://localhost:$WebUiPort",
    "http://127.0.0.1:$WebUiPort",
    "http://localhost:5173",
    "http://127.0.0.1:5173"
) -join ','

$env:HTTP_ADDR = ":$ApiPort"
$env:CORS_ORIGINS = $cors
$env:WEBUI_PORT = "$WebUiPort"
$env:API_PROXY_TARGET = "http://127.0.0.1:$ApiPort"
$env:VITE_API_BASE_URL = ''

if (-not (Test-Path -LiteralPath (Join-Path $WebDir 'node_modules'))) {
    Write-Step "Installing WebUI packages (npm install)"
    Push-Location $WebDir
    try {
        & $npm.Source install
        if ($LASTEXITCODE -ne 0) {
            throw "npm install failed."
        }
    } finally {
        Pop-Location
    }
}

Write-Step "Starting API hot-reload (air) on port $ApiPort"
$apiProc = Start-LoggedWatcher -FilePath $airPath -Arguments '-c .air.toml' -WorkingDirectory $ApiDir -LogPath $ApiLog -LauncherName 'start-api.cmd'

Write-Step "Starting WebUI hot-reload (Vite) on port $WebUiPort"
$webArgs = "run dev -- --host 0.0.0.0 --port $WebUiPort --strictPort"
$webProc = Start-LoggedWatcher -FilePath $npm.Source -Arguments $webArgs -WorkingDirectory $WebDir -LogPath $WebLog -LauncherName 'start-webui.cmd'

Wait-Port -Port $ApiPort -Name 'API'
Wait-Port -Port $WebUiPort -Name 'WebUI'

$state = [ordered]@{
    stack_name = 'knowledge-center'
    ports      = @{ api = $ApiPort; webui = $WebUiPort; postgres = 5434 }
    pids       = @{ api = $apiProc.Id; webui = $webProc.Id }
    commands   = @{ api = 'air'; webui = 'npm run dev' }
}
$state | ConvertTo-Json -Depth 5 | Set-Content -LiteralPath $StatePath -Encoding UTF8

Write-Ok "Hot-reload is running on this machine (host processes, not Docker for API/WebUI)"
Write-Host ""
Write-Host "API     http://127.0.0.1:$ApiPort"
Write-Host "WebUI   http://127.0.0.1:$WebUiPort"
Write-Host "Logs    $LogDir"
Write-Host "State   $StatePath"
Write-Host "Stop    re-run this script (it kills 4005/4006) or: taskkill /F /T /PID $($apiProc.Id) & taskkill /F /T /PID $($webProc.Id)"
