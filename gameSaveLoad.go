package main

import (
	"encoding/gob"
	"os"
	"path/filepath"
	"strings"
)

// load game from file
func loadGame() bool {
	// get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return false
	}

	// find all .sudo files in current directory
	files, err := filepath.Glob(filepath.Join(currentDir, "*.sudo"))
	if err != nil {
		return false
	}

	// if there is not 1 .sudo file - something is wrong
	if len(files) != 1 {
		return false
	}
	// open file for read
	file, err := os.Open(files[0])
	if err != nil {
		return false
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)

	if strings.HasSuffix(files[0], "basic.sudo") {
		var bas BasicSudoku
		if err = decoder.Decode(&bas); err != nil {
			return false
		}
		board = &bas
	} else if strings.HasSuffix(files[0], "diagonal.sudo") {
		var bas DiagonalSudoku
		if err = decoder.Decode(&bas); err != nil {
			return false
		}
		board = &bas
	} else if strings.HasSuffix(files[0], "two.sudo") {
		var bas TwoDoku
		if err = decoder.Decode(&bas); err != nil {
			return false
		}
		board = &bas
	}
	// update so the Display is true
	board.Move(0, 0)
	return true
}
