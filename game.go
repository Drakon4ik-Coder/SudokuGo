package main

import (
	"github.com/eiannone/keyboard"
	"strconv"
	"strings"
)

func game() bool {
	initGame()
	chosenPos := Vector2{0, 0}
	posChange := false
	boardSize := board.GetSize()
	for !board.IsComplete() {
		if board.Display() || posChange {
			posChange = false
			ClearConsole()
			board.Print(chosenPos.width, chosenPos.height)
		}
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		posChange = true

		if key == keyboard.KeyArrowUp && chosenPos.width > 0 {
			chosenPos.width--
		} else if key == keyboard.KeyArrowDown && chosenPos.width < boardSize-1 {
			chosenPos.width++
		} else if key == keyboard.KeyArrowRight && chosenPos.height < boardSize-1 {
			chosenPos.height++
		} else if key == keyboard.KeyArrowLeft && chosenPos.height > 0 {
			chosenPos.height--
		} else if key == keyboard.KeyEsc {
			ClearConsole()
			infoFont.Println("Press Esc second time to exit or BackSpace to get back to menu(any other to continue)")
			_, key, err := keyboard.GetKey()
			if err != nil {
				panic(err)
			}
			if key == keyboard.KeyEsc {
				return false
			} else if key == keyboard.KeyBackspace {
				return true
			}
		} else if '1' <= char && char <= '9' && int(char-'0') <= boardSize {
			board.Enter(chosenPos.width, chosenPos.height, int(char-'0'))
		} else if 'a' <= char && char <= 'z' && int(char-'a'+10) <= boardSize {
			board.Enter(chosenPos.width, chosenPos.height, int(char-'a'+10))
		} else if 'A' <= char && char <= 'Z' && int(char-'A'+10) <= boardSize {
			board.Enter(chosenPos.width, chosenPos.height, int(char-'A'+10))
		} else {
			posChange = false
		}
	}

	ClearConsole()
	board.Print(-1, -1)

	infoFont.Println("\nCongrats on finishing sudoku! Press Backspace to get back to menu, Esc to exit")

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

var board SudokuBoard

func initGame() bool {
	initBoard()
	return true
}

func initBoard() {
	boardType := gameOptions[0][gameParam[0]]
	boardSize, _ := strconv.Atoi(strings.Split(gameOptions[1][gameParam[1]], "x")[0])
	switch boardType {
	case "square":
		basic := &BasicSudoku{}
		basic.Init(boardSize, gameParam[2])
		board = basic
	}
}
