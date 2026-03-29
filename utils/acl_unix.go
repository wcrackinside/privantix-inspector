//go:build !windows

package utils

func GetFileACLs(path string) []string {
	// Unix systems typically manage access via basic owner/group/other POSIX perms
	// If setfacl is strictly required, we'd wrap it here. For MVP, return empty.
	return []string{}
}
