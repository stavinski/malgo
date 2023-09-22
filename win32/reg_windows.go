//go:build windows
// +build windows

// Registry win32 api calls
package win32

import "syscall"

const (
	HKLM = "HKEY_LOCAL_MACHINE"
	HKCU = "HKEY_CURRENT_USER"
)

var (
	modAdvapi32 = syscall.NewLazyDLL("Advapi32.dll")

	ProcRegGetValueA = modAdvapi32.MustFindProc("RegGetValueA")
	ProcRegSaveKeyA  = modAdvapi32.MustFindProc("RegSaveKeyA")
)
