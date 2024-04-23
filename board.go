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

// SudokuBoard is sudoku interface
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
	TimeEnd() bool      // Return whether the game time has ended
}

// Vector2 to store two dimensional vector values
type Vector2 struct {
	Xpos int
	Ypos int
}

// Change to store move made in Basic and Diagonal sudoku
type Change struct {
	Pos    Vector2 // position of change
	OldVal int     // value before move
	NewVal int     // value after move
}

// DoubleChange to store move made in TwoDoku
type DoubleChange struct {
	Main     bool // move was made in main board
	Adjacent bool // move was made in adjacent(9-th and 1-st) nonet
}

// BasicSudoku struct to implement basic sudoku board
type BasicSudoku struct {
	BoardShow     [][]int  // board to show to user
	Board         [][]int  // board with answers
	Size          int      // size of the board
	NonetSize     Vector2  // size of one nonet(width and height)
	Changed       bool     // changed since last Print function call
	CursorPos     Vector2  // current player position
	Actions       []Change // store player moves
	CurrentAction int      // current move
	TimeLeft      int      // time left on timer
}

// DiagonalSudoku struct to implement diagonal sudoku board
type DiagonalSudoku struct {
	BasicSudoku // same attributes as in basic sudoku
}

// TwoDoku struct to implement twodoku board
type TwoDoku struct {
	BoardMain     BasicSudoku    // first board
	BoardAdd      BasicSudoku    // second board
	Actions       []DoubleChange // store player moves
	CurrentAction int            // current move
}

// calculate a and b so a*b is x and a-b is minimum
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

// PreInit include same init steps for all boards
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

	// calculate nonet size
	tmpSize := math.Sqrt(float64(size))
	if tmpSize == math.Trunc(tmpSize) { // if square root is integer - nonet is square
		s.NonetSize.Xpos = int(tmpSize)
		s.NonetSize.Ypos = int(tmpSize)
	} else { // else - calculate closest factors
		w, h := findClosestFactors(size)
		s.NonetSize.Xpos = w
		s.NonetSize.Ypos = h
	}

	s.CursorPos = Vector2{0, 0}
	s.Actions = nil
	s.CurrentAction = 0
	s.TimeLeft = playTime
}

// Init to init the board
func (s *BasicSudoku) Init(size, difficulty, playTime int) {
	// preinit call
	s.PreInit(size, playTime)

	// fill sudoku
	s.FillSudoku(0)
	// copy main board to show board
	for i := 0; i < s.Size; i++ {
		for j := 0; j < s.Size; j++ {
			s.BoardShow[i][j] = s.Board[i][j]
		}
	}

	// empty grid
	s.EmptyGrid(difficulty)
}
func (s *DiagonalSudoku) Init(size, difficulty, playTime int) {
	// preinit call
	s.PreInit(size, playTime)

	// fill sudoku
	s.FillSudoku(0)
	// copy main board to show board
	for i := 0; i < s.Size; i++ {
		for j := 0; j < s.Size; j++ {
			s.BoardShow[i][j] = s.Board[i][j]
		}
	}

	// empty grid
	s.EmptyGrid(difficulty)
}
func (s *TwoDoku) Init(size, difficulty, playTime int) {
	// preinit call
	s.BoardMain.PreInit(size, playTime)
	s.BoardAdd.PreInit(size, playTime)
	s.BoardAdd.CursorPos = Vector2{-1, -1} // cursor is -1 -1 if not in scope of the board
	s.Actions = nil
	s.CurrentAction = 0

	// last steps to init board
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
	// copy main board's 9-th nonet to additional board's 1-st nonet
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			s.BoardAdd.Board[i][j] = s.BoardMain.Board[i+6][j+6]
		}
	}
	FinishInit(s.BoardAdd)
	// copy 9-th nonet from show board
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			s.BoardAdd.BoardShow[i][j] = s.BoardMain.BoardShow[i+6][j+6]
		}
	}
}

// Enter to enter value at current position and return the success bool
func (s *BasicSudoku) Enter(val int) bool {
	// if val is invalid or value is not changing - return
	if val > s.Size || s.Board[s.CursorPos.Xpos][s.CursorPos.Ypos] == s.BoardShow[s.CursorPos.Xpos][s.CursorPos.Ypos] || val == s.BoardShow[s.CursorPos.Xpos][s.CursorPos.Ypos] {
		return false
	}

	// clear next moves if there is any
	currAct := s.CurrentAction
	for currAct != len(s.Actions) {
		s.Actions = s.Actions[0:currAct]
	}

	// add new move
	s.Actions = append(s.Actions, Change{s.CursorPos, s.BoardShow[s.CursorPos.Xpos][s.CursorPos.Ypos], val})
	s.CurrentAction = len(s.Actions)
	s.Changed = true

	// check if new value is correct
	success := true
	if s.Board[s.CursorPos.Xpos][s.CursorPos.Ypos] != val {
		success = false
	}
	s.BoardShow[s.CursorPos.Xpos][s.CursorPos.Ypos] = val
	return success
}
func (s *TwoDoku) Enter(val int) bool {
	ret := false
	// enter value to the respective board
	if s.BoardMain.CursorPos.Xpos != -1 {
		ret = s.BoardMain.Enter(val)
	}
	if s.BoardAdd.CursorPos.Xpos != -1 {
		ret = s.BoardAdd.Enter(val)
	}

	// if there is a change - update move sequence
	if s.Display() {
		// clear next moves
		currAct := s.CurrentAction
		for currAct != len(s.Actions) {
			s.Actions = s.Actions[0:currAct]
		}
		// if main 9-th and additional 1-st nonets - there is 1 move in each board
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

// FillSudoku to fill main board with numbers
func (s *BasicSudoku) FillSudoku(position int) bool {
	// recursive stop when reached last box
	if position == s.Size*s.Size {
		return true
	}
	x := position / s.Size
	y := position % s.Size
	// skip filled boxes
	if s.Board[x][y] != 0 {
		return s.FillSudoku(position + 1)
	}

	// calculate available nums for current box
	available := s.AvailableNum(position, s.Board)
	lenAvailable := len(available)
	if lenAvailable == 0 {
		return false
	}

	// loop through available options
	for i := 0; i < lenAvailable; i++ {
		// assign new value
		s.Board[x][y] = available[i]
		// if solution is found
		if s.FillSudoku(position + 1) {
			return true
		}
	}
	// revert value to 0 if solution was not found and return to previous position
	s.Board[x][y] = 0
	return false
}
func (s *DiagonalSudoku) FillSudoku(position int) bool {
	// recursive stop when reached last box
	if position == s.Size*s.Size {
		return true
	}
	x := position / s.Size
	y := position % s.Size
	// skip filled boxes
	if s.Board[x][y] != 0 {
		return s.FillSudoku(position + 1)
	}

	// calculate available nums for current box
	available := s.AvailableNum(position, s.Board)
	lenAvailable := len(available)
	if lenAvailable == 0 {
		return false
	}

	// loop through available options
	for i := 0; i < lenAvailable; i++ {
		// assign new value
		s.Board[x][y] = available[i]
		// if solution is found
		if s.FillSudoku(position + 1) {
			return true
		}
	}
	// revert value to 0 if solution was not found and return to previous position
	s.Board[x][y] = 0
	return false
}

// AvailableNum to calculate all possible numbers for current box
func (s *BasicSudoku) AvailableNum(position int, board [][]int) []int {
	x := position / s.Size
	y := position % s.Size

	// check whether element belong to array
	contains := func(slice []int, value int) bool {
		for _, item := range slice {
			if item == value {
				return true
			}
		}
		return false
	}

	// elements to not put in current box
	taken := []int{}

	// calculate current nonet
	nonetStart := [2]int{x - x%s.NonetSize.Xpos, y - y%s.NonetSize.Ypos}
	// add nonet numbers to taken
	for i := 0; i < s.NonetSize.Xpos; i++ {
		for j := 0; j < s.NonetSize.Ypos; j++ {
			tmp := board[nonetStart[0]+i][nonetStart[1]+j]
			if tmp != 0 && !contains(taken, tmp) {
				taken = append(taken, tmp)
			}
		}
	}

	// add horizontal and vertical numbers to taken
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

	// allowed numbers
	var available []int

	// loop through all numbers and check is they are taken
	for i := 1; i <= s.Size; i++ {
		if !contains(taken, i) {
			available = append(available, i)
		}
	}

	// randomly shuffle numbers
	for i := range available {
		j := rand.Intn(i + 1)
		available[i], available[j] = available[j], available[i]
	}

	return available
}
func (s *DiagonalSudoku) AvailableNum(position int, board [][]int) []int {
	x := position / s.Size
	y := position % s.Size

	// check whether element belong to array
	contains := func(slice []int, value int) bool {
		for _, item := range slice {
			if item == value {
				return true
			}
		}
		return false
	}

	// elements to not put in current box
	taken := []int{}

	// calculate current nonet
	nonetStart := [2]int{x - x%s.NonetSize.Xpos, y - y%s.NonetSize.Ypos}
	// add nonet numbers to taken
	for i := 0; i < s.NonetSize.Xpos; i++ {
		for j := 0; j < s.NonetSize.Ypos; j++ {
			tmp := board[nonetStart[0]+i][nonetStart[1]+j]
			if tmp != 0 && !contains(taken, tmp) {
				taken = append(taken, tmp)
			}
		}
	}

	// add horizontal and vertical numbers to taken
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

	// if current box on left diagonal
	if x == y {
		// add left diagonal numbers to taken
		for i := 0; i < s.Size; i++ {
			tmpLeftDiagonal := board[i][i]
			if tmpLeftDiagonal != 0 && !contains(taken, tmpLeftDiagonal) {
				taken = append(taken, tmpLeftDiagonal)
			}
		}
	}

	// if current box on right diagonal
	if x == s.Size-y-1 {
		// add right diagonal numbers to taken
		for i := 0; i < s.Size; i++ {
			tmpRightDiagonal := board[s.Size-i-1][i]
			if tmpRightDiagonal != 0 && !contains(taken, tmpRightDiagonal) {
				taken = append(taken, tmpRightDiagonal)
			}
		}
	}

	// allowed numbers
	var available []int

	// loop through all numbers and check is they are taken
	for i := 1; i <= s.Size; i++ {
		if !contains(taken, i) {
			available = append(available, i)
		}
	}

	// randomly shuffle numbers
	for i := range available {
		j := rand.Intn(i + 1)
		available[i], available[j] = available[j], available[i]
	}

	return available
}

// number of solutions found
var solutions int

// EmptyGrid to empty show board
func (s *BasicSudoku) EmptyGrid(difficulty int) {
	// calculate how much boxes to empty in percentages
	emptyFinal := (difficulty + 1) * 25
	// copy of grid for removal
	copyGrid := BasicSudoku{}
	finished := false
	empty := 0
	// time of algorithm start
	startTime := time.Now()

	for !finished {
		// calculate random non empty box
		row, col := rand.Intn(s.Size), rand.Intn(s.Size)
		for s.BoardShow[row][col] == 0 {
			row, col = rand.Intn(s.Size), rand.Intn(s.Size)
		}

		// copy current grid and empty new position box
		copyGrid.Copy(s)
		solutions = 0
		copyGrid.BoardShow[row][col] = 0

		// calculate number of solutions
		copyGrid.SolveGrid(0)

		// if there is only 1 solution - accept changes
		if solutions == 1 {
			s.BoardShow[row][col] = 0
			empty++
		}

		// stop either when empty quota is reached or 2 seconds has passed
		if 100*empty/(s.Size*s.Size) > emptyFinal || time.Since(startTime).Seconds() > 2 {
			finished = true
		}

	}

}
func (s *DiagonalSudoku) EmptyGrid(difficulty int) {
	// calculate how much boxes to empty in percentages
	emptyFinal := (difficulty + 1) * 25
	// copy of grid for removal
	copyGrid := DiagonalSudoku{}
	finished := false
	empty := 0
	// time of algorithm start
	startTime := time.Now()

	for !finished {
		// calculate random non empty box
		row, col := rand.Intn(s.Size), rand.Intn(s.Size)
		for s.BoardShow[row][col] == 0 {
			row, col = rand.Intn(s.Size), rand.Intn(s.Size)
		}

		// copy current grid and empty new position box
		copyGrid.Copy(s)
		solutions = 0
		copyGrid.BoardShow[row][col] = 0

		// calculate number of solutions
		copyGrid.SolveGrid(0)

		// if there is only 1 solution - accept changes
		if solutions == 1 {
			s.BoardShow[row][col] = 0
			empty++
		}

		// stop either when empty quota is reached or 2 seconds has passed
		if 100*empty/(s.Size*s.Size) > emptyFinal || time.Since(startTime).Seconds() > 2 {
			finished = true
		}

	}

}

// SolveGrid to calculate number of solutions for current board
func (s *BasicSudoku) SolveGrid(position int) {
	// if reached end - increment solutions and return
	if position == s.Size*s.Size {
		solutions++
		return
	}

	x := position / s.Size
	y := position % s.Size

	// skip non-empty boxes
	if s.BoardShow[x][y] != 0 {
		s.SolveGrid(position + 1)
		return
	}

	// calculate available nums for current box and loop through them
	available := s.AvailableNum(position, s.BoardShow)
	lenAvailable := len(available)
	for i := 0; i < lenAvailable; i++ {
		s.BoardShow[x][y] = available[i]
		// recursively find all the solutions
		s.SolveGrid(position + 1)
	}

	// reset the empty box and return
	s.BoardShow[x][y] = 0
	return
}
func (s *DiagonalSudoku) SolveGrid(position int) {
	// if reached end - increment solutions and return
	if position == s.Size*s.Size {
		solutions++
		return
	}

	x := position / s.Size
	y := position % s.Size

	// skip non-empty boxes
	if s.BoardShow[x][y] != 0 {
		s.SolveGrid(position + 1)
		return
	}

	// calculate available nums for current box and loop through them
	available := s.AvailableNum(position, s.BoardShow)
	lenAvailable := len(available)
	for i := 0; i < lenAvailable; i++ {
		s.BoardShow[x][y] = available[i]
		// recursively find all the solutions
		s.SolveGrid(position + 1)
	}

	// reset the empty box and return
	s.BoardShow[x][y] = 0
	return
}

// SudokuPrintManual to print common parts of sudoku manual for every board
func SudokuPrintManual() {
	blueFont.Println("Move with arrows, enter with numbers 1-9(and A-C, depending on Board Size)")
	greenFont.Print("Green")
	blueFont.Println(" - solved")
	redFont.Print("Red")
	blueFont.Println(" - incorrect")
	purpleFont.Print("Purple")
	blueFont.Println(" - cursor")
}

// Print to output to console the show board as well as timer and instructions
func (s *BasicSudoku) Print() {
	if !s.Changed {
		return
	}
	SudokuPrintManual()
	s.Changed = false
	printFont := fmt.Printf
	// loop through show board
	for i, line := range s.BoardShow {
		// print horizontal borders for nonets
		if i%s.NonetSize.Xpos == 0 {
			// mid-board borders
			if i > 0 {
				blueFont.Print(strings.Repeat("|"+strings.Repeat("_", s.NonetSize.Ypos*2+1), s.Size/s.NonetSize.Ypos))
				blueFont.Println("|")
			} else { // first line border
				blueFont.Print(strings.Repeat("_"+strings.Repeat("_", s.NonetSize.Ypos*2+1), s.Size/s.NonetSize.Ypos))
				blueFont.Println("_")
			}
		}
		for j, element := range line {
			// print vertical borders
			if j%s.NonetSize.Ypos == 0 {
				blueFont.Print("| ")
			}
			// incorrect element
			if element != 0 && s.Board[i][j] != element {
				printFont = redFont.Printf
			} else if s.CursorPos.Xpos == i && s.CursorPos.Ypos == j { // element where cursor is located
				printFont = purpleFont.Printf
			} else if element != 0 { // correct element
				printFont = greenFont.Printf
			} else { // empty element
				printFont = fmt.Printf
			}
			// handle numbers over 9
			if element > 9 {
				// transform to letter format
				_, _ = printFont("%c ", 'A'+(element-10))
			} else {
				_, _ = printFont("%d ", element)
			}
		}
		blueFont.Print("|")
		fmt.Println()
	}
	// print last horizontal border
	blueFont.Print(strings.Repeat("|"+strings.Repeat("_", s.NonetSize.Ypos*2+1), s.Size/s.NonetSize.Ypos))
	blueFont.Println("|")
	// output time left if timer exist
	if s.TimeLeft != -1 {
		minutes := s.TimeLeft / 60
		seconds := s.TimeLeft % 60
		_, _ = greenFont.Printf("Time remaining: %02d:%02d\n", minutes, seconds)
	}
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
	// loop through show board
	for i, line := range s.BoardShow {
		// print horizontal borders for nonets
		if i%s.NonetSize.Xpos == 0 {
			// mid-board borders
			if i > 0 {
				blueFont.Print(strings.Repeat("|"+strings.Repeat("_", s.NonetSize.Ypos*2+1), s.Size/s.NonetSize.Ypos))
				blueFont.Println("|")
			} else { // first line border
				blueFont.Print(strings.Repeat("_"+strings.Repeat("_", s.NonetSize.Ypos*2+1), s.Size/s.NonetSize.Ypos))
				blueFont.Println("_")
			}
		}
		for j, element := range line {
			// print vertical borders
			if j%s.NonetSize.Ypos == 0 {
				blueFont.Print("| ")
			}
			// incorrect element
			if element != 0 && s.Board[i][j] != element {
				printFont = redFont.Printf
			} else if s.CursorPos.Xpos == i && s.CursorPos.Ypos == j { // element where cursor is located
				printFont = purpleFont.Printf
			} else if (i == j || i == s.Size-j-1) && element != 0 { // diagonal element
				printFont = diagonalFont.Printf
			} else if element != 0 { // correct element
				printFont = greenFont.Printf
			} else { // empty element
				printFont = fmt.Printf
			}
			// handle numbers over 9
			if element > 9 {
				// transform to letter format
				_, _ = printFont("%c ", 'A'+(element-10))
			} else {
				_, _ = printFont("%d ", element)
			}
		}
		blueFont.Print("|")
		fmt.Println()
	}
	// print last horizontal border
	blueFont.Print(strings.Repeat("|"+strings.Repeat("_", s.NonetSize.Ypos*2+1), s.Size/s.NonetSize.Ypos))
	blueFont.Println("|")
	// output time left if timer exist
	if s.TimeLeft != -1 {
		minutes := s.TimeLeft / 60
		seconds := s.TimeLeft % 60
		greenFont.Printf("Time remaining: %02d:%02d\n", minutes, seconds)
	}
}
func (s *TwoDoku) Print() {
	if !s.BoardMain.Changed && !s.BoardAdd.Changed {
		return
	}
	SudokuPrintManual()
	s.BoardMain.Changed = false
	s.BoardAdd.Changed = false
	printFont := fmt.Printf
	// print first 6 rows of main board
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
	// draw borders for 5 nonets
	blueFont.Print(strings.Repeat("|_______", 3))
	blueFont.Print("|")
	blueFont.Println(strings.Repeat("________", 2))
	// print last 3 and first 3 rows of main and additional boards as 5 nonets in one line
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
	// draw borders for 5 nonets
	blueFont.Print(strings.Repeat("|_______", 5))
	blueFont.Println("|")
	// print last 6 rows of additional board
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
	// draw last border
	fmt.Print(strings.Repeat("  ", 8))
	blueFont.Print(strings.Repeat("|_______", 3))
	blueFont.Println("|")
	// output time left if timer exist
	if s.BoardMain.TimeLeft != -1 {
		minutes := s.BoardMain.TimeLeft / 60
		seconds := s.BoardMain.TimeLeft % 60
		greenFont.Printf("Time remaining: %02d:%02d\n", minutes, seconds)
	}
}

// Copy to copy another sudoku state
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

// IsComplete returns whether all boxes in show board are filled correctly
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

// Display returns whether there is anything new to display
func (s *BasicSudoku) Display() bool {
	return s.Changed
}
func (s *TwoDoku) Display() bool {
	return s.BoardMain.Changed || s.BoardAdd.Changed
}

// Rules returns sudoku rules
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

// Move to move cursor
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
	// calculate pos from 0 to 15 based on in which board cursor is currently
	if s.BoardMain.CursorPos.Xpos == -1 { // second board
		currentPos = s.BoardAdd.CursorPos
		currentPos.Xpos += 6
		currentPos.Ypos += 6
	} else if s.BoardAdd.CursorPos.Xpos == -1 { // main board
		currentPos = s.BoardMain.CursorPos
	} else { // adjacent nonet
		currentPos = s.BoardMain.CursorPos
	}
	newPos := Vector2{currentPos.Xpos + row, currentPos.Ypos + col}
	// check for invalid new position
	if newPos.Xpos > 14 || newPos.Ypos > 14 || newPos.Xpos < 0 || newPos.Ypos < 0 || (newPos.Xpos > 8 && newPos.Ypos < 6) || (newPos.Xpos < 6 && newPos.Ypos > 8) {
		return
	}
	// assign new positions
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

// Undo to undo last move
func (s *BasicSudoku) Undo() {
	action := s.CurrentAction - 1
	// no previous move
	if action == -1 {
		return
	}
	s.CurrentAction = action
	prevChange := s.Actions[action]
	change := Vector2{s.CursorPos.Xpos - prevChange.Pos.Xpos, s.CursorPos.Ypos - prevChange.Pos.Ypos}
	// move to old position and assign old value
	s.Move(-change.Ypos, -change.Xpos)
	s.BoardShow[s.CursorPos.Xpos][s.CursorPos.Ypos] = prevChange.OldVal
}
func (s *TwoDoku) Undo() {
	action := s.CurrentAction - 1
	// no previous move
	if action == -1 {
		return
	}
	s.CurrentAction = action
	prevChange := s.Actions[action]
	// where move happened
	if prevChange.Main {
		s.BoardMain.Undo()
		s.BoardAdd.CursorPos = Vector2{-1, -1}
	} else {
		s.BoardAdd.Undo()
		s.BoardMain.CursorPos = Vector2{-1, -1}
		// if move was done in the adjacent nonet - undo one more time for another board
		if prevChange.Adjacent {
			s.Undo()
		}
	}
}

// Redo to redo the last undone move
func (s *BasicSudoku) Redo() {
	action := s.CurrentAction
	// no more moves to redo
	if action == len(s.Actions) {
		return
	}
	prevChange := s.Actions[action]
	change := Vector2{s.CursorPos.Xpos - prevChange.Pos.Xpos, s.CursorPos.Ypos - prevChange.Pos.Ypos}
	s.Move(-change.Ypos, -change.Xpos)
	s.BoardShow[s.CursorPos.Xpos][s.CursorPos.Ypos] = prevChange.NewVal
	// step to next move
	s.CurrentAction++
}
func (s *TwoDoku) Redo() {
	action := s.CurrentAction
	// no more moves to redo
	if action == len(s.Actions) {
		return
	}
	prevChange := s.Actions[action]
	s.CurrentAction++
	// check which board the move has occurred
	if prevChange.Main {
		s.BoardMain.Redo()
		s.BoardAdd.CursorPos = Vector2{-1, -1}
		// if move was made in adjacent nonet - redo for another board also
		if prevChange.Adjacent {
			s.Redo()
		}
	} else {
		s.BoardAdd.Redo()
		s.BoardMain.CursorPos = Vector2{-1, -1}
	}
}

// ClearFiles to delete all *.sudo files in local directory
func ClearFiles() error {
	// get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// find all filed
	files, err := filepath.Glob(filepath.Join(currentDir, "*.sudo"))
	if err != nil {
		return err
	}

	// loop and delete them
	for _, file := range files {
		err = os.Remove(file)
		if err != nil {
			return err
		}
	}
	return nil
}

// SaveGame to save game state in .sudo file
func (s *BasicSudoku) SaveGame() error {
	// clear directory
	err := ClearFiles()
	if err != nil {
		return err
	}
	// create file for save
	file, err := os.Create("basic.sudo")
	if err != nil {
		return err
	}
	defer file.Close()

	// Create an encoder and send struct for encoding
	enc := gob.NewEncoder(file)
	err = enc.Encode(s)
	if err != nil {
		dialog.Info(err.Error())
		return err
	}

	return nil
}
func (s *DiagonalSudoku) SaveGame() error {
	// clear directory
	err := ClearFiles()
	if err != nil {
		return err
	}
	// create file for save
	file, err := os.Create("diagonal.sudo")
	if err != nil {
		return err
	}
	defer file.Close()

	// Create an encoder and send struct for encoding
	enc := gob.NewEncoder(file)
	err = enc.Encode(s)
	if err != nil {
		dialog.Info(err.Error())
		return err
	}

	return nil
}
func (s *TwoDoku) SaveGame() error {
	// clear directory
	err := ClearFiles()
	if err != nil {
		return err
	}
	// create file for save
	file, err := os.Create("two.sudo")
	if err != nil {
		return err
	}
	defer file.Close()

	// Create an encoder and send struct for encoding
	enc := gob.NewEncoder(file)
	err = enc.Encode(s)
	if err != nil {
		dialog.Info(err.Error())
		return err
	}

	return nil
}

// RevealRandom to fill random empty box with answer
func (s *BasicSudoku) RevealRandom() {
	row := rand.Intn(s.Size)
	col := rand.Intn(s.Size)
	// function to look for the first empty box and reveal it
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
	// start from random position
	if look(row, col) {
		return
	}
	// if no box found - continue from the start
	look(0, 0)
}
func (s *TwoDoku) RevealRandom() {
	row := rand.Intn(9)
	col := rand.Intn(9)
	// function to look for the first empty box in both boards and reveal it
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
	// start from random position
	if look(row, col) {
		return
	}
	// if no box found - continue from the start
	look(0, 0)
}

// TimePass to decrement seconds from play time
func (s *BasicSudoku) TimePass(sec int) {
	// if not -1(no timer)
	if s.TimeLeft > 0 {
		s.Changed = true
		s.TimeLeft -= sec
		// if went over 0 - no time left
		if s.TimeLeft <= 0 {
			s.TimeLeft = 0
		}
	}
}
func (s *TwoDoku) TimePass(sec int) {
	s.BoardMain.TimePass(sec)
}

// TimeEnd to return whether the timer has ended
func (s *BasicSudoku) TimeEnd() bool {
	if s.TimeLeft == 0 {
		return true
	}
	return false
}
func (s *TwoDoku) TimeEnd() bool {
	return s.BoardMain.TimeEnd()
}
