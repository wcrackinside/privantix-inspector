package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var (
	advapi32                 = syscall.NewLazyDLL("advapi32.dll")
	procGetNamedSecurityInfo = advapi32.NewProc("GetNamedSecurityInfoW")
	procLookupAccountSid     = advapi32.NewProc("LookupAccountSidW")
)

func main() {
	path := "c:\\03_OtrosSistemas\\privantix-source-inspector\\go.mod"
	info, _ := os.Stat(path)
	perms := info.Mode().Perm().String()
	fmt.Println("Perms:", perms)

	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		fmt.Println("Err UTF16:", err)
		return
	}

	var ownerSid uintptr
	var secDesc uintptr

	ret, _, err2 := procGetNamedSecurityInfo.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		1, // SE_FILE_OBJECT
		1, // OWNER_SECURITY_INFORMATION
		uintptr(unsafe.Pointer(&ownerSid)),
		0, 0, 0,
		uintptr(unsafe.Pointer(&secDesc)),
	)
	if ret != 0 {
		fmt.Println("Err GetNamed:", ret, err2)
		return
	}
	defer syscall.LocalFree(syscall.Handle(secDesc))

	var name [256]uint16
	var domain [256]uint16
	nameLen := uint32(256)
	domainLen := uint32(256)
	var sidUse uint32

	ret2, _, err3 := procLookupAccountSid.Call(
		0,
		ownerSid,
		uintptr(unsafe.Pointer(&name[0])),
		uintptr(unsafe.Pointer(&nameLen)),
		uintptr(unsafe.Pointer(&domain[0])),
		uintptr(unsafe.Pointer(&domainLen)),
		uintptr(unsafe.Pointer(&sidUse)),
	)
	if ret2 == 0 {
		fmt.Println("Err Lookup:", ret2, err3)
		return
	}

	ownerName := syscall.UTF16ToString(name[:nameLen])
	domainName := syscall.UTF16ToString(domain[:domainLen])

	fmt.Println("Owner:", domainName+"\\"+ownerName)
}
