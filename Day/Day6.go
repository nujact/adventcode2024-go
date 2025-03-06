package Day

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Day6Cell represents a single cell in the matrix
// It tracks whether the cell is obstructed, has been visited,
// and how many times it has been visited from each direction
type Day6Cell struct {
	obstructed bool // Whether the cell contains an obstacle (#)
	visited    bool // Whether the cell has been visited at all
	visitedN   int  // Number of times visited from North direction
	visitedE   int  // Number of times visited from East direction
	visitedS   int  // Number of times visited from South direction
	visitedW   int  // Number of times visited from West direction
}

// Guard represents the moving guard in the matrix
// It tracks its position, current direction, and whether it's in a death loop
type Guard struct {
	row       int    // Current row position in the matrix
	col       int    // Current column position in the matrix
	direction string // Current direction of movement (N, E, S, W)
	deathLoop bool   // Whether the guard is stuck in a death loop
}

// Matrix represents the game board and contains the guard
// It stores the input strings, the cell matrix, and the guard
type Matrix struct {
	inputStrings []string     // Original input strings
	cellMatrix   [][]Day6Cell // 2D matrix of cells
	guard        *Guard       // The moving guard
}

// NewMatrix creates a new Matrix from the input string
// It initializes the cell matrix, sets up obstacles, and finds the guard's starting position
func NewMatrix(inputString string) *Matrix {
	inputStrings := strings.Split(inputString, "\n")
	cellMatrix := make([][]Day6Cell, len(inputStrings))
	for j := range cellMatrix {
		cellMatrix[j] = make([]Day6Cell, len(inputStrings[0]))
		for i := range cellMatrix[j] {
			cellMatrix[j][i] = Day6Cell{}
			inputStr := inputStrings[j][i]
			if inputStr == '#' {
				cellMatrix[j][i].obstructed = true
			}
		}
	}

	// Find guard position and initialize it
	guard := &Guard{}
	for j := range cellMatrix {
		for i := range cellMatrix[j] {
			if inputStrings[j][i] == '^' {
				guard.row = j
				guard.col = i
				guard.direction = "N" // Guard starts facing North
				cellMatrix[j][i].visited = true
			}
		}
	}

	return &Matrix{
		inputStrings: inputStrings,
		cellMatrix:   cellMatrix,
		guard:        guard,
	}
}

// Print displays the current state of the matrix
// Uses different characters to represent:
// # - obstacles
// ^ - guard's current position
// x - visited cells
// . - unvisited cells
func (m *Matrix) Print() {
	fmt.Print("\nMatrix:\n")
	for j := range m.cellMatrix {
		for i := range m.cellMatrix[j] {
			cell := m.cellMatrix[j][i]
			if cell.obstructed {
				fmt.Print("#")
			} else if m.guard.row == j && m.guard.col == i {
				fmt.Print("^")
			} else if cell.visited {
				fmt.Print("x")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

// CellReset resets all cells' visited states
// Used when testing different scenarios in part 2
func (m *Matrix) CellReset() {
	for j := range m.cellMatrix {
		for i := range m.cellMatrix[j] {
			m.cellMatrix[j][i].visitedW = 0
			m.cellMatrix[j][i].visitedE = 0
			m.cellMatrix[j][i].visitedN = 0
			m.cellMatrix[j][i].visitedS = 0
			m.cellMatrix[j][i].visited = false
		}
	}
}

// MoveGuard moves the guard according to its current direction
// Returns true if the move was valid and the guard should continue moving
// The guard's movement rules are:
// 1. Move in current direction if possible
// 2. If hitting an obstacle, back up and turn right
// 3. If visiting a cell too many times in same direction, enter death loop
func (m *Matrix) MoveGuard() bool {
	validMove := true
	m.guard.deathLoop = false

	// Check if move is within bounds
	switch m.guard.direction {
	case "N":
		if m.guard.row-1 < 0 {
			validMove = false
		}
	case "E":
		if m.guard.col+1 >= len(m.cellMatrix[0]) {
			validMove = false
		}
	case "S":
		if m.guard.row+1 >= len(m.cellMatrix) {
			validMove = false
		}
	case "W":
		if m.guard.col-1 < 0 {
			validMove = false
		}
	}

	if validMove {
		// Move guard
		switch m.guard.direction {
		case "N":
			m.guard.row--
		case "E":
			m.guard.col++
		case "S":
			m.guard.row++
		case "W":
			m.guard.col--
		}

		cell := &m.cellMatrix[m.guard.row][m.guard.col]
		if cell.obstructed {
			// Backup and turn right
			switch m.guard.direction {
			case "N":
				m.guard.row++
				m.guard.direction = "E"
			case "E":
				m.guard.col--
				m.guard.direction = "S"
			case "S":
				m.guard.row--
				m.guard.direction = "W"
			case "W":
				m.guard.col++
				m.guard.direction = "N"
			}
		} else {
			cell.visited = true
			switch m.guard.direction {
			case "N":
				cell.visitedN++
			case "E":
				cell.visitedE++
			case "S":
				cell.visitedS++
			case "W":
				cell.visitedW++
			}

			// Check for death loop
			if cell.visitedN > 1 || cell.visitedE > 1 || cell.visitedS > 1 || cell.visitedW > 1 {
				m.guard.deathLoop = true
				validMove = false
			}
		}
	}

	return validMove
}

// Day6 solves the Advent of Code 2024 Day 6 puzzle
// The puzzle involves a guard moving around a matrix and potentially getting stuck in death loops
// Part 1: Count the number of cells visited by the guard
// Part 2: Count the number of cells that can cause a death loop when blocked
func Day6() {
	fmt.Println("2024 Day6 start")

	// Open file and gather raw inputs
	inputMemory, shouldReturn := day6GetInput()
	if shouldReturn {
		return
	}

	fmt.Println("inputMemory size:", len(inputMemory))

	// Create matrix from input
	matrix := NewMatrix(inputMemory)

	// Process parts
	day6part1(matrix)
	matrix = NewMatrix(inputMemory) // Reset matrix for part 2
	day6part2(matrix)

	fmt.Println("2024 Day6 end")
}

// day6part1 processes part 1 of the puzzle
// Counts the number of cells visited by the guard as it moves around
// The guard moves until it can't move anymore or enters a death loop
func day6part1(matrix *Matrix) {
	// matrix.Print()

	for matrix.MoveGuard() {
		// Continue moving until guard can't move anymore
	}

	// matrix.Print()
	// fmt.Printf("Guard %d %d %s\n", matrix.guard.row, matrix.guard.col, matrix.guard.direction)

	// Count visited cells
	visitedCells := 0
	for j := range matrix.cellMatrix {
		for i := range matrix.cellMatrix[j] {
			if matrix.cellMatrix[j][i].visited {
				visitedCells++
			}
		}
	}
	fmt.Printf("visited count %d\n", visitedCells)
}

// day6part2 processes part 2 of the puzzle
// For each non-obstructed cell, tests if blocking it would cause a death loop
// A death loop occurs when the guard visits a cell too many times in the same direction
func day6part2(matrix *Matrix) {
	// matrix.Print()
	guardRowStart := matrix.guard.row
	guardColStart := matrix.guard.col

	deathLoopCount := 0
	for j := range matrix.cellMatrix {
		for i := range matrix.cellMatrix[j] {
			cell := &matrix.cellMatrix[j][i]
			matrix.guard.direction = "N"
			matrix.guard.row = guardRowStart
			matrix.guard.col = guardColStart
			matrix.guard.deathLoop = false
			matrix.CellReset()

			// Test if blocking this cell causes a death loop
			if !cell.obstructed && (matrix.guard.row != j || matrix.guard.col != i) {
				cell.obstructed = true
				for matrix.MoveGuard() {
					// Continue moving until guard can't move anymore
				}
				if matrix.guard.deathLoop {
					deathLoopCount++
				}
				cell.obstructed = false
			}
		}
	}

	// matrix.Print()
	fmt.Printf("deathloop count %d\n", deathLoopCount)
}

// day6GetInput reads the input file and returns its contents as a string
// Returns the input string and a boolean indicating if there was an error
// The input file should contain a matrix with:
// # - obstacles
// ^ - guard's starting position
// . - empty cells
func day6GetInput() (string, bool) {
	var inputFileNameBegins = "input" // "input" or "test"
	file, err := os.Open(fmt.Sprintf("../inputs/%s6.txt", inputFileNameBegins))
	if err != nil {
		fmt.Println(err)
		return "", true
	}
	defer file.Close()

	// Read all lines into memory
	var inputMemory string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		inputMemory += scanner.Text() + "\n"
	}
	// remove last \n
	inputMemory = inputMemory[:len(inputMemory)-1]

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return "", true
	}
	return inputMemory, false
}
