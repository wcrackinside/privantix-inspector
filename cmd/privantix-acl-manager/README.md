# Privantix ACL Manager (Fyne)

Interfaz gráfica para visualizar auditorías de permisos y gestionar backups/restauración de ACLs en Windows.

## Requisitos de compilación

Fyne requiere **CGO** y un compilador C (GCC). En Windows:

### Opción 1: MinGW-w64 (recomendado)

1. Instale [MinGW-w64](https://www.mingw-w64.org/) o [MSYS2](https://www.msys2.org/)
2. Agregue el directorio `bin` de MinGW al PATH (ej: `C:\msys64\mingw64\bin`)
3. Verifique: `gcc -v`
4. Compile:

```powershell
$env:CGO_ENABLED = "1"
go build -o privantix-acl-manager.exe ./cmd/privantix-acl-manager
```

### Opción 2: fyne-cross (Docker)

Si tiene Docker instalado:

```powershell
go install github.com/fyne-io/fyne-cross@latest
fyne-cross windows -arch=amd64
```

El ejecutable se genera en `fyne-cross/dist/windows-amd64/`.

## Uso

1. **Abrir auditoría JSON**: Cargue un archivo generado por `privantix-acl-audit.exe`
2. **Crear backup**: Guarda los permisos actuales con `icacls /save`
3. **Restaurar**: Restaura permisos desde un backup con `icacls /restore`

Solo funciona en Windows (usa `icacls`).
