//go:build windows
// +build windows

// Common types and function calls for the Win32 API
package win32

import (
	"syscall"
	"unicode/utf16"
	"unsafe"
)

const (
	NULL          = uintptr(0)
	ERROR_SUCCESS = uintptr(0)

	// Kernel32
	PROCESS_CREATE_THREAD     = uint32(2)
	PROCESS_QUERY_INFORMATION = uint32(0x0400)
	PROCESS_VM_OPERATION      = uint32(8)
	PROCESS_VM_READ           = uint32(10)
	PROCESS_VM_WRITE          = uint32(0x0020)
	MEM_COMMIT                = uint32(0x00001000)
	MEM_RESERVE               = uint32(0x00002000)
	MEM_RELEASE               = uint32(0x00008000)
	PAGE_EXECUTE_READWRITE    = uint32(0x40)
	PAGE_EXECUTE              = uint32(10)
	PAGE_EXECUTE_READ         = uint32(0x20)
	PAGE_READWRITE            = uint32(4)
	INFINITE                  = uint32(0xFFFFFFFF)
	THREAD_CREATE_SUSPENDED   = uint32(4)

	// WLAN
	WLAN_MAX_NAME_LENGTH           = uint32(256)
	WLAN_PROFILE_GET_PLAINTEXT_KEY = uint32(4)
)

var (
	modWlanapi  = syscall.NewLazyDLL("Wlanapi.dll")
	modKernel32 = syscall.NewLazyDLL("kernel32.dll")

	// WLAN
	ProcWlanOpenHandle     = modWlanapi.NewProc("WlanOpenHandle")
	ProcWlanCloseHandle    = modWlanapi.NewProc("WlanCloseHandle")
	ProcWlanEnumInterfaces = modWlanapi.NewProc("WlanEnumInterfaces")
	ProcWlanFreeMemory     = modWlanapi.NewProc("WlanFreeMemory")
	ProcWlanGetProfileList = modWlanapi.NewProc("WlanGetProfileList")
	ProcWlanGetProfile     = modWlanapi.NewProc("WlanGetProfile")

	// Kernel32
	ProcOpenProcess          = modKernel32.NewProc("OpenProcess")
	ProcCloseHandle          = modKernel32.NewProc("CloseHandle")
	ProcVirtualProtect       = modKernel32.NewProc("VirtualProtect")
	ProcVirtualAllocEx       = modKernel32.NewProc("VirtualAllocEx")
	ProcWriteProcessMemory   = modKernel32.NewProc("WriteProcessMemory")
	ProcCreateRemoteThreadEx = modKernel32.NewProc("CreateRemoteThreadEx")
	ProcResumeThread         = modKernel32.NewProc("ResumeThread")
	ProcLoadLibraryA         = modKernel32.NewProc("LoadLibraryA")
	ProcVirtualFreeEx        = modKernel32.NewProc("VirtualFreeEx")
	ProcWaitForSingleObject  = modKernel32.NewProc("WaitForSingleObject")
)

type GUID struct {
	Data1        uint32
	Data2, Data3 uint16
	Data4        [8]byte
}

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

// Convert UTF pointer to a Go string
//
// Taken from windows package
func UTF16PtrToString(p *uint16) string {
	if p == nil {
		return ""
	}
	if *p == 0 {
		return ""
	}

	// Find NUL terminator.
	n := 0
	for ptr := unsafe.Pointer(p); *(*uint16)(ptr) != 0; n++ {
		ptr = unsafe.Pointer(uintptr(ptr) + unsafe.Sizeof(*p))
	}

	return string(utf16.Decode(unsafe.Slice(p, n)))
}
