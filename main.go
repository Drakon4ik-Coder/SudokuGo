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

	// enter game if user didn't exit
	contin := true
	for contin {
		if menu() {
			contin = game()
		} else {
			contin = false
		}
	}
	ClearConsole()
	infoFont.Print("\tThanks for choosing to play our Sudoku. May you have a blessed day :)")
	//time.Sleep(time.Second * 5)
}
