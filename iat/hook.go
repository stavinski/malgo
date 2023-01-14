//go:build windows
// +build windows

package main

import "C"
import (
	"syscall"
	"unsafe"
)

//export HookedMessageBoxAFunc
func HookedMessageBoxAFunc(hWnd C.uint, lpText, lpCaption *C.char, uType C.uint) C.int {
	// display hook msgbox
	OriginalMessageBoxA.Call(uintptr(hWnd), uintptr(unsafe.Pointer(syscall.StringBytePtr("Hooked Call"))), uintptr(unsafe.Pointer(syscall.StringBytePtr("Hook"))), uintptr(uType))

	// call original
	res, _, _ := OriginalMessageBoxA.Call(uintptr(hWnd), uintptr(unsafe.Pointer(lpText)), uintptr(unsafe.Pointer(lpCaption)), uintptr(uType))
	return C.int(res)
}
