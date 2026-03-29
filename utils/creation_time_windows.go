//go:build windows

package utils

import (
	"syscall"
	"time"
	"unsafe"
)

var (
	kernel32     = syscall.NewLazyDLL("kernel32.dll")
	createFile   = kernel32.NewProc("CreateFileW")
	getFileInfo  = kernel32.NewProc("GetFileInformationByHandle")
	closeHandle  = kernel32.NewProc("CloseHandle")
)

// byHandleFileInformation matches BY_HANDLE_FILE_INFORMATION
type byHandleFileInformation struct {
	FileAttributes     uint32
	CreationTime       syscall.Filetime
	LastAccessTime     syscall.Filetime
	LastWriteTime      syscall.Filetime
	VolumeSerialNumber uint32
	FileSizeHigh       uint32
	FileSizeLow        uint32
	NumberOfLinks      uint32
	FileIndexHigh      uint32
	FileIndexLow       uint32
}

const (
	fileShareRead  = 0x00000001
	openExisting   = 3
	genericRead    = 0x80000000
)

func filetimeToTime(ft syscall.Filetime) time.Time {
	nsec100 := int64(ft.HighDateTime)<<32 + int64(ft.LowDateTime)
	// 100-nanosecond intervals since 1601-01-01; convert to Unix (since 1970-01-01)
	const winEpochSec = 11644473600
	const hundredsOfNsPerSec = 10000000
	sec := nsec100/hundredsOfNsPerSec - winEpochSec
	nsec := (nsec100 % hundredsOfNsPerSec) * 100
	return time.Unix(sec, nsec)
}

// GetFileCreationTime returns the file creation time for the given path on Windows.
func GetFileCreationTime(path string) (time.Time, error) {
	p, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return time.Time{}, err
	}
	handle, _, err := createFile.Call(
		uintptr(unsafe.Pointer(p)),
		uintptr(genericRead),
		uintptr(fileShareRead),
		0,
		uintptr(openExisting),
		0,
		0,
	)
	if handle == uintptr(syscall.InvalidHandle) {
		return time.Time{}, err
	}
	defer closeHandle.Call(handle)

	var info byHandleFileInformation
	ret, _, err := getFileInfo.Call(handle, uintptr(unsafe.Pointer(&info)))
	if ret == 0 {
		return time.Time{}, err
	}
	return filetimeToTime(info.CreationTime), nil
}
