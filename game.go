package main

import (
	"github.com/eiannone/keyboard"
	"time"
)

func timeControl(exit *bool) {
	for {
		time.Sleep(time.Second)
		if !*exit {
			board.TimePass(1)
		}
	}
}

func printSudo(exit *bool) {
	for {
		time.Sleep(time.Second / 70)
		if !*exit {
			if board.Display() {
				ClearConsole()
				blueFont.Println("Press Esc to exit or pause")
				board.Print()
			}
		}
	}
}

func readKeys(r *rune, k *keyboard.Key, b *bool) {
	for {
		*r, *k, _ = keyboard.GetKey()
		*b = true
	}
}

var char rune
var key keyboard.Key
var keyBool = false

func game() bool {
	if showRules() {
		return true
	}
	exit := false
	go timeControl(&exit)
	go printSudo(&exit)
	for {
		if board.IsComplete() || board.TimeEnd() {
			break
		}

		if keyBool {
			keyBool = false
			if key == keyboard.KeyArrowUp {
				board.Move(0, -1)
			} else if key == keyboard.KeyArrowDown {
				board.Move(0, 1)
			} else if key == keyboard.KeyArrowRight {
				board.Move(1, 0)
			} else if key == keyboard.KeyArrowLeft {
				board.Move(-1, 0)
			} else if key == keyboard.KeyEsc {
				exit = true
				ClearConsole()
				blueFont.Println("Press Esc second time to exit or BackSpace to get back to menu or Ctrl+S to save game(any other to continue)")
				for !keyBool {
				}
				if key == keyboard.KeyEsc {
					return false
				} else if key == keyboard.KeyBackspace {
					return true
				} else if key == keyboard.KeyCtrlS {
					_ = board.SaveGame()
					return false
				}
				exit = false
				board.Move(0, 0)
			} else if key == keyboard.KeyCtrlZ {
				board.Undo()
			} else if key == keyboard.KeyCtrlY {
				board.Redo()
			} else if key == keyboard.KeyCtrlR {
				board.RevealRandom()
			} else if '1' <= char && char <= '9' {
				board.Enter(int(char - '0'))
			} else if 'a' <= char && char <= 'z' {
				board.Enter(int(char - 'a' + 10))
			} else if 'A' <= char && char <= 'Z' {
				board.Enter(int(char - 'A' + 10))
			}
		}
	}
	exit = true
	ClearConsole()
	board.Print()

	if board.IsComplete() {
		greenFont.Println("\nCongrats on finishing sudoku!")
	} else {
		redFont.Println("\nSorry, you lost on time. Good luck next time ;)")
	}
	blueFont.Println("Press Backspace to get back to menu, Esc to exit")

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

func showRules() bool {
	ClearConsole()
	blueFont.Println(board.Rules())
	blueFont.Println("Press Enter to continue. Backspace to return to menu")
	for {
		_, key, _ := keyboard.GetKey()
		if key == keyboard.KeyEnter {
			return false
		} else if key == keyboard.KeyBackspace {
			return true
		}
	}
}

var board SudokuBoard
