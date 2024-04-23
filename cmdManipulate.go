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

// ClearConsole clear console
func ClearConsole() {
	// move to top left corner of console
	screen.MoveTopLeft()
	windth, height := screen.Size()
	// override all text with spaces
	fmt.Print(strings.Repeat(strings.Repeat(" ", windth)+"\n", height))
	screen.MoveTopLeft()
}

// DisableCursor make cursor invisible (supports linux and windows)
func DisableCursor() {
	if runtime.GOOS == "windows" {
		stdout := windows.Handle(os.Stdout.Fd())
		var originalMode uint32
		// enable ANSI escape codes
		windows.GetConsoleMode(stdout, &originalMode)
		windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	}
	// ANSI escape code to disable cursor
	fmt.Print("\033[?25l")
}

// EnableCursor make terminal cursor visible (supports linux and windows)
func EnableCursor() {
	// ANSI escape code to enable cursor
	fmt.Print("\033[?25h")
	ClearConsole()
}
