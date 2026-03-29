//go:build !windows

package utils

import (
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func GetFileOwnerAndPerms(info os.FileInfo, path string) (owner string, perms string) {
	perms = info.Mode().Perm().String()
	owner = "Unknown"

	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		uidStr := strconv.FormatUint(uint64(stat.Uid), 10)
		if u, err := user.LookupId(uidStr); err == nil {
			owner = u.Username
		} else {
			owner = uidStr
		}
	}
	return
}
