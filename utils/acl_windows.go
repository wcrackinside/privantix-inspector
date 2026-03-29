//go:build windows

package utils

import (
	"os/exec"
	"strings"
)

func GetFileACLs(path string) []string {
	out, err := exec.Command("icacls", path).CombinedOutput()
	if err != nil {
		return []string{}
	}
	
	lines := strings.Split(string(out), "\n")
	var results []string
	
	// icacls output format is typically:
	// filepath BUILTIN\Users:(I)(RX)
	//          NT AUTHORITY\SYSTEM:(I)(F)
	// Usually everything after the first space (or in the first line) contains the user and rights
	
	for i, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		
		// Ignore processed success messages in multiple languages
		if strings.Contains(l, "Successfully processed") || strings.Contains(l, "Procesado correctamente") || strings.Contains(l, "Successfully") {
			continue
		}
		
		if i == 0 {
			// First line includes the file path. Let's try to remove it
			// It's the path followed by a space, then the user
			if idx := strings.Index(l, " "); idx != -1 && strings.HasPrefix(l, path) {
				l = strings.TrimSpace(l[len(path):])
			} else {
				// Path might be quoted or different
				parts := strings.SplitN(l, " ", 2)
				if len(parts) == 2 {
					l = strings.TrimSpace(parts[1])
				}
			}
		}
		if l != "" {
			results = append(results, l)
		}
	}
	
	return results
}
