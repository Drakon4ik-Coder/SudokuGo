package main

import (
	"encoding/gob"
	"fmt"
	"github.com/tawesoft/golib/v2/dialog"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SudokuBoard interface {
	RevealRandom()      // Reveal random box
	Enter(val int) bool // Check if the val is the same as in the Board
	IsComplete() bool   // Check if the Board is complete
	Print()             // Print the Board
	Move(col, row int)  // Move the cursor if possible
	Rules() string      // Return rules of sudoku
	Display() bool      // Return whether there were any changes since last call of Print
	Undo()              // Undoes previous move
	Redo()              // Cancels last undo call
	SaveGame() error    // Save game state
	TimePass(sec int)   // Decrement left play time by sec seconds
}

type Vector2 struct {
	Xpos int
	Ypos int
}
type Change struct {
	Pos    Vector2
	OldVal int
	NewVal int
}
type DoubleChange struct {
	Main     bool
	Adjacent bool
}

type BasicSudoku struct {
	BoardShow     [][]int
	Board         [][]int
	Size          int
	NonetSize     Vector2
	Changed       bool
	CursorPos     Vector2
	Actions       []Change
	CurrentAction int
	TimeLeft      int
}

type DiagonalSudoku struct {
	BasicSudoku
}

type TwoDoku struct {
	BoardMain     BasicSudoku
	BoardAdd      BasicSudoku
	Actions       []DoubleChange
	CurrentAction int
}

func findClosestFactors(x int) (int, int) {
	var a, b int
	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}
	minDifference := math.MaxInt64

	// Iterate through possible values of a
	for i := 1; i <= int(math.Sqrt(float64(x))); i++ {
		if x%i == 0 {
			a = i
			b = x / i

			// Calculate the difference
			difference := abs(a - b)

			// Update if the difference is smaller
			if difference < minDifference {
				minDifference = difference
			}
		}
	}

	return a, x / a
}

func (s *BasicSudoku) PreInit(size int, playTime int) {
	s.Changed = true

	// Initialize Board with zeros
	createBoard := func() [][]int {
		tmp := make([][]int, size)
		for i := range tmp {
			tmp[i] = make([]int, size)
		}
		return tmp
	}
	s.Board = createBoard()
	s.BoardShow = createBoard()
	s.Size = size
	tmpSize := math.Sqrt(float64(size))
	if tmpSize == math.Trunc(tmpSize) {
		s.NonetSize.Xpos = int(tmpSize)
		s.NonetSize.Ypos = int(tmpSize)
	} else {
		w, h := findClosestFactors(size)
		s.NonetSize.Xpos = w
		s.NonetSize.Ypos = h
	}

	s.CursorPos = Vector2{0, 0}
	s.Actions = nil
	s.CurrentAction = 0
	s.TimeLeft = playTime
}

func (s *BasicSudoku) Init(size, difficulty, playTime int) {

	s.PreInit(size, playTime)

	s.FillSudoku(0)
	for i := 0; i < s.Size; i++ {
		for j := 0; j < s.Size; j++ {
			s.BoardShow[i][j] = s.Board[i][j]
		}
	}

	s.EmptyGrid(difficulty)
}
func (s *DiagonalSudoku) Init(size, difficulty, playTime int) {

	s.PreInit(size, playTime)

	s.FillSudoku(0)
	for i := 0; i < s.Size; i++ {
		for j := 0; j < s.Size; j++ {
			s.BoardShow[i][j] = s.Board[i][j]
		}
	}

	s.EmptyGrid(difficulty)
}
func (s *TwoDoku) Init(size, difficulty, playTime int) {

	s.BoardMain.PreInit(size, playTime)
	s.BoardAdd.PreInit(size, playTime)
	s.BoardAdd.CursorPos = Vector2{-1, -1}
	s.Actions = nil
	s.CurrentAction = 0

	FinishInit := func(board BasicSudoku) {
		board.FillSudoku(0)
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				board.BoardShow[i][j] = board.Board[i][j]
			}
		}
		board.EmptyGrid(difficulty)
	}
	FinishInit(s.BoardMain)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			s.BoardAdd.Board[i][j] = s.BoardMain.Board[i+6][j+6]
		}
	}
	FinishInit(s.BoardAdd)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			s.BoardAdd.BoardShow[i][j] = s.BoardMain.BoardShow[i+6][j+6]
		}
	}
}

func (s *BasicSudoku) Enter(val int) bool {
	if val > s.Size || s.Board[s.CursorPos.Xpos][s.CursorPos.Ypos] == s.BoardShow[s.CursorPos.Xpos][s.CursorPos.Ypos] || val == s.BoardShow[s.CursorPos.Xpos][s.CursorPos.Ypos] {
		return false
	}

	currAct := s.CurrentAction
	//dialog.Info(strconv.Itoa(currAct))
	for currAct != len(s.Actions) {
		s.Actions = s.Actions[0:currAct]
	}

	s.Actions = append(s.Actions, Change{s.CursorPos, s.BoardShow[s.CursorPos.Xpos][s.CursorPos.Ypos], val})
	s.CurrentAction = len(s.Actions)
	s.Changed = true

	success := true
	if s.Board[s.CursorPos.Xpos][s.CursorPos.Ypos] != val {
		success = false
	}
	s.BoardShow[s.CursorPos.Xpos][s.CursorPos.Ypos] = val
	return success
}
func (s *TwoDoku) Enter(val int) bool {
	ret := false
	if s.BoardMain.CursorPos.Xpos != -1 {
		ret = s.BoardMain.Enter(val)
	}
	if s.BoardAdd.CursorPos.Xpos != -1 {
		ret = s.BoardAdd.Enter(val)
	}

	if s.Display() {
		currAct := s.CurrentAction
		for currAct != len(s.Actions) {
			s.Actions = s.Actions[0:currAct]
		}
		if s.BoardMain.Changed && s.BoardAdd.Changed {
			s.Actions = append(s.Actions, DoubleChange{true, true})
			s.Actions = append(s.Actions, DoubleChange{false, true})
		} else if s.BoardMain.Changed {
			s.Actions = append(s.Actions, DoubleChange{true, false})
		} else if s.BoardAdd.Changed {
			s.Actions = append(s.Actions, DoubleChange{false, false})
		}
		s.CurrentAction = len(s.Actions)
	}
	return ret
}

func (s *BasicSudoku) FillSudoku(position int) bool {
	if position == s.Size*s.Size {
		return true
	}
	x := position / s.Size
	y := position % s.Size
	if s.Board[x][y] != 0 {
		return s.FillSudoku(position + 1)
	}

	available := s.AvailableNum(position, s.Board)
	lenAvailable := len(available)
	if lenAvailable == 0 {
		return false
	}

	for i := 0; i < lenAvailable; i++ {
		s.Board[x][y] = available[i]
		if s.FillSudoku(position + 1) {
			return true
		}
	}
	s.Board[x][y] = 0
	return false
}
func (s *DiagonalSudoku) FillSudoku(position int) bool {
	if position == s.Size*s.Size {
		return true
	}
	x := position / s.Size
	y := position % s.Size
	if s.Board[x][y] != 0 {
		return s.FillSudoku(position + 1)
	}

	available := s.AvailableNum(position, s.Board)
	lenAvailable := len(available)
	if lenAvailable == 0 {
		return false
	}

	for i := 0; i < lenAvailable; i++ {
		s.Board[x][y] = available[i]
		if s.FillSudoku(position + 1) {
			return true
		}
	}
	s.Board[x][y] = 0
	return false
}

func (s *BasicSudoku) AvailableNum(position int, board [][]int) []int {
	x := position / s.Size
	y := position % s.Size

	contains := func(slice []int, value int) bool {
		for _, item := range slice {
			if item == value {
				return true
			}
		}
		return false
	}

	taken := []int{}

	nonetStart := [2]int{x - x%s.NonetSize.Xpos, y - y%s.NonetSize.Ypos}
	for i := 0; i < s.NonetSize.Xpos; i++ {
		for j := 0; j < s.NonetSize.Ypos; j++ {
			tmp := board[nonetStart[0]+i][nonetStart[1]+j]
			if tmp != 0 && !contains(taken, tmp) {
				taken = append(taken, tmp)
			}
		}
	}

	for i := 0; i < s.Size; i++ {
		tmpHorizontal := board[i][y]
		if tmpHorizontal != 0 && !contains(taken, tmpHorizontal) {
			taken = append(taken, tmpHorizontal)
		}
		tmpVertical := board[x][i]
		if tmpVertical != 0 && !contains(taken, tmpVertical) {
			taken = append(taken, tmpVertical)
		}
	}

	var available []int

	for i := 1; i <= s.Size; i++ {
		if !contains(taken, i) {
			available = append(available, i)
		}
	}

	for i := range available {
		j := rand.Intn(i + 1)
		available[i], available[j] = available[j], available[i]
	}

	return available
}
func (s *DiagonalSudoku) AvailableNum(position int, board [][]int) []int {
	x := position / s.Size
	y := position % s.Size

	contains := func(slice []int, value int) bool {
		for _, item := range slice {
			if item == value {
				return true
			}
		}
		return false
	}

	taken := []int{}

	nonetStart := [2]int{x - x%s.NonetSize.Xpos, y - y%s.NonetSize.Ypos}
	for i := 0; i < s.NonetSize.Xpos; i++ {
		for j := 0; j < s.NonetSize.Ypos; j++ {
			tmp := board[nonetStart[0]+i][nonetStart[1]+j]
			if tmp != 0 && !contains(taken, tmp) {
				taken = append(taken, tmp)
			}
		}
	}

	for i := 0; i < s.Size; i++ {
		tmpHorizontal := board[i][y]
		if tmpHorizontal != 0 && !contains(taken, tmpHorizontal) {
			taken = append(taken, tmpHorizontal)
		}
		tmpVertical := board[x][i]
		if tmpVertical != 0 && !contains(taken, tmpVertical) {
			taken = append(taken, tmpVertical)
		}
	}

	if x == y {
		for i := 0; i < s.Size; i++ {
			tmpLeftDiagonal := board[i][i]
			if tmpLeftDiagonal != 0 && !contains(taken, tmpLeftDiagonal) {
				taken = append(taken, tmpLeftDiagonal)
			}
		}
	}

	if x == s.Size-y-1 {
		for i := 0; i < s.Size; i++ {
			tmpRightDiagonal := board[s.Size-i-1][i]
			if tmpRightDiagonal != 0 && !contains(taken, tmpRightDiagonal) {
				taken = append(taken, tmpRightDiagonal)
			}
		}
	}

	var available []int

	for i := 1; i <= s.Size; i++ {
		if !contains(taken, i) {
			available = append(available, i)
		}
	}

	for i := range available {
		j := rand.Intn(i + 1)
		available[i], available[j] = available[j], available[i]
	}

	return available
}

var solutions int

func (s *BasicSudoku) EmptyGrid(difficulty int) {
	emptyFinal := (difficulty + 1) * 25
	copyGrid := BasicSudoku{}
	finished := false
	empty := 0
	startTime := time.Now()

	for !finished {
		row, col := rand.Intn(s.Size), rand.Intn(s.Size)
		for s.BoardShow[row][col] == 0 {
			row, col = rand.Intn(s.Size), rand.Intn(s.Size)
		}

		copyGrid.Copy(s)
		solutions = 0
		copyGrid.BoardShow[row][col] = 0

		copyGrid.SolveGrid(0)

		if solutions == 1 {
			s.BoardShow[row][col] = 0
			empty++
		}

		if 100*empty/(s.Size*s.Size) > emptyFinal || time.Since(startTime).Seconds() > 2 {
			finished = true
		}

	}

}
func (s *DiagonalSudoku) EmptyGrid(difficulty int) {
	emptyFinal := (difficulty + 1) * 25
	copyGrid := DiagonalSudoku{}
	finished := false
	empty := 0
	startTime := time.Now()

	for !finished {
		row, col := rand.Intn(s.Size), rand.Intn(s.Size)
		for s.BoardShow[row][col] == 0 {
			row, col = rand.Intn(s.Size), rand.Intn(s.Size)
		}

		copyGrid.Copy(s)
		solutions = 0
		copyGrid.BoardShow[row][col] = 0

		copyGrid.SolveGrid(0)

		if solutions == 1 {
			s.BoardShow[row][col] = 0
			empty++
		}

		if 100*empty/(s.Size*s.Size) > emptyFinal || time.Since(startTime).Seconds() > 2 {
			finished = true
		}

	}

}

func (s *BasicSudoku) SolveGrid(position int) {
	if position == s.Size*s.Size {
		solutions++
		return
	}

	x := position / s.Size
	y := position % s.Size

	if s.BoardShow[x][y] != 0 {
		s.SolveGrid(position + 1)
		return
	}

	available := s.AvailableNum(position, s.BoardShow)
	lenAvailable := len(available)

	for i := 0; i < lenAvailable; i++ {
		s.BoardShow[x][y] = available[i]
		s.SolveGrid(position + 1)
	}

	s.BoardShow[x][y] = 0
	return
}
func (s *DiagonalSudoku) SolveGrid(position int) {
	if position == s.Size*s.Size {
		solutions++
		return
	}

	x := position / s.Size
	y := position % s.Size

	if s.BoardShow[x][y] != 0 {
		s.SolveGrid(position + 1)
		return
	}

	available := s.AvailableNum(position, s.BoardShow)
	lenAvailable := len(available)

	for i := 0; i < lenAvailable; i++ {
		s.BoardShow[x][y] = available[i]
		s.SolveGrid(position + 1)
	}

	s.BoardShow[x][y] = 0
	return
}

func SudokuPrintManual() {
	blueFont.Println("Move with arrows, enter with numbers 1-9(and A-C, depending on Board Size)")
	greenFont.Print("Green")
	blueFont.Println(" - solved")
	redFont.Print("Red")
	blueFont.Println(" - incorrect")
	purpleFont.Print("Purple")
	blueFont.Println(" - cursor")
}

func (s *BasicSudoku) Print() {
	if !s.Changed {
		return
	}
	SudokuPrintManual()
	s.Changed = false
	printFont := fmt.Printf
	for i, line := range s.BoardShow {
		if i%s.NonetSize.Xpos == 0 {
			if i > 0 {
				blueFont.Print(strings.Repeat("|"+strings.Repeat("_", s.NonetSize.Ypos*2+1), s.Size/s.NonetSize.Ypos))
				blueFont.Println("|")
			} else {
				blueFont.Print(strings.Repeat("_"+strings.Repeat("_", s.NonetSize.Ypos*2+1), s.Size/s.NonetSize.Ypos))
				blueFont.Println("_")
			}
		}
		for j, element := range line {
			if j%s.NonetSize.Ypos == 0 {
				blueFont.Print("| ")
			}
			if element != 0 && s.Board[i][j] != element {
				printFont = redFont.Printf
			} else if s.CursorPos.Xpos == i && s.CursorPos.Ypos == j {
				printFont = purpleFont.Printf
			} else if element != 0 {
				printFont = greenFont.Printf
			} else {
				printFont = fmt.Printf
			}
			if element > 9 {
				_, _ = printFont("%c ", 'A'+(element-10))
			} else {
				_, _ = printFont("%d ", element)
			}
		}
		blueFont.Print("|")
		fmt.Println()
	}
	blueFont.Print(strings.Repeat("|"+strings.Repeat("_", s.NonetSize.Ypos*2+1), s.Size/s.NonetSize.Ypos))
	blueFont.Println("|")
	//redFont.Println(s.CurrentAction)
	//for e := s.Actions.Front(); e != nil; e = e.Next() {
	//	redFont.Println(e.Value)
	//}
	//redFont.Println()
	//if s.CurrentAction != nil {
	//	redFont.Println(s.CurrentAction.Value)
	//}
}
func (s *DiagonalSudoku) Print() {
	if !s.Changed {
		return
	}
	SudokuPrintManual()
	diagonalFont.Print("Yellow")
	blueFont.Println(" - correct diagonal")

	s.Changed = false
	printFont := fmt.Printf
	for i, line := range s.BoardShow {
		if i%s.NonetSize.Xpos == 0 {
			if i > 0 {
				blueFont.Print(strings.Repeat("|"+strings.Repeat("_", s.NonetSize.Ypos*2+1), s.Size/s.NonetSize.Ypos))
				blueFont.Println("|")
			} else {
				blueFont.Print(strings.Repeat("_"+strings.Repeat("_", s.NonetSize.Ypos*2+1), s.Size/s.NonetSize.Ypos))
				blueFont.Println("_")
			}
		}
		for j, element := range line {
			if j%s.NonetSize.Ypos == 0 {
				blueFont.Print("| ")
			}
			if element != 0 && s.Board[i][j] != element {
				printFont = redFont.Printf
			} else if s.CursorPos.Xpos == i && s.CursorPos.Ypos == j {
				printFont = purpleFont.Printf
			} else if (i == j || i == s.Size-j-1) && element != 0 {
				printFont = diagonalFont.Printf
			} else if element != 0 {
				printFont = greenFont.Printf
			} else {
				printFont = fmt.Printf
			}
			if element > 9 {
				_, _ = printFont("%c ", 'A'+(element-10))
			} else {
				_, _ = printFont("%d ", element)
			}
		}
		blueFont.Print("|")
		fmt.Println()
	}
	blueFont.Print(strings.Repeat("|"+strings.Repeat("_", s.NonetSize.Ypos*2+1), s.Size/s.NonetSize.Ypos))
	blueFont.Println("|")
}
func (s *TwoDoku) Print() {
	if !s.BoardMain.Changed && !s.BoardAdd.Changed {
		return
	}
	SudokuPrintManual()
	s.BoardMain.Changed = false
	s.BoardAdd.Changed = false
	printFont := fmt.Printf
	for i := 0; i < 6; i++ {
		if i%3 == 0 {
			if i > 0 {
				blueFont.Print(strings.Repeat("|_______", 3))
				blueFont.Println("|")
			} else {
				blueFont.Print(strings.Repeat("________", 3))
				blueFont.Println("_")
			}
		}
		for j := 0; j < 9; j++ {
			if j%3 == 0 {
				blueFont.Print("| ")
			}
			element := s.BoardMain.BoardShow[i][j]
			if element != 0 && s.BoardMain.Board[i][j] != element {
				printFont = redFont.Printf
			} else if s.BoardMain.CursorPos.Xpos == i && s.BoardMain.CursorPos.Ypos == j {
				printFont = purpleFont.Printf
			} else if element != 0 {
				printFont = greenFont.Printf
			} else {
				printFont = fmt.Printf
			}
			_, _ = printFont("%d ", element)
		}
		blueFont.Print("|")
		fmt.Println()
	}
	blueFont.Print(strings.Repeat("|_______", 3))
	blueFont.Print("|")
	blueFont.Println(strings.Repeat("________", 2))
	for i := 6; i < 9; i++ {
		for j := 0; j < 15; j++ {
			var value int
			var element int
			if j < 9 {
				element = s.BoardMain.BoardShow[i][j]
				value = s.BoardMain.Board[i][j]
			} else {
				element = s.BoardAdd.BoardShow[i-6][j%9+3]
				value = s.BoardAdd.Board[i-6][j%9+3]
			}

			if j%3 == 0 {
				blueFont.Print("| ")
			}

			if element != 0 && value != element {
				printFont = redFont.Printf
			} else if s.BoardMain.CursorPos.Xpos == i && s.BoardMain.CursorPos.Ypos == j || s.BoardAdd.CursorPos.Xpos == i-6 && s.BoardAdd.CursorPos.Ypos == j-6 {
				printFont = purpleFont.Printf
			} else if element != 0 {
				printFont = greenFont.Printf
			} else {
				printFont = fmt.Printf
			}
			_, _ = printFont("%d ", element)
		}
		blueFont.Print("|")
		fmt.Println()
	}
	blueFont.Print(strings.Repeat("|_______", 5))
	blueFont.Println("|")
	for i := 3; i < 9; i++ {
		if i%3 == 0 && i != 3 {
			fmt.Print(strings.Repeat("  ", 8))
			blueFont.Print(strings.Repeat("|_______", 3))
			blueFont.Println("|")

		}
		fmt.Print(strings.Repeat("  ", 8))
		for j := 0; j < 9; j++ {
			if j%3 == 0 {
				blueFont.Print("| ")
			}
			element := s.BoardAdd.BoardShow[i][j]
			if element != 0 && s.BoardAdd.Board[i][j] != element {
				printFont = redFont.Printf
			} else if s.BoardAdd.CursorPos.Xpos == i && s.BoardAdd.CursorPos.Ypos == j {
				printFont = purpleFont.Printf
			} else if element != 0 {
				printFont = greenFont.Printf
			} else {
				printFont = fmt.Printf
			}
			_, _ = printFont("%d ", element)
		}
		blueFont.Print("|")
		fmt.Println()
	}
	fmt.Print(strings.Repeat("  ", 8))
	blueFont.Print(strings.Repeat("|_______", 3))
	blueFont.Println("|")
	//for e := s.Actions.Front(); e != nil; e = e.Next() {
	//	redFont.Println(e.Value)
	//}
	//redFont.Println()
	//if s.CurrentAction != nil {
	//	redFont.Println(s.CurrentAction.Value)
	//}
	//redFont.Println(s.BoardMain.CursorPos)
	//redFont.Println(s.BoardAdd.CursorPos)
}

func (s *BasicSudoku) Copy(s2 *BasicSudoku) {
	s.Size = s2.Size
	s.NonetSize = s2.NonetSize

	createBoard := func() [][]int {
		tmp := make([][]int, s.Size)
		for i := range tmp {
			tmp[i] = make([]int, s.Size)
		}
		return tmp
	}
	s.Board = createBoard()
	s.BoardShow = createBoard()

	for i := 0; i < s.Size; i++ {
		for j := 0; j < s.Size; j++ {
			s.Board[i][j] = s2.Board[i][j]
			s.BoardShow[i][j] = s2.BoardShow[i][j]
		}
	}

}
func (s *DiagonalSudoku) Copy(s2 *DiagonalSudoku) {
	s.Size = s2.Size
	s.NonetSize = s2.NonetSize

	createBoard := func() [][]int {
		tmp := make([][]int, s.Size)
		for i := range tmp {
			tmp[i] = make([]int, s.Size)
		}
		return tmp
	}
	s.Board = createBoard()
	s.BoardShow = createBoard()

	for i := 0; i < s.Size; i++ {
		for j := 0; j < s.Size; j++ {
			s.Board[i][j] = s2.Board[i][j]
			s.BoardShow[i][j] = s2.BoardShow[i][j]
		}
	}

}

func (s *BasicSudoku) IsComplete() bool {
	for i, line := range s.BoardShow {
		for j, element := range line {
			if element == 0 || s.Board[i][j] != element {
				return false
			}
		}
	}
	s.Changed = true
	s.CursorPos = Vector2{-1, -1}
	return true
}
func (s *TwoDoku) IsComplete() bool {
	return s.BoardMain.IsComplete() && s.BoardAdd.IsComplete()
}

func (s *BasicSudoku) Display() bool {
	return s.Changed
}
func (s *TwoDoku) Display() bool {
	return s.BoardMain.Changed || s.BoardAdd.Changed
}

func (s *BasicSudoku) Rules() string {
	return "Sudoku is played on a 9x9(or other sizes) grid where each row, column,\n" +
		"and 3x3(can differ) region must contain all digits from 1 to 9(possibly up to C) without repetition.\n" +
		"Use logic to fill in the empty cells based on the filled cells.\n" +
		"No guessing is allowed, and each puzzle has exactly one unique solution."
}
func (s *DiagonalSudoku) Rules() string {
	return "Diagonal Sudoku is played on a 9x9 grid where each row, column, diagonal\n" +
		"and 3x3 region must contain all digits from 1 to 9 without repetition.\n" +
		"Use logic to fill in the empty cells based on the filled cells.\n" +
		"No guessing is allowed, and each puzzle has exactly one unique solution."
}
func (s *TwoDoku) Rules() string {
	return "Twodoku consists of two 9x9 Sudoku puzzles\n" +
		"that share the same 3x3 region(9-th and 1-st for the corresponding Board).\n" +
		"Each row, column, and 3x3(can differ) region must contain all digits\n" +
		"from 1 to 9(possibly up to C) without repetition.\n" +
		"Use logic to fill in the empty cells based on the filled cells.\n" +
		"No guessing is allowed, and each puzzle has exactly one unique solution."
}

func (s *BasicSudoku) Move(col, row int) {
	newPos := Vector2{s.CursorPos.Xpos + row, s.CursorPos.Ypos + col}
	if newPos.Xpos < s.Size && newPos.Ypos < s.Size && newPos.Xpos >= 0 && newPos.Ypos >= 0 {
		s.CursorPos.Xpos = newPos.Xpos
		s.CursorPos.Ypos = newPos.Ypos
		s.Changed = true
	}
}
func (s *TwoDoku) Move(col, row int) {
	var currentPos Vector2
	if s.BoardMain.CursorPos.Xpos == -1 {
		currentPos = s.BoardAdd.CursorPos
		currentPos.Xpos += 6
		currentPos.Ypos += 6
	} else if s.BoardAdd.CursorPos.Xpos == -1 {
		currentPos = s.BoardMain.CursorPos
	} else {
		currentPos = s.BoardMain.CursorPos
	}
	newPos := Vector2{currentPos.Xpos + row, currentPos.Ypos + col}
	// invalid position
	if newPos.Xpos > 14 || newPos.Ypos > 14 || newPos.Xpos < 0 || newPos.Ypos < 0 || (newPos.Xpos > 8 && newPos.Ypos < 6) || (newPos.Xpos < 6 && newPos.Ypos > 8) {
		return
	}
	s.BoardAdd.CursorPos.Xpos = newPos.Xpos - 6
	s.BoardAdd.CursorPos.Ypos = newPos.Ypos - 6
	s.BoardMain.CursorPos.Xpos = newPos.Xpos
	s.BoardMain.CursorPos.Ypos = newPos.Ypos
	if (newPos.Xpos <= 8 && newPos.Ypos < 6) || (newPos.Xpos < 6 && newPos.Ypos <= 8) {
		s.BoardAdd.CursorPos = Vector2{-1, -1}

	} else if (newPos.Xpos > 8 && newPos.Ypos >= 6) || (newPos.Xpos >= 6 && newPos.Ypos > 8) {
		s.BoardMain.CursorPos = Vector2{-1, -1}
	}
	s.BoardMain.Changed = true
}

func (s *BasicSudoku) Undo() {
	action := s.CurrentAction - 1
	if action == -1 {
		return
	}
	s.CurrentAction = action
	prevChange := s.Actions[action]
	change := Vector2{s.CursorPos.Xpos - prevChange.Pos.Xpos, s.CursorPos.Ypos - prevChange.Pos.Ypos}
	//dialog.Info(fmt.Sprintf("%d %d", change.Xpos, change.Ypos))
	s.Move(-change.Ypos, -change.Xpos)
	s.BoardShow[s.CursorPos.Xpos][s.CursorPos.Ypos] = prevChange.OldVal
}
func (s *TwoDoku) Undo() {
	action := s.CurrentAction - 1
	if action == -1 {
		return
	}
	s.CurrentAction = action
	prevChange := s.Actions[action]
	if prevChange.Main {
		s.BoardMain.Undo()
		s.BoardAdd.CursorPos = Vector2{-1, -1}
	} else {
		s.BoardAdd.Undo()
		s.BoardMain.CursorPos = Vector2{-1, -1}
		if prevChange.Adjacent {
			s.Undo()
		}
	}
	s.Move(0, 0)
}

func (s *BasicSudoku) Redo() {
	action := s.CurrentAction
	if action == len(s.Actions) {
		return
	}
	prevChange := s.Actions[action]
	change := Vector2{s.CursorPos.Xpos - prevChange.Pos.Xpos, s.CursorPos.Ypos - prevChange.Pos.Ypos}
	//dialog.Info(fmt.Sprintf("%d %d", change.Xpos, change.Ypos))
	s.Move(-change.Ypos, -change.Xpos)
	s.BoardShow[s.CursorPos.Xpos][s.CursorPos.Ypos] = prevChange.NewVal

	s.CurrentAction++
}
func (s *TwoDoku) Redo() {
	action := s.CurrentAction
	if action == len(s.Actions) {
		return
	}
	prevChange := s.Actions[action]
	s.CurrentAction++
	if prevChange.Main {
		s.BoardMain.Redo()
		s.BoardAdd.CursorPos = Vector2{-1, -1}
		if prevChange.Adjacent {
			s.Redo()
		}
	} else {
		s.BoardAdd.Redo()
		s.BoardMain.CursorPos = Vector2{-1, -1}
	}
	s.Move(0, 0)
}

func ClearFiles() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	files, err := filepath.Glob(filepath.Join(currentDir, "*.sudo"))
	if err != nil {
		return err
	}

	for _, file := range files {
		err = os.Remove(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *BasicSudoku) SaveGame() error {
	err := ClearFiles()
	if err != nil {
		return err
	}
	file, err := os.Create("basic.sudo")
	if err != nil {
		return err
	}
	defer file.Close()

	gob.Register(BasicSudoku{})

	// Create an encoder and send some values.
	enc := gob.NewEncoder(file)

	err = enc.Encode(s)
	if err != nil {
		dialog.Info(err.Error())
		return err
	}

	return nil
}
func (s *DiagonalSudoku) SaveGame() error {
	err := ClearFiles()
	if err != nil {
		return err
	}
	file, err := os.Create("diagonal.sudo")
	if err != nil {
		return err
	}
	defer file.Close()

	gob.Register(DiagonalSudoku{})

	// Create an encoder and send some values.
	enc := gob.NewEncoder(file)

	err = enc.Encode(s)
	if err != nil {
		dialog.Info(err.Error())
		return err
	}

	return nil
}
func (s *TwoDoku) SaveGame() error {
	err := ClearFiles()
	if err != nil {
		return err
	}
	file, err := os.Create("two.sudo")
	if err != nil {
		return err
	}
	defer file.Close()

	gob.Register(TwoDoku{})

	// Create an encoder and send some values.
	enc := gob.NewEncoder(file)

	err = enc.Encode(s)
	if err != nil {
		dialog.Info(err.Error())
		return err
	}

	return nil
}

func (s *BasicSudoku) RevealRandom() {
	row := rand.Intn(s.Size)
	col := rand.Intn(s.Size)
	look := func(x, y int) bool {
		for i := x; i < s.Size; i++ {
			for j := y; j < s.Size; j++ {
				if s.BoardShow[i][j] != s.Board[i][j] {
					s.CursorPos = Vector2{i, j}
					s.Enter(s.Board[i][j])
					return true
				}
			}
		}
		return false
	}
	if look(row, col) {
		return
	}
	look(0, 0)
}
func (s *TwoDoku) RevealRandom() {
	row := rand.Intn(9)
	col := rand.Intn(9)
	look := func(x, y int) bool {
		for i := x; i < 9; i++ {
			for j := y; j < 9; j++ {
				if s.BoardMain.BoardShow[i][j] != s.BoardMain.Board[i][j] {
					s.BoardMain.CursorPos = Vector2{i, j}
					if i >= 6 && j >= 6 {
						s.BoardAdd.CursorPos = Vector2{i - 6, j - 6}
					} else {
						s.BoardAdd.CursorPos = Vector2{-1, -1}
					}
					s.Enter(s.BoardMain.Board[i][j])
					return true
				}
				if s.BoardAdd.BoardShow[i][j] != s.BoardAdd.Board[i][j] {
					s.BoardAdd.CursorPos = Vector2{i, j}
					if i <= 2 && j <= 2 {
						s.BoardMain.CursorPos = Vector2{i + 6, j + 6}
					} else {
						s.BoardMain.CursorPos = Vector2{-1, -1}
					}
					s.Enter(s.BoardAdd.Board[i][j])
					return true
				}
			}
		}
		return false
	}
	if look(row, col) {
		return
	}
	look(0, 0)
}

func (s *BasicSudoku) TimePass(sec int) {
	if s.TimeLeft > 0 {
		s.TimeLeft -= sec
		if s.TimeLeft <= 0 {
			s.TimeLeft = 0
		}
	}
}
func (s *TwoDoku) TimePass(sec int) {
	s.BoardMain.TimePass(sec)
}
