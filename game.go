package main

import (
	"github.com/eiannone/keyboard"
	"strconv"
	"strings"
)

func game() bool {
	if showRules() {
		return true
	}
	for {
		if board.Display() {
			ClearConsole()
			blueFont.Println("Press Esc to exit")
			board.Print()
			if board.IsComplete() {
				break
			}
		}
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyArrowUp {
			board.Move(0, -1)
		} else if key == keyboard.KeyArrowDown {
			board.Move(0, 1)
		} else if key == keyboard.KeyArrowRight {
			board.Move(1, 0)
		} else if key == keyboard.KeyArrowLeft {
			board.Move(-1, 0)
		} else if key == keyboard.KeyEsc {
			ClearConsole()
			blueFont.Println("Press Esc second time to exit or BackSpace to get back to menu or Ctrl+S to save game(any other to continue)")
			_, key, err := keyboard.GetKey()
			if err != nil {
				panic(err)
			}
			if key == keyboard.KeyEsc {
				return false
			} else if key == keyboard.KeyBackspace {
				return true
			} else if key == keyboard.KeyCtrlS {
				board.SaveGame()
				return false
			}
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

	ClearConsole()
	board.Print()

	blueFont.Println("\nCongrats on finishing sudoku! Press Backspace to get back to menu, Esc to exit")

	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		if key == keyboard.KeyEsc {
			return false
		} else if key == keyboard.KeyBackspace {
			return true
		}
	}
}

func showRules() bool {
	ClearConsole()
	blueFont.Println(board.Rules())
	blueFont.Println("Press Enter to continue. Backspace to return to menu")
	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		if key == keyboard.KeyEnter {
			return false
		} else if key == keyboard.KeyBackspace {
			return true
		}
	}
}

var board SudokuBoard

func initGame() bool {
	initBoard()
	return true
}

func initBoard() {
	boardType := gameOptions[0][gameParam[0]]
	boardSize, _ := strconv.Atoi(strings.Split(gameOptions[1][gameParam[1]], "x")[0])
	time := -1
	if gameOptions[3][gameParam[3]] != "∞" {
		time, _ = strconv.Atoi(strings.Split(gameOptions[3][gameParam[3]], " min")[0])
		time *= 60
	}
	switch boardType {
	case "square":
		basic := &BasicSudoku{}
		basic.Init(boardSize, gameParam[2], time)
		board = basic
	case "diagonal":
		diagonal := &DiagonalSudoku{}
		diagonal.Init(9, gameParam[2], time)
		board = diagonal
	case "twodoku":
		twodoku := &TwoDoku{}
		twodoku.Init(9, gameParam[2], time)
		board = twodoku
	}
}
