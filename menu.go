package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
)

func init() {
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
}

func menuListen(stopChannel chan bool) {
	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		if key == keyboard.KeyArrowDown && selected < outputLimit[1] {
			selected++
		} else if key == keyboard.KeyArrowUp && selected > outputLimit[0] {
			selected--
		}
	}
}

// position of
var selected int
var output = [4]string{"Welcome to Sudoku!\n", " New Game", " Load Game", " Exit"}
var outputLimit = [2]int{1, 3}

func menu() bool {
	selected = 1
	lastSelected := -1
	stopChannel := make(chan bool)
	go menuListen(stopChannel)
	for {
		if lastSelected != selected {
			ClearConsole()
			for index, element := range output {
				if index == selected {
					fmt.Print(" >")
				}
				fmt.Println(element)
			}
			lastSelected = selected
		}
	}
	return true
}
