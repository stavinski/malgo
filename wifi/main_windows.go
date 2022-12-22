//go:build windows
// +build windows

package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/stavinski/malgo/win32"
)

func main() {
	var negotiatedVer uint32
	var hWlan uintptr
	res, _, err := win32.ProcWlanOpenHandle.Call(uintptr(uint32(2)), win32.NULL, uintptr(unsafe.Pointer(&negotiatedVer)), uintptr(unsafe.Pointer(&hWlan)))
	if res != win32.ERROR_SUCCESS {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}
	defer win32.ProcWlanCloseHandle.Call(uintptr(hWlan), win32.NULL)
	interfaces := &win32.WLAN_INTERFACE_INFO_LIST{}
	res, _, err = win32.ProcWlanEnumInterfaces.Call(hWlan, win32.NULL, uintptr(unsafe.Pointer(&interfaces)))
	if res != win32.ERROR_SUCCESS {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}
	defer win32.ProcWlanFreeMemory.Call(uintptr(unsafe.Pointer(interfaces)))

	fmt.Println("%v ifaces found.", interfaces.NumberOfItems)

}
