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

type UNICODE_STRING struct {
	Length, MaximumLength uint16
	Buffer                uintptr
}

type SYSTEM_PROCESS_INFORMATION struct {
	NextEntryOffset, NumberOfThreads  uint32
	Reserved                          [3]uint64
	CreateTime, UserTime, KernelTime  uint64
	ImageName                         UNICODE_STRING
	BasePriority                      uint32
	ProcessId, InheritedFromProcessId uintptr
}

type TOKEN_PRIVILEGES struct {
	PrivilegeCount uint32
	Privileges     [1]LUID_AND_ATTRIBUTES
}

type LUID_AND_ATTRIBUTES struct {
	Luid       LUID
	Attributes uint32
}

type LUID struct {
	LowPart  uint32
	HighPart int32
}

const (
	TRUE          = uintptr(1)
	FALSE         = uintptr(0)
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

	// ntdll
	SystemBasicInformation                = 0
	SystemPerformanceInformation          = 2
	SystemTimeOfDayInformation            = 3
	SystemProcessInformation              = 5
	SystemProcessorPerformanceInformation = 8
	SystemInterruptInformation            = 23
	SystemExceptionInformation            = 33
	SystemRegistryQuotaInformation        = 37
	SystemLookasideInformation            = 45
	SystemCodeIntegrityInformation        = 103
	SystemPolicyInformation               = 134

	// advapi32
	SE_BACKUP_NAME          = "SeBackupPrivilege"
	SE_PRIVILEGE_ENABLED    = uint32(2)
	TOKEN_ADJUST_PRIVILEGES = uint32(0x00000020)
	TOKEN_QUERY             = uint32(0x00000008)
	REG_STANDARD_FORMAT     = uint32(1)
	REG_LATEST_FORMAT       = uint32(2)
	REG_NO_COMPRESSION      = uint32(4)
)

var (
	// Kernel32
	modKernel32              = syscall.NewLazyDLL("kernel32.dll")
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

	// ntdll
	modNtdll                     = syscall.NewLazyDLL("ntdll.dll")
	ProcNtQuerySystemInformation = modNtdll.NewProc("NtQuerySystemInformation")

	// user32
	modUser32       = syscall.NewLazyDLL("user32.dll")
	ProcMessageBoxA = modUser32.NewProc("MessageBoxA")
	ProcMessageBoxW = modUser32.NewProc("MessageBoxW")

	// advapi32
	modAdvapi32               = syscall.NewLazyDLL("advapi32.dll")
	ProcAdjustTokenPrivileges = modAdvapi32.NewProc("AdjustTokenPrivileges")
	ProcLookupPrivilegeValueA = modAdvapi32.NewProc("LookupPrivilegeValueA")
	ProcOpenProcessToken      = modAdvapi32.NewProc("OpenProcessToken")
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
