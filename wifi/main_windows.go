//go:build windows
// +build windows

// Retrieve WLAN profile passwords from Windows using Win32 API calls
package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/stavinski/malgo/win32"
)

// XML marshal type for a WLAN profile
type ProfileXML struct {
	// TBC
}

func main() {
	var (
		negotiatedVer uint32
		hWlan         uintptr
		res           uintptr
		err           error
	)

	res, _, err = win32.ProcWlanOpenHandle.Call(uintptr(uint32(2)),
		win32.NULL,
		uintptr(unsafe.Pointer(&negotiatedVer)),
		uintptr(unsafe.Pointer(&hWlan)))
	if res != win32.ERROR_SUCCESS {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}
	defer win32.ProcWlanCloseHandle.Call(hWlan, win32.NULL)
	interfaces := &win32.WLAN_INTERFACE_INFO_LIST{}
	res, _, err = win32.ProcWlanEnumInterfaces.Call(hWlan, win32.NULL, uintptr(unsafe.Pointer(&interfaces)))
	if res != win32.ERROR_SUCCESS {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}
	// release interfaces memory at the end of the func
	defer win32.ProcWlanFreeMemory.Call(uintptr(unsafe.Pointer(interfaces)))

	fmt.Printf("%v ifaces found.\n", interfaces.NumberOfItems)

	for i := uint32(0); i < interfaces.NumberOfItems; i++ {
		var profList *win32.WLAN_PROFILE_INFO_LIST
		iface := interfaces.InterfaceInfo[i]
		res, _, err = win32.ProcWlanGetProfileList.Call(hWlan,
			uintptr(unsafe.Pointer(&iface.InterfaceGUID)),
			win32.NULL,
			uintptr(unsafe.Pointer(&profList)))
		if res != win32.ERROR_SUCCESS {
			fmt.Fprintf(os.Stderr, err.Error())
			return
		}

		for j := uint32(0); j < profList.NumberOfItems; j++ {
			var profXMLPtr uintptr
			profListEntry := profList.ProfileInfo[j]
			fmt.Printf("Profile: %v\n", syscall.UTF16ToString(profListEntry.ProfileName[:]))
			flags := win32.WLAN_PROFILE_GET_PLAINTEXT_KEY
			res, _, err = win32.ProcWlanGetProfile.Call(hWlan,
				uintptr(unsafe.Pointer(&iface.InterfaceGUID)),
				uintptr(unsafe.Pointer(&profListEntry.ProfileName)),
				win32.NULL,
				uintptr(unsafe.Pointer(profXMLPtr)),
				uintptr(unsafe.Pointer(&flags)),
				win32.NULL)
			if res != win32.ERROR_SUCCESS {
				fmt.Fprintf(os.Stderr, err.Error())
				return
			}

			profXML := win32.UTF16PtrToString((*uint16)(unsafe.Pointer(profXMLPtr)))
			fmt.Println(profXML)

			// release XML string memory
			win32.ProcWlanFreeMemory.Call(profXMLPtr)
		}

		// release profList memory
		win32.ProcWlanFreeMemory.Call(uintptr(unsafe.Pointer(profList)))
	}

}
