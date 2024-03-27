package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
)

// print with highlighted font
var optionFont *color.Color
var infoFont *color.Color

func init() {
	infoFont = color.New(color.FgCyan)
	optionFont = color.New(color.FgMagenta)
}

// init keyboard listening
func init() {
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
}

// detect arrows Up and Down and Enter and change corresponding global variables
func menuListen() {
	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		if key == keyboard.KeyArrowDown && selected < outputLimit[1] {
			selected++
		} else if key == keyboard.KeyArrowUp && selected > outputLimit[0] {
			selected--
		} else if key == keyboard.KeyEnter {
			chosen = true
			break
		}
	}
}

// select option
var selected int

// where options start and end
var outputLimit [2]int

// was the choice made
var chosen bool

// create menu and return true if succeeded, false if user exited
func menu() bool {
	// start menu options with output
	outputMenuStart := [4]string{"\tWelcome to Sudoku! (operate with Up and Down, then press Enter to confirm)\n", " New Game", " Load Game", " Exit"}

	// initialise menu data
	selected = 1
	outputLimit = [2]int{1, 3}
	chosen = false

	// check if option was changed since last print
	lastSelected := -1

	// run background key listener
	go menuListen()

	// iterate until option is chosen
	for {
		if lastSelected != selected {
			ClearConsole()
			for index, element := range outputMenuStart {
				if index == selected {
					optionFont.Println(" >" + element)
				} else if index < outputLimit[0] {
					infoFont.Println(element)
				} else {
					fmt.Println(element)
				}
			}
			lastSelected = selected
		}
		if chosen {
			break
		}
	}
	switch selected {
	case 1:
		return newGameMenu()
	case 2:
		return loadGame()
	case 3:
		return false
	default:
		panic("How did you do it?!")

	}
	return true
}

func newGameMenu() bool {
	return true
}
