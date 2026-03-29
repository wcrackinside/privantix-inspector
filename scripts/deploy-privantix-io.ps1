# Despliegue de www.privantix.io — alineado con FileZilla (FTP estandar, puerto 21).
# FileZilla: Host privantix.io, Protocol 0 = FTP, User suptime_privantix.io
# Ruta remota del sitio: /www/privantix.io/privantix_site
#
# No incluye contrasena. Antes: $env:FTP_PASSWORD = "..."
#
# Uso:
#   $env:FTP_PASSWORD = "tu_clave"
#   .\scripts\deploy-privantix-io.ps1
# Prueba:
#   .\scripts\deploy-privantix-io.ps1 -DryRun
#
# Alternativa FTPS en cloudlogin (si tu panel indica ftps1.us.cloudlogin.co):
#   .\scripts\deploy-privantix-io.ps1 -CloudloginFtps

param(
    [string] $FtpUser = "",
    [switch] $DryRun,
    [switch] $CloudloginFtps
)

$u = if ($FtpUser) { $FtpUser } elseif ($env:FTP_USER) { $env:FTP_USER } else { "suptime_privantix.io" }

$child = Join-Path $PSScriptRoot "deploy-ftp.ps1"

if ($CloudloginFtps) {
    & $child `
        -FtpHost "ftps1.us.cloudlogin.co" `
        -FtpPort 21 `
        -FtpUser $u `
        -RemotePath "/www/privantix.io/privantix_site" `
        -DryRun:$DryRun `
        -UseFtps
} else {
    & $child `
        -FtpHost "privantix.io" `
        -FtpPort 21 `
        -FtpUser $u `
        -RemotePath "/www/privantix.io/privantix_site" `
        -DryRun:$DryRun
}
