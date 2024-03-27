package main

import (
	"fmt"
	"github.com/inancgumus/screen"
	"strings"
)

func init() {
	screen.Clear()
	screen.MoveTopLeft()
}

func ClearConsole() {
	screen.MoveTopLeft()
	windth, height := screen.Size()
	fmt.Print(strings.Repeat(strings.Repeat(" ", windth)+"\n", height))
	screen.MoveTopLeft()
}
