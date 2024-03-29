package main

import (
	"strconv"
	"strings"
)

func frame() bool {
	return true
}

func game() bool {
	initGame()
	frameContinue := frame()
	for frameContinue {
		frameContinue = frame()
		ClearConsole()
	}
	return true
}

var boardType string
var boardSize int
var boardNormal [][]int
var boardTwodoku [][][]int
var boardTriangle [][]int

func initGame() bool {
	initBoard()
	return true
}

func initBoard() {
	boardType = gameOptions[0][gameParam[0]]
	boardSize, _ = strconv.Atoi(strings.Split(gameOptions[1][gameParam[1]], "x")[0])
	if boardType == "square" || boardType == "diagonal" {
		boardNormal = make([][]int, boardSize)
		for i := range boardNormal {
			boardNormal[i] = make([]int, boardSize)
		}
	} else if boardType == "twodoku" {
		boardTwodoku = make([][][]int, 2)
		for i := 0; i < 2; i++ {
			boardTwodoku[i] = make([][]int, boardSize)
			for j := range boardTwodoku[i] {
				boardTwodoku[i][j] = make([]int, boardSize)
			}
		}
	} else if boardType == "triangle" {
		boardTriangle = make([][]int, boardSize)
		for i := range boardTriangle {
			boardTriangle[i] = make([]int, 2*i+1)
		}
	}
}
