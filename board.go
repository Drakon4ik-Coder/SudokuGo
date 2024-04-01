package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

type SudokuBoard interface {
	//Init(size int)                // Initialize the board
	//Set(row, col, val int) bool   // Set value at a specific position
	//Get(row, col int) int         // Get value at a specific position
	//RevealRandom()                // Reveal random box
	Enter(row, col, val int) bool // Check if the val is the same as in the board
	IsComplete() bool             // Check if the board is complete
	Print(row, col int)           // Print the board
	Display() bool                // Return whether there were any changes since last call of Print
	GetSize() int
	Rules() string // Get size of the board
}

type Vector2 struct {
	height int
	width  int
}

type BasicSudoku struct {
	boardShow [][]int
	board     [][]int
	size      int
	nonetSize Vector2
	changed   bool
}

type DiagonalSudoku struct {
	BasicSudoku
}

type TwoDoku struct {
	// two boards matrices are stored in three-dimensional array
	board1 BasicSudoku
	board2 BasicSudoku
}

type TriangularSudoku struct {
	BasicSudoku
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

func (s *BasicSudoku) Init(size, difficulty int) {

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
		s.nonetSize.width = int(tmpSize)
		s.nonetSize.height = int(tmpSize)
	} else {
		w, h := findClosestFactors(size)
		s.nonetSize.width = w
		s.nonetSize.height = h
	}
	s.FillSudoku(0)
	for i := 0; i < s.size; i++ {
		for j := 0; j < s.size; j++ {
			s.boardShow[i][j] = s.board[i][j]
		}
	}

	s.EmptyGrid(difficulty)
	print("2")
}
func (s *DiagonalSudoku) Init(size, difficulty int) {

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
		s.nonetSize.width = int(tmpSize)
		s.nonetSize.height = int(tmpSize)
	} else {
		w, h := findClosestFactors(size)
		s.nonetSize.width = w
		s.nonetSize.height = h
	}
	s.FillSudoku(0)
	for i := 0; i < s.size; i++ {
		for j := 0; j < s.size; j++ {
			s.boardShow[i][j] = s.board[i][j]
		}
	}

	s.EmptyGrid(difficulty)
	print("2")
}

func (s *BasicSudoku) Enter(row, col, val int) bool {
	if s.board[row][col] == s.boardShow[row][col] {
		return true
	}

	s.changed = true

	success := true
	if s.board[row][col] != val {
		success = false
	}
	s.boardShow[row][col] = val
	return success
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

	nonetStart := [2]int{x - x%s.nonetSize.width, y - y%s.nonetSize.height}
	for i := 0; i < s.nonetSize.width; i++ {
		for j := 0; j < s.nonetSize.height; j++ {
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

	nonetStart := [2]int{x - x%s.nonetSize.width, y - y%s.nonetSize.height}
	for i := 0; i < s.nonetSize.width; i++ {
		for j := 0; j < s.nonetSize.height; j++ {
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

func SudokuPrint() {
	blueFont.Println("Move with arrows, enter with numbers 1-9(and A-C, depending on board size)")
	greenFont.Print("Green")
	blueFont.Println(" - solved")
	redFont.Print("Red")
	blueFont.Println(" - incorrect")
	purpleFont.Print("Purple")
	blueFont.Println(" - cursor")
}

func (s *BasicSudoku) Print(row, col int) {
	SudokuPrint()
	s.changed = false
	printFont := fmt.Printf
	for i, line := range s.boardShow {
		if i%s.nonetSize.width == 0 {
			if i > 0 {
				blueFont.Print(strings.Repeat("|"+strings.Repeat("_", s.nonetSize.height*2+1), s.size/s.nonetSize.height))
				blueFont.Println("|")
			} else {
				blueFont.Print(strings.Repeat("_"+strings.Repeat("_", s.nonetSize.height*2+1), s.size/s.nonetSize.height))
				blueFont.Println("_")
			}
		}
		for j, element := range line {
			if j%s.nonetSize.height == 0 {
				blueFont.Print("| ")
			}
			if element != 0 && s.board[i][j] != element {
				printFont = redFont.Printf
			} else if row == i && col == j {
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
	blueFont.Print(strings.Repeat("|"+strings.Repeat("_", s.nonetSize.height*2+1), s.size/s.nonetSize.height))
	blueFont.Println("|")
}
func (s *DiagonalSudoku) Print(row, col int) {
	SudokuPrint()
	diagonalFont.Print("Yellow")
	blueFont.Println(" - correct diagonal")

	s.changed = false
	printFont := fmt.Printf
	for i, line := range s.boardShow {
		if i%s.nonetSize.width == 0 {
			if i > 0 {
				blueFont.Print(strings.Repeat("|"+strings.Repeat("_", s.nonetSize.height*2+1), s.size/s.nonetSize.height))
				blueFont.Println("|")
			} else {
				blueFont.Print(strings.Repeat("_"+strings.Repeat("_", s.nonetSize.height*2+1), s.size/s.nonetSize.height))
				blueFont.Println("_")
			}
		}
		for j, element := range line {
			if j%s.nonetSize.height == 0 {
				blueFont.Print("| ")
			}
			if element != 0 && s.board[i][j] != element {
				printFont = redFont.Printf
			} else if row == i && col == j {
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
	blueFont.Print(strings.Repeat("|"+strings.Repeat("_", s.nonetSize.height*2+1), s.size/s.nonetSize.height))
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
	return true
}

func (s *BasicSudoku) Display() bool {
	return s.changed
}

func (s *BasicSudoku) GetSize() int {
	return s.size
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
