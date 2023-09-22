//go:build windows
// +build windows

// Dump the SAM and SECURITY from registry
package main

const (
	SAM_KEY      = win32.HKLM + "\\SAM"
	SECURITY_KEY = win32.HKLM + "\\SECURITY"
)

func main() {

}
