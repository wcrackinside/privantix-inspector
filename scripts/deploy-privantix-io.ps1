# Despliegue de www.privantix.io (cloudlogin FTPS).
# No incluye contrasena. Antes: $env:FTP_PASSWORD = "..."
#
# Uso:
#   $env:FTP_PASSWORD = "tu_clave"
#   .\scripts\deploy-privantix-io.ps1
# Prueba sin subir:
#   .\scripts\deploy-privantix-io.ps1 -DryRun

param(
    [string] $FtpUser = "",
    [switch] $DryRun,
    [switch] $PlainFtp
)

$u = if ($FtpUser) { $FtpUser } elseif ($env:FTP_USER) { $env:FTP_USER } else { "suptime_privantix.io" }

$child = Join-Path $PSScriptRoot "deploy-ftp.ps1"
& $child `
    -FtpHost "ftps1.us.cloudlogin.co" `
    -FtpUser $u `
    -RemotePath "/www/privantix.io/privantix_site" `
    -DryRun:$DryRun `
    -UseFtps:$(-not $PlainFtp)
