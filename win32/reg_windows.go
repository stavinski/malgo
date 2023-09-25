//go:build windows
// +build windows

// Registry win32 api calls
package win32

const (
	HKEY_CLASSES_ROOT   = 0x80000000
	HKEY_CURRENT_USER   = 0x80000001
	HKEY_LOCAL_MACHINE  = 0x80000002
	HKEY_USERS          = 0x80000003
	HKEY_CURRENT_CONFIG = 0x80000005

	KEY_ALL_ACCESS = 0xf003f
	KEY_READ       = 0x20019
)

var (
	ProcRegGetValueA  = modAdvapi32.NewProc("RegGetValueA")
	ProcRegSaveKeyA   = modAdvapi32.NewProc("RegSaveKeyA")
	ProcRegOpenKeyExA = modAdvapi32.NewProc("RegOpenKeyExA")
	ProcRegCloseKey   = modAdvapi32.NewProc("RegCloseKey")
	ProcRegSaveKeyExA = modAdvapi32.NewProc("RegSaveKeyExA")
)
