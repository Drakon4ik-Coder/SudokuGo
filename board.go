package main

import (
	"container/list"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

type SudokuBoard interface {
	//RevealRandom()                // Reveal random box
	Enter(val int) bool // Check if the val is the same as in the board
	IsComplete() bool   // Check if the board is complete
	Print()             // Print the board
	Move(col, row int)  // Move the cursor if possible
	Rules() string      // Return rules of sudoku
	Display() bool      // Return whether there were any changes since last call of Print
	Undo()              // Undoes previous move
	Redo()              // Cancels last redo call
}

type Vector2 struct {
	x int
	y int
}
type Change struct {
	pos    Vector2
	oldVal int
	newVal int
}

type BasicSudoku struct {
	boardShow     [][]int
	board         [][]int
	size          int
	nonetSize     Vector2
	changed       bool
	cursorPos     Vector2
	actions       *list.List
	currentAction *list.Element
}

type DiagonalSudoku struct {
	BasicSudoku
}

type TwoDoku struct {
	boardMain BasicSudoku
	boardAdd  BasicSudoku
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

func (s *BasicSudoku) PreInit(size int) {
	s.changed = true

	// Initialize board with zeros
	createBoard := func() [][]int {
		tmp := make([][]int, size)
		for i := range tmp {
			tmp[i] = make([]int, size)
		}
		return tmp
	}
	s.board = createBoard()
	s.boardShow = createBoard()
	s.size = size
	tmpSize := math.Sqrt(float64(size))
	if tmpSize == math.Trunc(tmpSize) {
		s.nonetSize.x = int(tmpSize)
		s.nonetSize.y = int(tmpSize)
	} else {
		w, h := findClosestFactors(size)
		s.nonetSize.x = w
		s.nonetSize.y = h
	}

	s.cursorPos = Vector2{0, 0}
	s.actions = list.New()

}

func (s *BasicSudoku) Init(size, difficulty int) {

	s.PreInit(size)

	s.FillSudoku(0)
	for i := 0; i < s.size; i++ {
		for j := 0; j < s.size; j++ {
			s.boardShow[i][j] = s.board[i][j]
		}
	}

	s.EmptyGrid(difficulty)
}
func (s *DiagonalSudoku) Init(size, difficulty int) {

	s.PreInit(size)

	s.FillSudoku(0)
	for i := 0; i < s.size; i++ {
		for j := 0; j < s.size; j++ {
			s.boardShow[i][j] = s.board[i][j]
		}
	}

	s.EmptyGrid(difficulty)
}
func (s *TwoDoku) Init(size, difficulty int) {

	s.boardMain.PreInit(size)
	s.boardAdd.PreInit(size)
	s.boardAdd.cursorPos = Vector2{-1, -1}

	FinishInit := func(board BasicSudoku) {
		board.FillSudoku(0)
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				board.boardShow[i][j] = board.board[i][j]
			}
		}
		board.EmptyGrid(difficulty)
	}
	FinishInit(s.boardMain)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			s.boardAdd.board[i][j] = s.boardMain.board[i+6][j+6]
		}
	}
	FinishInit(s.boardAdd)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			s.boardAdd.boardShow[i][j] = s.boardMain.boardShow[i+6][j+6]
		}
	}
}

func (s *BasicSudoku) Enter(val int) bool {
	if val > s.size || s.board[s.cursorPos.x][s.cursorPos.y] == s.boardShow[s.cursorPos.x][s.cursorPos.y] || val == s.boardShow[s.cursorPos.x][s.cursorPos.y] {
		return true
	}

	currAct := s.currentAction
	for currAct != nil {
		nextAct := currAct.Next()
		s.actions.Remove(currAct)
		currAct = nextAct
	}

	s.actions.PushBack(Change{s.cursorPos, s.boardShow[s.cursorPos.x][s.cursorPos.y], val})
	s.currentAction = nil
	s.changed = true

	success := true
	if s.board[s.cursorPos.x][s.cursorPos.y] != val {
		success = false
	}
	s.boardShow[s.cursorPos.x][s.cursorPos.y] = val
	return success
}
func (s *TwoDoku) Enter(val int) bool {
	if s.boardAdd.cursorPos.x == -1 {
		return s.boardMain.Enter(val)
	} else if s.boardMain.cursorPos.x == -1 {
		return s.boardAdd.Enter(val)
	} else {
		s.boardAdd.Enter(val)
		return s.boardMain.Enter(val)
	}
}

func (s *BasicSudoku) FillSudoku(position int) bool {
	if position == s.size*s.size {
		return true
	}
	x := position / s.size
	y := position % s.size
	if s.board[x][y] != 0 {
		return s.FillSudoku(position + 1)
	}

	available := s.AvailableNum(position, s.board)
	lenAvailable := len(available)
	if lenAvailable == 0 {
		return false
	}

	for i := 0; i < lenAvailable; i++ {
		s.board[x][y] = available[i]
		if s.FillSudoku(position + 1) {
			return true
		}
	}
	s.board[x][y] = 0
	return false
}
func (s *DiagonalSudoku) FillSudoku(position int) bool {
	if position == s.size*s.size {
		return true
	}
	x := position / s.size
	y := position % s.size
	if s.board[x][y] != 0 {
		return s.FillSudoku(position + 1)
	}

	available := s.AvailableNum(position, s.board)
	lenAvailable := len(available)
	if lenAvailable == 0 {
		return false
	}

	for i := 0; i < lenAvailable; i++ {
		s.board[x][y] = available[i]
		if s.FillSudoku(position + 1) {
			return true
		}
	}
	s.board[x][y] = 0
	return false
}

func (s *BasicSudoku) AvailableNum(position int, board [][]int) []int {
	x := position / s.size
	y := position % s.size

	contains := func(slice []int, value int) bool {
		for _, item := range slice {
			if item == value {
				return true
			}
		}
		return false
	}

	taken := []int{}

	nonetStart := [2]int{x - x%s.nonetSize.x, y - y%s.nonetSize.y}
	for i := 0; i < s.nonetSize.x; i++ {
		for j := 0; j < s.nonetSize.y; j++ {
			tmp := board[nonetStart[0]+i][nonetStart[1]+j]
			if tmp != 0 && !contains(taken, tmp) {
				taken = append(taken, tmp)
			}
		}
	}

	for i := 0; i < s.size; i++ {
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

	for i := 1; i <= s.size; i++ {
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
	x := position / s.size
	y := position % s.size

	contains := func(slice []int, value int) bool {
		for _, item := range slice {
			if item == value {
				return true
			}
		}
		return false
	}

	taken := []int{}

	nonetStart := [2]int{x - x%s.nonetSize.x, y - y%s.nonetSize.y}
	for i := 0; i < s.nonetSize.x; i++ {
		for j := 0; j < s.nonetSize.y; j++ {
			tmp := board[nonetStart[0]+i][nonetStart[1]+j]
			if tmp != 0 && !contains(taken, tmp) {
				taken = append(taken, tmp)
			}
		}
	}

	for i := 0; i < s.size; i++ {
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
		for i := 0; i < s.size; i++ {
			tmpLeftDiagonal := board[i][i]
			if tmpLeftDiagonal != 0 && !contains(taken, tmpLeftDiagonal) {
				taken = append(taken, tmpLeftDiagonal)
			}
		}
	}

	if x == s.size-y-1 {
		for i := 0; i < s.size; i++ {
			tmpRightDiagonal := board[s.size-i-1][i]
			if tmpRightDiagonal != 0 && !contains(taken, tmpRightDiagonal) {
				taken = append(taken, tmpRightDiagonal)
			}
		}
	}

	var available []int

	for i := 1; i <= s.size; i++ {
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
		row, col := rand.Intn(s.size), rand.Intn(s.size)
		for s.boardShow[row][col] == 0 {
			row, col = rand.Intn(s.size), rand.Intn(s.size)
		}

		copyGrid.Copy(s)
		solutions = 0
		copyGrid.boardShow[row][col] = 0

		copyGrid.SolveGrid(0)

		if solutions == 1 {
			s.boardShow[row][col] = 0
			empty++
		}

		if 100*empty/(s.size*s.size) > emptyFinal || time.Since(startTime).Seconds() > 2 {
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
		row, col := rand.Intn(s.size), rand.Intn(s.size)
		for s.boardShow[row][col] == 0 {
			row, col = rand.Intn(s.size), rand.Intn(s.size)
		}

		copyGrid.Copy(s)
		solutions = 0
		copyGrid.boardShow[row][col] = 0

		copyGrid.SolveGrid(0)

		if solutions == 1 {
			s.boardShow[row][col] = 0
			empty++
		}

		if 100*empty/(s.size*s.size) > emptyFinal || time.Since(startTime).Seconds() > 2 {
			finished = true
		}

	}

}

func (s *BasicSudoku) SolveGrid(position int) {
	if position == s.size*s.size {
		solutions++
		return
	}

	x := position / s.size
	y := position % s.size

	if s.boardShow[x][y] != 0 {
		s.SolveGrid(position + 1)
		return
	}

	available := s.AvailableNum(position, s.boardShow)
	lenAvailable := len(available)

	for i := 0; i < lenAvailable; i++ {
		s.boardShow[x][y] = available[i]
		s.SolveGrid(position + 1)
	}

	s.boardShow[x][y] = 0
	return
}
func (s *DiagonalSudoku) SolveGrid(position int) {
	if position == s.size*s.size {
		solutions++
		return
	}

	x := position / s.size
	y := position % s.size

	if s.boardShow[x][y] != 0 {
		s.SolveGrid(position + 1)
		return
	}

	available := s.AvailableNum(position, s.boardShow)
	lenAvailable := len(available)

	for i := 0; i < lenAvailable; i++ {
		s.boardShow[x][y] = available[i]
		s.SolveGrid(position + 1)
	}

	s.boardShow[x][y] = 0
	return
}

func SudokuPrintManual() {
	blueFont.Println("Move with arrows, enter with numbers 1-9(and A-C, depending on board size)")
	greenFont.Print("Green")
	blueFont.Println(" - solved")
	redFont.Print("Red")
	blueFont.Println(" - incorrect")
	purpleFont.Print("Purple")
	blueFont.Println(" - cursor")
}

func (s *BasicSudoku) Print() {
	if !s.changed {
		return
	}
	SudokuPrintManual()
	s.changed = false
	printFont := fmt.Printf
	for i, line := range s.boardShow {
		if i%s.nonetSize.x == 0 {
			if i > 0 {
				blueFont.Print(strings.Repeat("|"+strings.Repeat("_", s.nonetSize.y*2+1), s.size/s.nonetSize.y))
				blueFont.Println("|")
			} else {
				blueFont.Print(strings.Repeat("_"+strings.Repeat("_", s.nonetSize.y*2+1), s.size/s.nonetSize.y))
				blueFont.Println("_")
			}
		}
		for j, element := range line {
			if j%s.nonetSize.y == 0 {
				blueFont.Print("| ")
			}
			if element != 0 && s.board[i][j] != element {
				printFont = redFont.Printf
			} else if s.cursorPos.x == i && s.cursorPos.y == j {
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
	blueFont.Print(strings.Repeat("|"+strings.Repeat("_", s.nonetSize.y*2+1), s.size/s.nonetSize.y))
	blueFont.Println("|")
	//for e := s.actions.Front(); e != nil; e = e.Next() {
	//	redFont.Println(e.Value)
	//}
	//redFont.Println()
	//if s.currentAction != nil {
	//	redFont.Println(s.currentAction.Value)
	//}
}
func (s *DiagonalSudoku) Print() {
	if !s.changed {
		return
	}
	SudokuPrintManual()
	diagonalFont.Print("Yellow")
	blueFont.Println(" - correct diagonal")

	s.changed = false
	printFont := fmt.Printf
	for i, line := range s.boardShow {
		if i%s.nonetSize.x == 0 {
			if i > 0 {
				blueFont.Print(strings.Repeat("|"+strings.Repeat("_", s.nonetSize.y*2+1), s.size/s.nonetSize.y))
				blueFont.Println("|")
			} else {
				blueFont.Print(strings.Repeat("_"+strings.Repeat("_", s.nonetSize.y*2+1), s.size/s.nonetSize.y))
				blueFont.Println("_")
			}
		}
		for j, element := range line {
			if j%s.nonetSize.y == 0 {
				blueFont.Print("| ")
			}
			if element != 0 && s.board[i][j] != element {
				printFont = redFont.Printf
			} else if s.cursorPos.x == i && s.cursorPos.y == j {
				printFont = purpleFont.Printf
			} else if (i == j || i == s.size-j-1) && element != 0 {
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
	blueFont.Print(strings.Repeat("|"+strings.Repeat("_", s.nonetSize.y*2+1), s.size/s.nonetSize.y))
	blueFont.Println("|")
}
func (s *TwoDoku) Print() {
	if !s.boardMain.changed && !s.boardAdd.changed {
		return
	}
	SudokuPrintManual()
	s.boardMain.changed = false
	s.boardAdd.changed = false
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
			element := s.boardMain.boardShow[i][j]
			if element != 0 && s.boardMain.board[i][j] != element {
				printFont = redFont.Printf
			} else if s.boardMain.cursorPos.x == i && s.boardMain.cursorPos.y == j {
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
				element = s.boardMain.boardShow[i][j]
				value = s.boardMain.board[i][j]
			} else {
				element = s.boardAdd.boardShow[i-6][j%9+3]
				value = s.boardAdd.board[i-6][j%9+3]
			}

			if j%3 == 0 {
				blueFont.Print("| ")
			}

			if element != 0 && value != element {
				printFont = redFont.Printf
			} else if s.boardMain.cursorPos.x == i && s.boardMain.cursorPos.y == j || s.boardAdd.cursorPos.x == i-6 && s.boardAdd.cursorPos.y == j-6 {
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
			element := s.boardAdd.boardShow[i][j]
			if element != 0 && s.boardAdd.board[i][j] != element {
				printFont = redFont.Printf
			} else if s.boardAdd.cursorPos.x == i && s.boardAdd.cursorPos.y == j {
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
}

func (s *BasicSudoku) Copy(s2 *BasicSudoku) {
	s.size = s2.size
	s.nonetSize = s2.nonetSize

	createBoard := func() [][]int {
		tmp := make([][]int, s.size)
		for i := range tmp {
			tmp[i] = make([]int, s.size)
		}
		return tmp
	}
	s.board = createBoard()
	s.boardShow = createBoard()

	for i := 0; i < s.size; i++ {
		for j := 0; j < s.size; j++ {
			s.board[i][j] = s2.board[i][j]
			s.boardShow[i][j] = s2.boardShow[i][j]
		}
	}

}
func (s *DiagonalSudoku) Copy(s2 *DiagonalSudoku) {
	s.size = s2.size
	s.nonetSize = s2.nonetSize

	createBoard := func() [][]int {
		tmp := make([][]int, s.size)
		for i := range tmp {
			tmp[i] = make([]int, s.size)
		}
		return tmp
	}
	s.board = createBoard()
	s.boardShow = createBoard()

	for i := 0; i < s.size; i++ {
		for j := 0; j < s.size; j++ {
			s.board[i][j] = s2.board[i][j]
			s.boardShow[i][j] = s2.boardShow[i][j]
		}
	}

}

func (s *BasicSudoku) IsComplete() bool {
	for i, line := range s.boardShow {
		for j, element := range line {
			if element == 0 || s.board[i][j] != element {
				return false
			}
		}
	}
	s.changed = true
	s.cursorPos = Vector2{-1, -1}
	return true
}
func (s *TwoDoku) IsComplete() bool {
	return s.boardMain.IsComplete() && s.boardAdd.IsComplete()
}

func (s *BasicSudoku) Display() bool {
	return s.changed
}
func (s *TwoDoku) Display() bool {
	return s.boardMain.changed || s.boardAdd.changed
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
		"that share the same 3x3 region(9-th and 1-st for the corresponding board).\n" +
		"Each row, column, and 3x3(can differ) region must contain all digits\n" +
		"from 1 to 9(possibly up to C) without repetition.\n" +
		"Use logic to fill in the empty cells based on the filled cells.\n" +
		"No guessing is allowed, and each puzzle has exactly one unique solution."
}

func (s *BasicSudoku) Move(col, row int) {
	newPos := Vector2{s.cursorPos.x + row, s.cursorPos.y + col}
	if newPos.x < s.size && newPos.y < s.size && newPos.x >= 0 && newPos.y >= 0 {
		s.cursorPos.x = newPos.x
		s.cursorPos.y = newPos.y
		s.changed = true
	}
}
func (s *TwoDoku) Move(col, row int) {
	var currentPos Vector2
	if s.boardMain.cursorPos.x == -1 {
		currentPos = s.boardAdd.cursorPos
		currentPos.x += 6
		currentPos.y += 6
	} else if s.boardAdd.cursorPos.x == -1 {
		currentPos = s.boardMain.cursorPos
	} else {
		currentPos = s.boardMain.cursorPos
	}
	newPos := Vector2{currentPos.x + row, currentPos.y + col}
	// invalid position
	if newPos.x > 14 || newPos.y > 14 || newPos.x < 0 || newPos.y < 0 || (newPos.x > 8 && newPos.y < 6) || (newPos.x < 6 && newPos.y > 8) {
		return
	}
	s.boardAdd.cursorPos.x = newPos.x - 6
	s.boardAdd.cursorPos.y = newPos.y - 6
	s.boardMain.cursorPos.x = newPos.x
	s.boardMain.cursorPos.y = newPos.y
	if (newPos.x <= 8 && newPos.y < 6) || (newPos.x < 6 && newPos.y <= 8) {
		s.boardAdd.cursorPos = Vector2{-1, -1}

	} else if (newPos.x > 8 && newPos.y >= 6) || (newPos.x >= 6 && newPos.y > 8) {
		s.boardMain.cursorPos = Vector2{-1, -1}
	}
	s.boardMain.changed = true
}

func (s *BasicSudoku) Undo() {
	action := s.currentAction
	if action == nil {
		action = s.actions.Back()
	} else {
		action = action.Prev()
	}
	if action == nil {
		return
	}
	s.currentAction = action
	prevChange := action.Value.(Change)
	change := Vector2{s.cursorPos.x - prevChange.pos.x, s.cursorPos.y - prevChange.pos.y}
	//dialog.Info(fmt.Sprintf("%d %d", change.x, change.y))
	s.Move(-change.y, -change.x)
	s.boardShow[s.cursorPos.x][s.cursorPos.y] = prevChange.oldVal
}
func (s *TwoDoku) Undo() {
	//TODO implement me
	panic("implement me")
}

func (s *BasicSudoku) Redo() {
	action := s.currentAction
	if action == nil {
		return
	}
	prevChange := action.Value.(Change)
	change := Vector2{s.cursorPos.x - prevChange.pos.x, s.cursorPos.y - prevChange.pos.y}
	//dialog.Info(fmt.Sprintf("%d %d", change.x, change.y))
	s.Move(-change.y, -change.x)
	s.boardShow[s.cursorPos.x][s.cursorPos.y] = prevChange.newVal

	s.currentAction = action.Next()
}
func (s *TwoDoku) Redo() {
	//TODO implement me
	panic("implement me")
}
