//go:build windows
// +build windows

// Retrieve WLAN profile passwords from Windows using Win32 API calls
package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/stavinski/malgo/win32"
)

// XML marshal type for a WLAN profile
type ProfileXML struct {
	// TBC
}

func main() {
	var (
		ver           uint32
		negotiatedVer uint32
		wlanHandle    uintptr
		res           uintptr
		err           error
	)

	ver = 2
	res, _, err = win32.ProcWlanOpenHandle.Call(uintptr(ver),
		win32.NULL,
		uintptr(unsafe.Pointer(&negotiatedVer)),
		uintptr(unsafe.Pointer(&wlanHandle)))
	if res != win32.ERROR_SUCCESS {
		if res == 1062 {
			// WLAN service not running
			return
		}
		fmt.Fprint(os.Stderr, err.Error())
		return
	}
	defer win32.ProcWlanCloseHandle.Call(wlanHandle, win32.NULL)
	interfaces := &win32.WLAN_INTERFACE_INFO_LIST{}
	res, _, err = win32.ProcWlanEnumInterfaces.Call(wlanHandle, win32.NULL, uintptr(unsafe.Pointer(&interfaces)))
	if res != win32.ERROR_SUCCESS {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}
	// release interfaces memory at the end of the func
	defer win32.ProcWlanFreeMemory.Call(uintptr(unsafe.Pointer(interfaces)))

	fmt.Printf("%v ifaces found.\n", interfaces.NumberOfItems)

	// for i := uint32(0); i < interfaces.NumberOfItems; i++ {
	// 	var profList *win32.WLAN_PROFILE_INFO_LIST
	// 	iface := interfaces.InterfaceInfo[i]
	// 	res, _, err = win32.ProcWlanGetProfileList.Call(wlanHandle,
	// 		uintptr(unsafe.Pointer(&iface.InterfaceGUID)),
	// 		win32.NULL,
	// 		uintptr(unsafe.Pointer(&profList)))
	// 	if res != win32.ERROR_SUCCESS {
	// 		fmt.Fprint(os.Stderr, err.Error())
	// 		return
	// 	}

	// 	for j := uint32(0); j < profList.NumberOfItems; j++ {
	// 		var profXMLPtr uintptr
	// 		profListEntry := profList.ProfileInfo[j]
	// 		fmt.Printf("Profile: %v\n", syscall.UTF16ToString(profListEntry.ProfileName[:]))
	// 		flags := win32.WLAN_PROFILE_GET_PLAINTEXT_KEY
	// 		res, _, err = win32.ProcWlanGetProfile.Call(wlanHandle,
	// 			uintptr(unsafe.Pointer(&iface.InterfaceGUID)),
	// 			uintptr(unsafe.Pointer(&profListEntry.ProfileName)),
	// 			win32.NULL,
	// 			profXMLPtr,
	// 			uintptr(unsafe.Pointer(&flags)),
	// 			win32.NULL)
	// 		if res != win32.ERROR_SUCCESS {
	// 			fmt.Fprint(os.Stderr, err.Error())
	// 			return
	// 		}

	// 		profXML := win32.UTF16PtrToString((*uint16)(unsafe.Pointer(profXMLPtr)))
	// 		fmt.Println(profXML)

	// 		// release XML string memory
	// 		win32.ProcWlanFreeMemory.Call(profXMLPtr)
	// 	}

	// 	// release profList memory
	// 	win32.ProcWlanFreeMemory.Call(uintptr(unsafe.Pointer(profList)))
	// }

}
