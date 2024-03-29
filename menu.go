package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
)

// fonts for highlighting
var focusFont *color.Color
var infoFont *color.Color
var optionFont *color.Color

// init keyboard listening and fonts
func init() {
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}

	infoFont = color.New(color.FgCyan)
	focusFont = color.New(color.FgMagenta)
	optionFont = color.New(color.FgGreen)
}

// create menu and return true if succeeded, false if user exited
func menu() bool {
	// start menu options with output
	outputMenuStart := [4]string{"\tWelcome to Sudoku! (operate with Up and Down, then press Enter to confirm)\n", " New Game", " Load Game", " Exit"}

	// initialise menu data
	selected := 1
	outputLimit := [2]int{1, 3}

	// check if option was changed since last print
	lastSelected := -1

	// function for drawing frame
	var Draw func() = func() {
		ClearConsole()
		for index, element := range outputMenuStart {
			if index < outputLimit[0] {
				infoFont.Println(element)
			} else if selected == index {
				focusFont.Println("> " + element)
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

	// choose what to do next
	switch selected {
	// new game
	case 1:
		return newGameMenu()
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
	{"square", "diagonal", "twodoku", "triangle"},
	{"12x12", "9x9", "6x6", "4x4"},
	{"easy", "medium", "hard"},
	{"∞", "5 min", "10 min", "15 min", "30 min"},
	{},
	{},
}

func newGameMenu() bool {
	// start menu options with output
	outputMenuOptions := [7]string{"\tChoose game options! (operate with arrows, then press Enter to confirm either Start or Exit)\n", " Shape", " Size", " Difficulty", " Clock", " Play", " Exit"}

	/*initialise menu data*/
	selected := 1
	outputLimit := [2]int{1, 6}
	gameParam = [4]int{0, 1, 0, 0}

	// check if option was changed since last print
	lastSelected := -1

	// function for drawing frame
	var Draw func() = func() {
		ClearConsole()
		for index, element := range outputMenuOptions {
			if index < outputLimit[0] {
				infoFont.Println(element)
				continue
			} else if selected == index {
				focusFont.Print("> " + element)
			} else {
				fmt.Print(element)
			}
			tmpPos := index - outputLimit[0]
			tmpLen := len(gameOptions[tmpPos])

			if tmpLen > 0 {
				optionFont.Print(" < " + gameOptions[tmpPos][gameParam[tmpPos]] + " >")
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
		} else if key == keyboard.KeyArrowRight {
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

		} else if key == keyboard.KeyArrowLeft {
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
