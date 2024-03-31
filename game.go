package main

import (
	"github.com/eiannone/keyboard"
	"strconv"
	"strings"
)

func game() {
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
		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyArrowUp && chosenPos.width > 0 {
			chosenPos.width--
			posChange = true
		} else if key == keyboard.KeyArrowDown && chosenPos.width < boardSize-1 {
			chosenPos.width++
			posChange = true
		} else if key == keyboard.KeyArrowRight && chosenPos.height < boardSize-1 {
			chosenPos.height++
			posChange = true
		} else if key == keyboard.KeyArrowLeft && chosenPos.height > 0 {
			chosenPos.height--
			posChange = true
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
