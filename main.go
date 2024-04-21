package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"os"
)

func main() {
	// close keyboard output
	defer keyboard.Close()
	// enable cursor after program finish
	defer EnableCursor()
	// thread to wait for keys
	go readKeys(&char, &key, &keyBool)
	// disable console cursor
	if DisableCursor() {
		fmt.Print("Unfortunately we cannot disable your cursor blink. Press Enter to continue and Ctrl+Q to exit")
		for {
			if keyBool {
				if key == keyboard.KeyEnter {
					break
				} else if key == keyboard.KeyCtrlQ {
					os.Exit(101)
				}
			}
		}
	}

	// program loop
	contin := true
	for contin {
		if menu() {
			// enter game if user didn't exit
			contin = game()
		} else {
			contin = false
		}
	}
	ClearConsole()
	blueFont.Print("\tThanks for choosing to play our Sudoku. May you have a blessed day :)")
}
