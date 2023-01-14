//go:build windows
// +build windows

package win32

import "syscall"

const (
	WLAN_MAX_NAME_LENGTH           = uint32(256)
	WLAN_PROFILE_GET_PLAINTEXT_KEY = uint32(4)
)

var (
	modWlanapi = syscall.NewLazyDLL("Wlanapi.dll")

	ProcWlanOpenHandle     = modWlanapi.NewProc("WlanOpenHandle")
	ProcWlanCloseHandle    = modWlanapi.NewProc("WlanCloseHandle")
	ProcWlanEnumInterfaces = modWlanapi.NewProc("WlanEnumInterfaces")
	ProcWlanFreeMemory     = modWlanapi.NewProc("WlanFreeMemory")
	ProcWlanGetProfileList = modWlanapi.NewProc("WlanGetProfileList")
	ProcWlanGetProfile     = modWlanapi.NewProc("WlanGetProfile")
)

type WLAN_PROFILE_INFO_LIST struct {
	NumberOfItems, Index uint32
	ProfileInfo          [1]WLAN_PROFILE_INFO
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
	InterfaceInfo        [1]WLAN_INTERFACE_INFO
}
