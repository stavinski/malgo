//go:build windows
// +build windows

// Common types and function calls for the Win32 API
package win32

import (
	"syscall"
	"unicode/utf16"
	"unsafe"
)

type GUID struct {
	Data1        uint32
	Data2, Data3 uint16
	Data4        [8]byte
}

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
)

var (
	modKernel32 = syscall.NewLazyDLL("kernel32.dll")

	// Kernel32
	ProcOpenProcess          = modKernel32.NewProc("OpenProcess")
	ProcGetCurrentProcess    = modKernel32.NewProc("GetCurrentProcess")
	ProcCloseHandle          = modKernel32.NewProc("CloseHandle")
	ProcVirtualProtect       = modKernel32.NewProc("VirtualProtect")
	ProcVirtualAllocEx       = modKernel32.NewProc("VirtualAllocEx")
	ProcWriteProcessMemory   = modKernel32.NewProc("WriteProcessMemory")
	ProcCreateRemoteThreadEx = modKernel32.NewProc("CreateRemoteThreadEx")
	ProcResumeThread         = modKernel32.NewProc("ResumeThread")
	ProcLoadLibraryA         = modKernel32.NewProc("LoadLibraryA")
	ProcVirtualFreeEx        = modKernel32.NewProc("VirtualFreeEx")
	ProcWaitForSingleObject  = modKernel32.NewProc("WaitForSingleObject")
	ProcGetModuleHandleA     = modKernel32.NewProc("GetModuleHandleA")
)

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

// Convert byte pointer to a Go string
func BytePtrToString(p *uint8) string {
	if p == nil {
		return ""
	}
	if *p == 0 {
		return ""
	}

	// Find NUL terminator.
	n := 0
	for ptr := unsafe.Pointer(p); *(*uint8)(ptr) != 0; n++ {
		ptr = unsafe.Pointer(uintptr(ptr) + unsafe.Sizeof(*p))
	}

	return string(unsafe.Slice(p, n))
}
