//go:build windows

package utils

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	advapi32                 = syscall.NewLazyDLL("advapi32.dll")
	procGetNamedSecurityInfo = advapi32.NewProc("GetNamedSecurityInfoW")
	procLookupAccountSid     = advapi32.NewProc("LookupAccountSidW")
)

func GetFileOwnerAndPerms(info os.FileInfo, path string) (owner string, perms string) {
	perms = info.Mode().Perm().String()
	owner = "Unknown"

	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return
	}

	var ownerSid uintptr
	var secDesc uintptr

	ret, _, _ := procGetNamedSecurityInfo.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		1, // SE_FILE_OBJECT
		1, // OWNER_SECURITY_INFORMATION
		uintptr(unsafe.Pointer(&ownerSid)),
		0, 0, 0,
		uintptr(unsafe.Pointer(&secDesc)),
	)
	if ret != 0 {
		return
	}
	defer syscall.LocalFree(syscall.Handle(secDesc))

	var name [256]uint16
	var domain [256]uint16
	nameLen := uint32(256)
	domainLen := uint32(256)
	var sidUse uint32

	ret2, _, _ := procLookupAccountSid.Call(
		0,
		ownerSid,
		uintptr(unsafe.Pointer(&name[0])),
		uintptr(unsafe.Pointer(&nameLen)),
		uintptr(unsafe.Pointer(&domain[0])),
		uintptr(unsafe.Pointer(&domainLen)),
		uintptr(unsafe.Pointer(&sidUse)),
	)
	if ret2 != 0 {
		ownerName := syscall.UTF16ToString(name[:nameLen])
		domainName := syscall.UTF16ToString(domain[:domainLen])
		if domainName != "" {
			owner = domainName + "\\" + ownerName
		} else {
			owner = ownerName
		}
	}

	return
}
