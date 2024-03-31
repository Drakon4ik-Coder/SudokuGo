package main

import (
	"strconv"
	"strings"
)

func frame() {
	board.Print(2, 3)
}

func game() bool {
	initGame()
	for {
		if board.Display() {
			ClearConsole()
			frame()
		}
	}
	return true
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
