package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"os"
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
	if !DisableCursor() {
		fmt.Print("Unfortunately we cannot disable your cursor blink. Press Enter to continue and Ctrl+Q to exit")
		for {
			_, key, err := keyboard.GetKey()
			if err != nil {
				panic(err)
			}
			if key == keyboard.KeyEnter {
				break
			} else if key == keyboard.KeyCtrlQ {
				os.Exit(101)
			}
		}
	}
	// close keyboard output
	defer keyboard.Close()

	// call menu
	if !menu() {
		fmt.Print("Thanks for choosing to play our Sudoku. May you have a blessed day :)")
	}
	//game()
}
