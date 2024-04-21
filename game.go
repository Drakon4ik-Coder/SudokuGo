package main

import (
	"github.com/eiannone/keyboard"
	"time"
)

// timer thread function
func timeControl(exit *bool) {
	for {
		// wait 1 second and decrement board play time if game in not on pause
		time.Sleep(time.Second)
		if !*exit {
			board.TimePass(1)
		}
	}
}

// print thread function
func printSudo(exit *bool) {
	// maximum fps
	fps := time.Duration(100)
	for {
		// print if game in not on pause and there is info to display
		time.Sleep(time.Second / fps)
		if !*exit && board.Display() {
			ClearConsole()
			blueFont.Println("Press Esc to exit or pause")
			board.Print()
		}
	}
}

// key reader thread function
func readKeys(r *rune, k *keyboard.Key, b *bool) {
	for {
		// wait for rune and key
		*r, *k, _ = keyboard.GetKey()
		// bool key found
		*b = true
	}
}

// char, key and whether it was gotten
var char rune
var key keyboard.Key
var keyBool = false

// main game function
func game() bool {
	if showRules() {
		return true
	}
	// game on pause
	pause := false
	go timeControl(&pause)
	go printSudo(&pause)
	// game loop
	for {
		// if timer has ended or board is finished
		if board.IsComplete() || board.TimeEnd() {
			break
		}

		if keyBool {
			keyBool = false
			// move with arrows
			if key == keyboard.KeyArrowUp {
				board.Move(0, -1)
			} else if key == keyboard.KeyArrowDown {
				board.Move(0, 1)
			} else if key == keyboard.KeyArrowRight {
				board.Move(1, 0)
			} else if key == keyboard.KeyArrowLeft {
				board.Move(-1, 0)
			} else if key == keyboard.KeyEsc { // pause game
				pause = true
				ClearConsole()
				blueFont.Println("Press Esc second time to pause or BackSpace to get back to menu or Ctrl+S to save game(any other to continue)")
				for {
					if keyBool {
						keyBool = false
						if key == keyboard.KeyEsc {
							return false
						} else if key == keyboard.KeyBackspace {
							return true
						} else if key == keyboard.KeyCtrlS {
							_ = board.SaveGame()
							return false
						} else {
							break
						}
					}
				}
				pause = false
				board.Move(0, 0)
			} else if key == keyboard.KeyCtrlZ { // undo move
				board.Undo()
			} else if key == keyboard.KeyCtrlY { // redo move
				board.Redo()
			} else if key == keyboard.KeyCtrlR { // reveal random element
				board.RevealRandom()
			} else if '1' <= char && char <= '9' { // enter 0 to 9
				board.Enter(int(char - '0'))
			} else if 'a' <= char && char <= 'z' {
				board.Enter(int(char - 'a' + 10))
			} else if 'A' <= char && char <= 'Z' {
				board.Enter(int(char - 'A' + 10))
			}
		}
	}
	pause = true
	ClearConsole()
	board.Print()

	// decide whether the user lost or won
	if board.IsComplete() {
		greenFont.Println("\nCongrats on finishing sudoku!")
	} else {
		redFont.Println("\nSorry, you lost on time. Good luck next time ;)")
	}
	blueFont.Println("Press Backspace to get back to menu, Esc to pause")

	for {
		if keyBool {
			if key == keyboard.KeyEsc {
				return false
			} else if key == keyboard.KeyBackspace {
				return true
			}
		}
	}
}

// show board rule
func showRules() bool {
	ClearConsole()
	blueFont.Println(board.Rules())
	blueFont.Println("Press Enter to continue. Backspace to return to menu")
	for {
		if keyBool {
			keyBool = false
			if key == keyboard.KeyEnter {
				return false
			} else if key == keyboard.KeyBackspace {
				return true
			}
		}
	}
}

// game board
var board SudokuBoard
