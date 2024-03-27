package main

import (
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"runtime"
)

func DisableCursor() bool {
	if runtime.GOOS == "windows" {
		stdout := windows.Handle(os.Stdout.Fd())
		var originalMode uint32
		// enable ANSI escape codes
		windows.GetConsoleMode(stdout, &originalMode)
		windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
		// ANSI escape code
		fmt.Print("\033[?25l")

	} else if runtime.GOOS == "linux" {
		fmt.Print("\033[?25l")
	} else {
		return false
	}
	return true
}
