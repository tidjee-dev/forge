//go:build !windows

package ink

import (
	"syscall"
	"unsafe"
)

// isTerminalFd reports whether the file descriptor fd refers to a terminal on
// Unix-like operating systems. It uses the TCGETS ioctl (equivalent to
// tcgetattr) which succeeds only when fd is a TTY.
func isTerminalFd(fd uintptr) bool {
	var termios syscall.Termios
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		syscall.TCGETS,
		uintptr(unsafe.Pointer(&termios)),
	)
	return errno == 0
}
