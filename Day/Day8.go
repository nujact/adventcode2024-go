// Package Day implements solutions for Advent of Code 2024 challenges.
// Day8 specifically solves a puzzle involving antenna wave interference patterns.
// The puzzle involves a matrix where each cell can contain an antenna broadcasting
// at a specific frequency. When two antennas have the same frequency, they create
// interference points (anti-nodes) at specific positions relative to their line of sight.
package Day

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Day8Cell represents a single cell in the antenna matrix.
// Each cell can contain an antenna with a specific frequency and tracks interference points (anti-nodes)
// from other antennas with matching frequencies.
type Day8Cell struct {
	antennaFrequency string   // The frequency of the antenna (if present)
	antiNodeList     []string // List of frequencies that create interference points in this cell
	brotherList      [][]int  // List of coordinates [row,col] for other antennas with matching frequency
}

// Day8Matrix represents the game board containing cells with antennas and their interference patterns.
// The matrix tracks both the original input and the processed cell data for calculating interference patterns.
type Day8Matrix struct {
	inputStrings []string      // Original input strings representing the raw matrix
	cellMatrix   [][]*Day8Cell // 2D matrix of cells containing antenna and interference data
}

// day8NewMatrix creates a new Matrix from the input string.
// It parses the input string into a 2D matrix, initializing cells with antennas
// where specified in the input. Empty cells are represented by '.' in the input.
//
// Parameters:
//   - inputString: Raw string containing the matrix layout, with '.' for empty cells
//     and other characters representing antenna frequencies
//
// Returns:
//   - *Day8Matrix: Initialized matrix with antenna positions and empty cells
func day8NewMatrix(inputString string) *Day8Matrix {
	inputStrings := strings.Split(inputString, "\n")
	cellMatrix := make([][]*Day8Cell, len(inputStrings))
	for j := range cellMatrix {
		cellMatrix[j] = make([]*Day8Cell, len(inputStrings[0]))
		for i := range cellMatrix[j] {
			cell := &Day8Cell{
				antiNodeList: make([]string, 0),
			}
			inputStr := inputStrings[j][i]
			if inputStr != '.' {
				cell.antennaFrequency = string(inputStr)
				cell.antiNodeList = append(cell.antiNodeList, string(inputStr))
			}
			cellMatrix[j][i] = cell
		}
	}

	return &Day8Matrix{
		inputStrings: inputStrings,
		cellMatrix:   cellMatrix,
	}
}

// Print displays the current state of the matrix to standard output.
// The display uses the following format:
//   - '.' for empty cells
//   - '#' for cells containing interference points (anti-nodes)
//   - The frequency character for cells containing antennas
//
// This method is useful for visualizing the matrix state during processing
// and debugging the interference pattern calculations.
func (m *Day8Matrix) Print() {
	fmt.Print("\nMatrix:\n")
	for _, matrixRow := range m.cellMatrix {
		for _, cell := range matrixRow {
			if cell.antennaFrequency != "" {
				fmt.Print(cell.antennaFrequency)
			} else if len(cell.antiNodeList) > 0 {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// getBrotherList finds all cells containing antennas with the same frequency.
// For a given antenna, this method locates all other antennas in the matrix
// that share the same frequency, which is necessary for calculating interference patterns.
//
// Parameters:
//   - antennaFrequency: The frequency to search for in the matrix
//   - row: Current antenna's row position to exclude from results
//   - col: Current antenna's column position to exclude from results
//
// Returns:
//   - [][]int: List of [row,col] coordinates for matching antennas, excluding the input position
func (m *Day8Matrix) getBrotherList(antennaFrequency string, row, col int) [][]int {
	brotherList := make([][]int, 0)
	for j := range m.cellMatrix {
		for i := range m.cellMatrix[j] {
			if m.cellMatrix[j][i].antennaFrequency == antennaFrequency && (j != row || i != col) {
				brotherList = append(brotherList, []int{j, i})
			}
		}
	}
	return brotherList
}

// calcAntiNodesByWave calculates interference points (anti-nodes) for a specific wave number.
// Anti-nodes are calculated in both directions from the line between two matching antennas.
// The wave number determines how far from the antenna line the anti-nodes appear.
//
// Parameters:
//   - wave: The wave number determining distance from the antenna line
//   - brother: Coordinates [row,col] of the matching antenna
//   - row: Row coordinate of the current antenna
//   - col: Column coordinate of the current antenna
//   - cell: The current antenna cell being processed
//
// Returns:
//   - bool: true if valid anti-nodes were found and added, false if no valid points were found
func (m *Day8Matrix) calcAntiNodesByWave(wave int, brother []int, row, col int, cell *Day8Cell) bool {
	brotherRow := brother[0]
	brotherCol := brother[1]

	// Calculate anti-nodes in both directions
	fromAntiNodeRow := row - ((brotherRow - row) * wave)
	fromAntiNodeCol := col - ((brotherCol - col) * wave)
	foundNodes := false

	// Check and add anti-node in first direction
	if fromAntiNodeRow >= 0 && fromAntiNodeCol >= 0 &&
		fromAntiNodeRow < len(m.cellMatrix) && fromAntiNodeCol < len(m.cellMatrix[0]) {
		m.cellMatrix[fromAntiNodeRow][fromAntiNodeCol].antiNodeList = append(
			m.cellMatrix[fromAntiNodeRow][fromAntiNodeCol].antiNodeList,
			cell.antennaFrequency,
		)
		foundNodes = true
	}

	// Calculate and add anti-node in second direction
	intoAntiNodeRow := brotherRow + ((brotherRow - row) * wave)
	intoAntiNodeCol := brotherCol + ((brotherCol - col) * wave)
	if intoAntiNodeRow >= 0 && intoAntiNodeCol >= 0 &&
		intoAntiNodeRow < len(m.cellMatrix) && intoAntiNodeCol < len(m.cellMatrix[0]) {
		m.cellMatrix[intoAntiNodeRow][intoAntiNodeCol].antiNodeList = append(
			m.cellMatrix[intoAntiNodeRow][intoAntiNodeCol].antiNodeList,
			cell.antennaFrequency,
		)
		foundNodes = true
	}

	return foundNodes
}

// calcAntiNodes processes the entire matrix to find all interference points.
// For each antenna in the matrix, this method:
// 1. Finds all other antennas with matching frequency (brothers)
// 2. Calculates interference points at increasing wave numbers until no valid points are found
// 3. Updates the antiNodeList for each cell where interference occurs
//
// This is the main processing function that identifies all interference patterns in the matrix.
func (m *Day8Matrix) calcAntiNodes() {
	fmt.Println("CalcAntiNodes")
	for row := range m.cellMatrix {
		for col := range m.cellMatrix[row] {
			cell := m.cellMatrix[row][col]
			if cell.antennaFrequency != "" {
				// Cell has an antenna, find its brothers
				// fmt.Printf("\n%d,%d antennaFreq %s", row, col, cell.antennaFrequency)
				cell.brotherList = m.getBrotherList(cell.antennaFrequency, row, col)
				// fmt.Printf(" brothers: %d", len(cell.brotherList))

				// Calculate anti-nodes for each brother
				for _, brother := range cell.brotherList {
					wave := 1
					for m.calcAntiNodesByWave(wave, brother, row, col, cell) {
						wave++
					}
				}
			}
		}
	}
}

// Day8 solves the Advent of Code 2024 Day 8 puzzle.
// The puzzle involves finding interference patterns in a matrix of antennas.
// Each antenna broadcasts at a specific frequency, and interference points (anti-nodes)
// occur at specific positions relative to pairs of antennas with matching frequencies.
//
// The solution follows these steps:
// 1. Read the input matrix from file
// 2. Create a matrix structure with antenna positions
// 3. Calculate all interference points
// 4. Count the total number of cells containing interference points
//
// The final answer is the count of cells that contain at least one interference point.
func Day8() {
	fmt.Println("2024 Day8 start")

	// Open file and gather raw inputs
	inputMemory, shouldReturn := day8GetInput()
	if shouldReturn {
		return
	}

	fmt.Println("inputMemory size:", len(inputMemory))

	// Create matrix from input
	matrix := day8NewMatrix(inputMemory)
	matrix.calcAntiNodes()

	// matrix.Print()

	// Count cells with anti-nodes
	countAntiNodes := 0
	for row := range matrix.cellMatrix {
		for col := range matrix.cellMatrix[row] {
			cell := matrix.cellMatrix[row][col]
			if len(cell.antiNodeList) > 0 {
				countAntiNodes++
				// fmt.Printf(" %d,%d %s antiNodes: %v\n",
				// 	row, col, cell.antennaFrequency, cell.antiNodeList)
			}
		}
	}
	fmt.Printf("countAntiNodes: %d\n", countAntiNodes)

	fmt.Println("2024 Day8 end")
}

// day8GetInput reads the puzzle input file and returns its contents.
// The input file should contain a matrix where:
//   - '.' represents empty cells
//   - Any other character represents an antenna with that frequency
//
// Returns:
//   - string: The contents of the input file
//   - bool: true if an error occurred, false otherwise
//
// The function handles file operations and error checking, returning appropriate values
// to indicate success or failure of the input reading process.
func day8GetInput() (string, bool) {
	var inputFileNameBegins = "input" // "input" or "test"
	file, err := os.Open(fmt.Sprintf("inputs/%s8.txt", inputFileNameBegins))
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
