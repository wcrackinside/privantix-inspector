# Deploy static site (privantix_site) to production FTP using WinSCP.
# Install WinSCP: https://winscp.net/
# Do NOT commit passwords. Use environment variable FTP_PASSWORD or prompt.
#
# Usage:
#   $env:FTP_PASSWORD = "secret"
#   .\scripts\deploy-ftp.ps1 -FtpHost "ftp.tudominio.com" -FtpUser "usuario" -RemotePath "/public_html"
#
# Dry run (list what would happen):
#   .\scripts\deploy-ftp.ps1 -FtpHost "..." -FtpUser "..." -RemotePath "/public_html" -DryRun
# FTPS (TLS): .\scripts\deploy-ftp.ps1 -FtpHost "ftps1.us.cloudlogin.co" ... -UseFtps

param(
    [Parameter(Mandatory = $true)]
    [string] $FtpHost,
    [Parameter(Mandatory = $true)]
    [string] $FtpUser,
    [int] $FtpPort = 21,
    [string] $RemotePath = "/",
    [string] $LocalPath = "",
    [switch] $DryRun,
    [switch] $UseFtps
)

$ErrorActionPreference = "Stop"

$repoRoot = Split-Path $PSScriptRoot -Parent
if (-not $LocalPath) {
    $LocalPath = Join-Path $repoRoot "privantix_site"
}

if (-not (Test-Path $LocalPath)) {
    Write-Error "No existe la carpeta local: $LocalPath"
}

$winscp = Get-Command "winscp.com" -ErrorAction SilentlyContinue
if (-not $winscp) {
    $defaultWinScp = "${env:ProgramFiles(x86)}\WinSCP\WinSCP.com"
    if (Test-Path $defaultWinScp) {
        $winscp = @{ Source = $defaultWinScp }
    }
}

if (-not $winscp) {
    Write-Host @"
No se encontro WinSCP (winscp.com).

Opciones:
  1) Instala WinSCP desde https://winscp.net/ y vuelve a ejecutar este script.
  2) Usa FileZilla / cliente FTP y sube manualmente la carpeta:
       $LocalPath
     al directorio remoto (ej. public_html o www).
  3) Con WSL/Linux: lftp -c "open -u USER,PASS ftp://HOST; mirror -R --delete ./privantix_site /remote/path"

"@
    exit 1
}

$exe = if ($winscp.Source) { $winscp.Source } else { $winscp.Path }

$pass = $env:FTP_PASSWORD
if (-not $pass) {
    $secure = Read-Host -AsSecureString "Contrasena FTP"
    $BSTR = [System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($secure)
    $pass = [System.Runtime.InteropServices.Marshal]::PtrToStringAuto($BSTR)
}

$localWin = $LocalPath -replace '/', '\'
$rp = $RemotePath.Trim()
if (-not $rp) { $rp = "/" }
if ($rp[0] -ne '/') { $rp = "/" + $rp }

# Sesion FTP o FTPS; luego cd al directorio remoto y sync al directorio actual (.)
$scheme = if ($UseFtps) { "ftps" } else { "ftp" }
$sessionUrl = "${scheme}://$([uri]::EscapeDataString($FtpUser)):$([uri]::EscapeDataString($pass))@${FtpHost}:${FtpPort}/"

$tempScript = [System.IO.Path]::GetTempFileName()
$syncCmd = if ($DryRun) { "ls" } else { "synchronize remote -delete `"$localWin`" ." }
$lines = @(
    "option batch abort",
    "option confirm off",
    "open $sessionUrl",
    "cd $rp",
    $syncCmd,
    "exit"
)
Set-Content -Path $tempScript -Value $lines -Encoding UTF8

try {
    Write-Host "Local:  $localWin"
    Write-Host "Remoto: ${scheme}://${FtpHost}:${FtpPort}${RemotePath}"
    & $exe /script=$tempScript
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
}
finally {
    Remove-Item $tempScript -Force -ErrorAction SilentlyContinue
}

Write-Host "Listo."
