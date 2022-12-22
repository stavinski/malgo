//go:build windows
// +build windows

// Common types and function calls for the Win32 API
package win32

import "syscall"

const (
	NULL                 = uintptr(0)
	ERROR_SUCCESS        = uintptr(0)
	WLAN_MAX_NAME_LENGTH = 256
)

var (
	modWlanapi = syscall.NewLazyDLL("Wlanapi.dll")

	// TODO: make these internal and expose a proper func for each
	ProcWlanOpenHandle     = modWlanapi.NewProc("WlanOpenHandle")
	ProcWlanCloseHandle    = modWlanapi.NewProc("procWlanCloseHandle")
	ProcWlanEnumInterfaces = modWlanapi.NewProc("WlanEnumInterfaces")
	ProcWlanFreeMemory     = modWlanapi.NewProc("WlanFreeMemory")
)

type GUID struct {
	Data1        uint32
	Data2, Data3 uint16
	Data4        [8]byte
}

type WLAN_PROFILE_INFO struct {
	ProfileName [WLAN_MAX_NAME_LENGTH]uint16
	Flags       uint32
}

type WLAN_INTERFACE_INFO struct {
	InterfaceGUID        GUID
	InterfaceDescription [WLAN_MAX_NAME_LENGTH]uint16
	InterfaceState       uint32
}

type WLAN_INTERFACE_INFO_LIST struct {
	NumberOfItems, Index uint32
	InterfaceInfo        [64]WLAN_INTERFACE_INFO
}
