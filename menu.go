package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
	"strconv"
	"strings"
)

// fonts for highlighting
var purpleFont *color.Color
var blueFont *color.Color
var greenFont *color.Color
var redFont *color.Color
var diagonalFont *color.Color

// init keyboard listening and fonts
func init() {
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}

	blueFont = color.New(color.FgCyan)
	purpleFont = color.New(color.FgMagenta)
	greenFont = color.New(color.FgGreen)
	redFont = color.New(color.FgRed)
	diagonalFont = color.New(color.FgYellow)
}

// create menu and return true if succeeded, false if user exited
func menu() bool {
	// start menu options with output
	outputMenuStart := [4]string{"\tWelcome to Sudoku! (operate with Up and Down, then press Enter to confirm)\n", " New Game", " Load Game", " Exit"}

	// initialise menu data
	selected := 1
	outputLimit := [2]int{1, 3}

	// check if option was Changed since last print
	lastSelected := -1

	// function for drawing frame
	Draw := func() {
		ClearConsole()
		for index, element := range outputMenuStart {
			if index < outputLimit[0] {
				blueFont.Println(element)
			} else if selected == index {
				purpleFont.Println("> " + element)
			} else {
				fmt.Println(element)
			}
		}
	}

	// draw menu for the first time
	Draw()

	// iterate until option is chosen
	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyArrowUp && selected > outputLimit[0] {
			selected--
		} else if key == keyboard.KeyArrowDown && selected < outputLimit[1] {
			selected++
		} else if key == keyboard.KeyEnter {
			break
		}
		if lastSelected != selected {
			Draw()
			lastSelected = selected
		}
	}

	// choose what to do Nextx
	switch selected {
	// new game
	case 1:
		if newGameMenu() {
			blueFont.Println("Loading...")
			return initGame()
		} else {
			return false
		}
	// load old game
	case 2:
		return loadGame()
	// exit option
	case 3:
		return false
	// should be impossible, but just in case
	default:
		panic("How did you do it?!")

	}
	return true
}

var gameParam [4]int

var gameOptions = [][]string{
	{"square", "diagonal", "twodoku"},
	{"12x12", "9x9", "6x6", "4x4"},
	{"easy", "medium", "hard"},
	{"∞", "5 min", "10 min", "15 min", "30 min"},
	{},
	{},
}

func newGameMenu() bool {
	// start menu options with output
	outputMenuOptions := [7]string{"Choose game options! (operate with arrows, then press Enter to confirm either Start or Exit)\n", "Shape", "Size", "Difficulty", "Clock", "Play", "Exit"}

	/*initialise menu data*/
	selected := 5
	outputLimit := [2]int{1, 6}
	gameParam = [4]int{0, 1, 0, 0}

	// check if option was Changed since last print
	lastSelected := -1

	// function for drawing frame
	Draw := func() {
		ClearConsole()
		for index, element := range outputMenuOptions {
			if index < outputLimit[0] {
				blueFont.Println(element)
				continue
			} else if selected == index {
				purpleFont.Print("> " + element)
			} else {
				fmt.Print(element)
			}
			tmpPos := index - outputLimit[0]
			tmpLen := len(gameOptions[tmpPos])

			if tmpLen > 0 {
				if element == "Size" && gameParam[0] != 0 {
					greenFont.Print(" < 9x9 >")
				} else {
					greenFont.Print(" < " + gameOptions[tmpPos][gameParam[tmpPos]] + " >")
				}
			}
			fmt.Println()
		}
	}

	// draw menu for the first time
	Draw()

	// iterate until option is chosen
	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyArrowUp && selected > outputLimit[0] {
			selected--
		} else if key == keyboard.KeyArrowDown && selected < outputLimit[1] {
			selected++
		} else if key == keyboard.KeyArrowRight && !(gameParam[0] != 0 && selected == 2) {
			// position in gameOptions and gameParam
			tmpPos := selected - outputLimit[0]
			// number of scroll options
			tmpLen := len(gameOptions[tmpPos])

			// if there are option
			if tmpLen > 0 {
				// cycle all options
				gameParam[tmpPos] = (gameParam[tmpPos] + 1) % tmpLen
				// trigger frame redraw
				lastSelected = -1
			}

		} else if key == keyboard.KeyArrowLeft && !(gameParam[0] != 0 && selected == 2) {
			// position in gameOptions and gameParam
			tmpPos := selected - outputLimit[0]
			// number of scroll options
			tmpLen := len(gameOptions[tmpPos])
			// if there are option
			if tmpLen > 0 {
				// cycle all options
				gameParam[tmpPos] = func() int {
					newVal := gameParam[tmpPos] - 1
					if newVal < 0 {
						return tmpLen + newVal
					}
					return newVal
				}()
				// trigger frame redraw
				lastSelected = -1
			}

		} else if key == keyboard.KeyEnter {
			// exit only in cases of Play or Exit
			switch selected {
			// Play
			case 5:
				return true
			// Exit
			case 6:
				return false
			}
		}
		if lastSelected != selected {
			Draw()
			lastSelected = selected
		}
	}
	return true
}

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
