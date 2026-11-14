package main

import (
	"os/exec"
	"syscall"
	"unsafe"
)

func configureServer(command *exec.Cmd) { command.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} }
func showError(err error) {
	syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(err.Error()))), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("象映笔记"))), 0x10)
}
