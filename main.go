package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
)

func frame() bool {
	return true
}

func game() bool {
	frameContinue := frame()
	for frameContinue {
		frameContinue = frame()
		ClearConsole()
	}
	return true
}

func main() {
	// disable console cursor
	fmt.Print("\033[?25l")
	defer keyboard.Close()
	menu()
	//game()
}
