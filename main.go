package main

import (
	"github.com/eiannone/keyboard"
)

func main() {
	// close keyboard output
	defer keyboard.Close()
	// enable cursor after program finish
	defer EnableCursor()
	// thread to wait for keys
	go readKeys(&char, &key, &keyBool)
	// disable console cursor
	DisableCursor()

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
