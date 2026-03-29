//go:build windows

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type DirACL struct {
	Path              string   `json:"path"`
	Name              string   `json:"name"`
	Depth             int      `json:"depth"`
	Owner             string   `json:"owner"`
	Permissions       string   `json:"permissions"`
	ACLs              []string `json:"acls"`
	TrustedPrincipals []string `json:"trusted_principals,omitempty"`
	OutsidePrincipals []string `json:"outside_principals,omitempty"`
	ComplianceStatus  string   `json:"compliance_status,omitempty"`
}

type AuditResult struct {
	Path   string   `json:"path"`
	Dirs   []DirACL `json:"dirs"`
	Errors []string `json:"errors,omitempty"`
}

func main() {
	a := app.NewWithID("privantix.acl.manager")
	w := a.NewWindow("Privantix ACL Manager - Gestión de permisos")
	w.Resize(fyne.NewSize(1000, 650))

	state := &appState{audit: nil}

	// Tabla de directorios
	table := widget.NewTable(
		func() (int, int) {
			if state.audit == nil {
				return 0, 0
			}
			return len(state.audit.Dirs) + 1, 5
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			if state.audit == nil {
				label.SetText("")
				return
			}
			if id.Row == 0 {
				headers := []string{"Ruta", "Nombre", "Propietario", "Permisos", "Cumplimiento"}
				if id.Col < len(headers) {
					label.SetText(headers[id.Col])
				}
				label.TextStyle = fyne.TextStyle{Bold: true}
				return
			}
			idx := id.Row - 1
			if idx >= len(state.audit.Dirs) {
				return
			}
			d := state.audit.Dirs[idx]
			switch id.Col {
			case 0:
				label.SetText(d.Path)
			case 1:
				label.SetText(d.Name)
			case 2:
				label.SetText(d.Owner)
			case 3:
				label.SetText(d.Permissions)
			case 4:
				label.SetText(d.ComplianceStatus)
			default:
				label.SetText("")
			}
		},
	)
	table.SetColumnWidth(0, 350)
	table.SetColumnWidth(1, 120)
	table.SetColumnWidth(2, 180)
	table.SetColumnWidth(3, 100)
	table.SetColumnWidth(4, 110)

	// Panel de detalle ACLs
	detailLabel := widget.NewLabel("Seleccione un directorio para ver los ACLs")
	detailLabel.Wrapping = fyne.TextWrapWord
	detailScroll := container.NewScroll(detailLabel)
	detailScroll.SetMinSize(fyne.NewSize(0, 140))

	table.OnSelected = func(id widget.TableCellID) {
		if state.audit == nil || id.Row == 0 || id.Row > len(state.audit.Dirs) {
			return
		}
		d := state.audit.Dirs[id.Row-1]
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Ruta: %s\n", d.Path))
		sb.WriteString(fmt.Sprintf("Propietario: %s\n", d.Owner))
		sb.WriteString(fmt.Sprintf("Permisos: %s\n\n", d.Permissions))
		sb.WriteString("ACLs:\n")
		for _, acl := range d.ACLs {
			if strings.Contains(acl, "procesaron") || strings.Contains(acl, "processed") {
				continue
			}
			sb.WriteString("  • ")
			sb.WriteString(acl)
			sb.WriteString("\n")
		}
		if len(d.TrustedPrincipals) > 0 || len(d.OutsidePrincipals) > 0 {
			sb.WriteString("\n")
			if len(d.TrustedPrincipals) > 0 {
				sb.WriteString("Confianza: " + strings.Join(d.TrustedPrincipals, ", ") + "\n")
			}
			if len(d.OutsidePrincipals) > 0 {
				sb.WriteString("Externos: " + strings.Join(d.OutsidePrincipals, ", ") + "\n")
			}
		}
		detailLabel.SetText(sb.String())
	}

	// Botones
	openBtn := widget.NewButton("Abrir auditoría JSON", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()
			data, err := os.ReadFile(uriToPath(reader.URI()))
			if err != nil {
				dialog.ShowError(fmt.Errorf("error leyendo archivo: %w", err), w)
				return
			}
			var audit AuditResult
			if err := json.Unmarshal(data, &audit); err != nil {
				dialog.ShowError(fmt.Errorf("JSON inválido: %w", err), w)
				return
			}
			state.audit = &audit
			state.loadedPath = reader.URI().Path()
			table.Refresh()
			w.SetTitle(fmt.Sprintf("Privantix ACL Manager - %s (%d directorios)", filepath.Base(state.loadedPath), len(audit.Dirs)))
		}, w)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		fd.Resize(fyne.NewSize(700, 500))
		fd.Show()
	})

	createBackupBtn := widget.NewButton("Crear backup de permisos", func() {
		if state.audit == nil {
			dialog.ShowInformation("Sin datos", "Abra primero un archivo de auditoría JSON.", w)
			return
		}
		rootPath := state.audit.Path
		fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil || writer == nil {
				return
			}
			defer writer.Close()
			backupPath := uriToPath(writer.URI())
			if !strings.HasSuffix(strings.ToLower(backupPath), ".txt") {
				backupPath += ".txt"
			}
			cmd := exec.Command("icacls", rootPath, "/save", backupPath, "/t", "/c")
			out, err := cmd.CombinedOutput()
			if err != nil {
				dialog.ShowError(fmt.Errorf("error creando backup: %v\n%s", err, string(out)), w)
				return
			}
			dialog.ShowInformation("Backup creado", fmt.Sprintf("Permisos guardados en:\n%s", backupPath), w)
		}, w)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
		fd.SetFileName("permisos_backup.txt")
		fd.Resize(fyne.NewSize(700, 500))
		fd.Show()
	})

	restoreBtn := widget.NewButton("Restaurar desde backup", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			reader.Close()
			backupPath := uriToPath(reader.URI())
			dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
				if err != nil || uri == nil {
					return
				}
				rootPath := uriToPath(uri)
				cmd := exec.Command("icacls", rootPath, "/restore", backupPath)
				out, err := cmd.CombinedOutput()
				if err != nil {
					dialog.ShowError(fmt.Errorf("error restaurando: %v\n%s", err, string(out)), w)
					return
				}
				dialog.ShowInformation("Restauración completada", "Los permisos se restauraron correctamente.", w)
			}, w)
		}, w)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
		fd.Resize(fyne.NewSize(700, 500))
		fd.Show()
	})

	toolbar := container.NewHBox(openBtn, createBackupBtn, restoreBtn)

	split := container.NewVSplit(table, detailScroll)
	split.SetOffset(0.72)

	mainContent := container.NewBorder(toolbar, nil, nil, nil, split)
	w.SetContent(mainContent)
	w.ShowAndRun()
}

type appState struct {
	audit      *AuditResult
	loadedPath string
}

func uriToPath(u fyne.URI) string {
	if u == nil {
		return ""
	}
	path := u.Path()
	if path == "" {
		path = u.String()
	}
	path = strings.TrimPrefix(path, "file://")
	path = strings.TrimPrefix(path, "file:///")
	if len(path) > 2 && path[0] == '/' && path[2] == ':' {
		path = path[1:]
	}
	return filepath.FromSlash(path)
}
