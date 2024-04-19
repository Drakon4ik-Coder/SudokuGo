package main

import (
	"encoding/gob"
	"github.com/tawesoft/golib/v2/dialog"
	"os"
	"path/filepath"
	"strings"
)

func loadGame() bool {
	currentDir, err := os.Getwd()
	if err != nil {
		return false
	}

	files, err := filepath.Glob(filepath.Join(currentDir, "*.sudo"))
	if err != nil {
		return false
	}

	if len(files) != 1 {
		return false
	}
	file, err := os.Open(files[0])
	if err != nil {
		return false
	}
	defer file.Close()

	gob.Register(BasicSudoku{})

	decoder := gob.NewDecoder(file)

	if strings.HasSuffix(files[0], "basic.sudo") {
		var bas BasicSudoku
		if err = decoder.Decode(&bas); err != nil {
			dialog.Info("Fuck")
			return false
		}
		board = &bas
	} else if strings.HasSuffix(files[0], "diagonal.sudo") {
		var bas DiagonalSudoku
		if err = decoder.Decode(&bas); err != nil {
			dialog.Info("Fuck")
			return false
		}
		board = &bas
	} else if strings.HasSuffix(files[0], "two.sudo") {
		var bas TwoDoku
		if err = decoder.Decode(&bas); err != nil {
			dialog.Info("Fuck")
			return false
		}
		board = &bas
	}
	board.Move(0, 0)
	return true
}
