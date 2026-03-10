package ink

import (
	"io"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

type outputColorMode uint32

const (
	colorModeAuto outputColorMode = iota
	colorModeAlways
	colorModeNever
)

var globalColorMode atomic.Uint32

func SetGlobalColorMode(mode outputColorMode) {
	globalColorMode.Store(uint32(mode))
}

const (
	AutoColorMode   = colorModeAuto
	AlwaysColorMode = colorModeAlways
	NeverColorMode  = colorModeNever
)

var (
	envDisabled bool
	envForced   bool
	envOnce     sync.Once
)

func initEnv() {
	if v, ok := os.LookupEnv("NO_COLOR"); ok && v != "" {
		envDisabled = true
		return
	}
	// TERM=dumb means the terminal cannot interpret escape sequences.
	if strings.EqualFold(os.Getenv("TERM"), "dumb") {
		envDisabled = true
		return
	}
	// COLORTERM=truecolor / 24bit signals that the terminal is capable.
	ct := strings.ToLower(os.Getenv("COLORTERM"))
	if ct == "truecolor" || ct == "24bit" {
		envForced = true
	}
}

func isFd(w io.Writer) bool {
	type fder interface {
		Fd() uintptr
	}
	f, ok := w.(fder)
	if !ok {
		return false
	}
	return isTerminal(f.Fd())
}

func isTerminal(fd uintptr) bool {
	return isTerminalFd(fd)
}

func isColorEnabled(w io.Writer) bool {
	envOnce.Do(initEnv)

	switch outputColorMode(globalColorMode.Load()) {
	case colorModeAlways:
		return true
	case colorModeNever:
		return false
	case colorModeAuto:
		if envDisabled {
			return false
		}
		if envForced {
			return true
		}
	}
	return isFd(w)
}
