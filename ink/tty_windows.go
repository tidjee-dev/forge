//go:build windows

package ink

import (
	"syscall"
	"unsafe"
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
)

// isTerminalFd reports whether the file descriptor fd refers to a terminal on
// Windows. It calls GetConsoleMode, which succeeds only when the handle is a
// real console (not a pipe or file redirection).
func isTerminalFd(fd uintptr) bool {
	handle := syscall.Handle(fd)
	var mode uint32
	r, _, _ := procGetConsoleMode.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&mode)),
	)
	return r != 0
}
