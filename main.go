package main

import (
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
	defer keyboard.Close()
	menu()
	//game()
}
