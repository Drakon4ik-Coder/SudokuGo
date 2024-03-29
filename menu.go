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

	// draw menu for the first time
	ClearConsole()
	for index, element := range outputMenuStart {
		if index < outputLimit[0] {
			infoFont.Println(element)
			continue
		} else if selected == index {
			focusFont.Print("> ")
		}
		focusFont.Println(element)
	}

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
			ClearConsole()
			for index, element := range outputMenuStart {
				if index < outputLimit[0] {
					infoFont.Println(element)
					continue
				} else if selected == index {
					focusFont.Print("> ")
				}
				focusFont.Println(element)
			}
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

var newGameParam [4]int

func newGameMenu() bool {
	// start menu options with output
	outputMenuOptions := [7]string{"\tChoose game options! (operate with arrows, then press Enter to confirm either Start or Exit)\n", " Shape", " Size", " Difficulty", " Clock", " Play", " Exit"}
	scrollOptions := [][]string{
		{"square", "diagonal", "twodoku", "triangle"},
		{"12x12", "9x9", "6x6", "4x4"},
		{"easy", "medium", "hard"},
		{"∞", "5 min", "10 min", "15 min", "30 min"},
		{},
		{},
	}

	/*initialise menu data*/
	selected := 1
	outputLimit := [2]int{1, 6}
	newGameParam = [4]int{0, 1, 0, 0}

	// check if option was changed since last print
	lastSelected := -1

	// draw menu for the first time
	ClearConsole()
	for index, element := range outputMenuOptions {
		if index < outputLimit[0] {
			infoFont.Println(element)
			continue
		} else if selected == index {
			focusFont.Print("> ")
		}
		focusFont.Print(element)
		tmpPos := index - outputLimit[0]
		tmpLen := len(scrollOptions[tmpPos])

		if tmpLen > 0 {
			optionFont.Print(" < " + scrollOptions[tmpPos][newGameParam[tmpPos]] + " >")
		}
		fmt.Println()
	}

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
			// position in scrollOptions and newGameParam
			tmpPos := selected - outputLimit[0]
			// number of scroll options
			tmpLen := len(scrollOptions[tmpPos])

			// if there are option
			if tmpLen > 0 {
				// cycle all options
				newGameParam[tmpPos] = (newGameParam[tmpPos] + 1) % tmpLen
				// trigger frame redraw
				lastSelected = -1
			}

		} else if key == keyboard.KeyArrowLeft {
			// position in scrollOptions and newGameParam
			tmpPos := selected - outputLimit[0]
			// number of scroll options
			tmpLen := len(scrollOptions[tmpPos])
			// if there are option
			if tmpLen > 0 {
				// cycle all options
				newGameParam[tmpPos] = func() int {
					newVal := newGameParam[tmpPos] - 1
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
			ClearConsole()
			for index, element := range outputMenuOptions {
				if index < outputLimit[0] {
					infoFont.Println(element)
					continue
				} else if selected == index {
					focusFont.Print("> ")
				}
				focusFont.Print(element)
				tmpPos := index - outputLimit[0]
				tmpLen := len(scrollOptions[tmpPos])

				if tmpLen > 0 {
					optionFont.Print(" < " + scrollOptions[tmpPos][newGameParam[tmpPos]] + " >")
				}
				fmt.Println()
			}
			lastSelected = selected
		}
	}
	return true
}
