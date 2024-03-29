package main

import (
	"fmt"
	"github.com/inancgumus/screen"
	"golang.org/x/sys/windows"
	"os"
	"runtime"
	"strings"
)

func init() {
	screen.Clear()
	screen.MoveTopLeft()
}

func ClearConsole() {
	screen.MoveTopLeft()
	windth, height := screen.Size()
	fmt.Print(strings.Repeat(strings.Repeat(" ", windth)+"\n", height))
	screen.MoveTopLeft()
}

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

// make terminal cursor visible
func EnableCursor() bool {
	if runtime.GOOS == "windows" {
		fmt.Print("\033[?25h")

	} else if runtime.GOOS == "linux" {
		fmt.Print("\033[?25h")
	} else {
		ClearConsole()
		return false
	}
	ClearConsole()
	return true
}
