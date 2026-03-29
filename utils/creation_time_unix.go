//go:build !windows

package utils

import (
	"os"
	"time"
)

// GetFileCreationTime returns the file creation time for the given path.
// On Unix, creation time is not always available from the kernel; ModTime is returned as fallback.
func GetFileCreationTime(path string) (time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}
