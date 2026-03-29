//go:build !windows

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Privantix ACL Manager requiere Windows (usa icacls para backup/restore de permisos).")
	os.Exit(1)
}
