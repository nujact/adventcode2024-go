package Day

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Cell represents a single character in the input matrix
// value: the character at this position
// starList: maps compass directions to words formed in that direction
type Cell struct {
	value    string
	starList map[string]string
}

// CompassDirections defines all possible directions to search for words
// Used for both 4-letter (XMAS) and 3-letter (MAS) word searches
var CompassDirections = []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}

// Day4 solves the Advent of Code 2024 Day 4 puzzle
// The puzzle involves searching for "XMAS" and "MAS" patterns in a character matrix
// Part 1: Find all occurrences of "XMAS" in any direction
// Part 2: Find all occurrences of "MAS" in diagonal directions around "A" characters
func Day4() {
	fmt.Println("2024 Day4 start")

	// Open file and gather raw inputs
	var inputFileNameBegins = "input" // "input" or "test"
	file, err := os.Open(fmt.Sprintf("../inputs/%s4.txt", inputFileNameBegins))
	if err != nil {
		fmt.Println(err)
		return
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
		return
	}

	fmt.Println("inputMemory size:", len(inputMemory))

	// Convert input to matrices for processing
	inputMatrix := getInputMatrix(inputMemory)
	cellMatrix := getCellMatrix(inputMemory)

	// Print matrix dimensions
	fmt.Printf("inputMatrix[X][Y]: X %d Y %d\n", len(inputMatrix), len(inputMatrix[0]))

	// // Print input matrix
	// for i := 0; i < len(inputMatrix); i++ {
	// 	for j := 0; j < len(inputMatrix[i]); j++ {
	// 		fmt.Print(string(inputMatrix[i][j]))
	// 	}
	// 	fmt.Println()
	// }

	// fmt.Println("\ncellMatrix:")
	// for row := 0; row < len(cellMatrix); row++ {
	// 	for col := 0; col < len(cellMatrix[row]); col++ {
	// 		fmt.Print(cellMatrix[row][col].value)
	// 	}
	// 	fmt.Println()
	// }

	// Process parts
	part1(cellMatrix)
	part2(cellMatrix)

	fmt.Println("2024 Day4 end")
}

// part1 solves the first part of the puzzle
// Searches for the word "XMAS" in all 8 compass directions
// Each cell's starList contains 4-letter words formed in each direction
func part1(cellMatrix [][]Cell) {
	// Calculate 4-letter words in compass directions around each cell
	for row := 0; row < len(cellMatrix); row++ {
		for col := 0; col < len(cellMatrix[row]); col++ {
			//fmt.Printf("%s r,c %d,%d ", cellMatrix[row][col].value, row, col)
			calcStarList(cellMatrix, row, col)
			//fmt.Printf("starList: %d %v\n", len(cellMatrix[row][col].starList), cellMatrix[row][col].starList)
		}
		//fmt.Println()
	}
	//fmt.Println()

	// Count hits of XMAS in all starlists
	xmasCount := 0
	for row := 0; row < len(cellMatrix); row++ {
		for col := 0; col < len(cellMatrix[row]); col++ {
			for _, starWord := range cellMatrix[row][col].starList {
				if starWord == "XMAS" {
					xmasCount++
				}
			}
		}
	}
	fmt.Printf("xmasCount: %d\n\n", xmasCount)
}

// part2 solves the second part of the puzzle
// Searches for "MAS" in diagonal directions around "A" characters
// Counts positions where 2 or more "MAS" words are found in diagonal directions
func part2(cellMatrix [][]Cell) {
	// Calculate 3-letter words in compass directions around each cell
	for row := 0; row < len(cellMatrix); row++ {
		for col := 0; col < len(cellMatrix[row]); col++ {
			//fmt.Printf("%s r,c %d,%d ", cellMatrix[row][col].value, row, col)
			calcMasList(cellMatrix, row, col)
			//fmt.Printf("masList: %d %v\n", len(cellMatrix[row][col].starList), cellMatrix[row][col].starList)
		}
		//fmt.Println()
	}
	//fmt.Println()

	// Count hits of MAS in all starlists
	masCount := 0
	for row := 0; row < len(cellMatrix); row++ {
		for col := 0; col < len(cellMatrix[row]); col++ {
			// Search for A in middle of X-MAS
			if cellMatrix[row][col].value == "A" {
				var nwWord, neWord, swWord, seWord string

				// Check if valid to get corner
				if row-1 >= 0 && col-1 >= 0 {
					nwWord = cellMatrix[row-1][col-1].starList["SE"]
				}
				if row-1 >= 0 && col+1 < len(cellMatrix[row]) {
					neWord = cellMatrix[row-1][col+1].starList["SW"]
				}
				if row+1 < len(cellMatrix) && col-1 >= 0 {
					swWord = cellMatrix[row+1][col-1].starList["NE"]
				}
				if row+1 < len(cellMatrix) && col+1 < len(cellMatrix[row]) {
					seWord = cellMatrix[row+1][col+1].starList["NW"]
				}

				// Count MAS in corners
				foundMas := 0
				if nwWord == "MAS" {
					foundMas++
				}
				if neWord == "MAS" {
					foundMas++
				}
				if swWord == "MAS" {
					foundMas++
				}
				if seWord == "MAS" {
					foundMas++
				}
				// If 2 or more MAS found in corners, increment mas count
				if foundMas >= 2 {
					masCount++
				}
			}
		}
	}

	fmt.Printf("masCount: %d\n\n", masCount)
}

// calcMasList calculates 3-letter words in diagonal directions for a given cell
// Only processes NW, NE, SW, SE directions (diagonals)
func calcMasList(cellMatrix [][]Cell, row, col int) {
	for _, compassDirection := range CompassDirections {
		// Only care about NW NE SW SE
		switch compassDirection {
		case "N", "E", "S", "W":
			continue

		case "NE":
			if row-2 < 0 || col+2 >= len(cellMatrix[row]) {
				continue
			}
			cellMatrix[row][col].starList[compassDirection] = cellMatrix[row][col].value +
				cellMatrix[row-1][col+1].value +
				cellMatrix[row-2][col+2].value

		case "SE":
			if row+2 >= len(cellMatrix) || col+2 >= len(cellMatrix[row]) {
				continue
			}
			cellMatrix[row][col].starList[compassDirection] = cellMatrix[row][col].value +
				cellMatrix[row+1][col+1].value +
				cellMatrix[row+2][col+2].value

		case "SW":
			if row+2 >= len(cellMatrix) || col-2 < 0 {
				continue
			}
			cellMatrix[row][col].starList[compassDirection] = cellMatrix[row][col].value +
				cellMatrix[row+1][col-1].value +
				cellMatrix[row+2][col-2].value

		case "NW":
			if row-2 < 0 || col-2 < 0 {
				continue
			}
			cellMatrix[row][col].starList[compassDirection] = cellMatrix[row][col].value +
				cellMatrix[row-1][col-1].value +
				cellMatrix[row-2][col-2].value
		}
	}
}

// calcStarList calculates 4-letter words in all directions for a given cell
// Processes all 8 compass directions (N, NE, E, SE, S, SW, W, NW)
func calcStarList(cellMatrix [][]Cell, row, col int) {
	for _, compassDirection := range CompassDirections {
		switch compassDirection {
		case "N":
			if row-3 < 0 {
				continue
			}
			cellMatrix[row][col].starList[compassDirection] = cellMatrix[row][col].value +
				cellMatrix[row-1][col].value +
				cellMatrix[row-2][col].value +
				cellMatrix[row-3][col].value

		case "NE":
			if row-3 < 0 || col+3 >= len(cellMatrix[row]) {
				continue
			}
			cellMatrix[row][col].starList[compassDirection] = cellMatrix[row][col].value +
				cellMatrix[row-1][col+1].value +
				cellMatrix[row-2][col+2].value +
				cellMatrix[row-3][col+3].value

		case "E":
			if col+3 >= len(cellMatrix[row]) {
				continue
			}
			cellMatrix[row][col].starList[compassDirection] = cellMatrix[row][col].value +
				cellMatrix[row][col+1].value +
				cellMatrix[row][col+2].value +
				cellMatrix[row][col+3].value

		case "SE":
			if row+3 >= len(cellMatrix) || col+3 >= len(cellMatrix[row]) {
				continue
			}
			cellMatrix[row][col].starList[compassDirection] = cellMatrix[row][col].value +
				cellMatrix[row+1][col+1].value +
				cellMatrix[row+2][col+2].value +
				cellMatrix[row+3][col+3].value

		case "S":
			if row+3 >= len(cellMatrix) {
				continue
			}
			cellMatrix[row][col].starList[compassDirection] = cellMatrix[row][col].value +
				cellMatrix[row+1][col].value +
				cellMatrix[row+2][col].value +
				cellMatrix[row+3][col].value

		case "SW":
			if row+3 >= len(cellMatrix) || col-3 < 0 {
				continue
			}
			cellMatrix[row][col].starList[compassDirection] = cellMatrix[row][col].value +
				cellMatrix[row+1][col-1].value +
				cellMatrix[row+2][col-2].value +
				cellMatrix[row+3][col-3].value

		case "W":
			if col-3 < 0 {
				continue
			}
			cellMatrix[row][col].starList[compassDirection] = cellMatrix[row][col].value +
				cellMatrix[row][col-1].value +
				cellMatrix[row][col-2].value +
				cellMatrix[row][col-3].value

		case "NW":
			if row-3 < 0 || col-3 < 0 {
				continue
			}
			cellMatrix[row][col].starList[compassDirection] = cellMatrix[row][col].value +
				cellMatrix[row-1][col-1].value +
				cellMatrix[row-2][col-2].value +
				cellMatrix[row-3][col-3].value
		}
	}
}

// getInputMatrix converts the input string into a 2D matrix of runes
// Used for initial processing and visualization
func getInputMatrix(inputMemory string) [][]rune {
	// Split by new line into array
	inputRows := strings.Split(inputMemory, "\n")
	// Convert to matrix of characters
	inputMatrix := make([][]rune, len(inputRows))
	for i := 0; i < len(inputRows); i++ {
		inputMatrix[i] = []rune(inputRows[i])
	}
	return inputMatrix
}

// getCellMatrix converts the input string into a 2D matrix of Cells
// Each Cell contains the character value and a map of words in different directions
func getCellMatrix(inputMemory string) [][]Cell {
	// Split by new line into array
	inputRows := strings.Split(inputMemory, "\n")
	// Convert to matrix of Cells
	cellMatrix := make([][]Cell, len(inputRows))
	for row := 0; row < len(inputRows); row++ {
		cellMatrix[row] = make([]Cell, len(inputRows[row]))
		for col := 0; col < len(inputRows[row]); col++ {
			cellMatrix[row][col] = Cell{
				value:    string(inputRows[row][col]),
				starList: make(map[string]string),
			}
			// Initialize starList with empty strings for all directions
			for _, direction := range CompassDirections {
				cellMatrix[row][col].starList[direction] = ""
			}
		}
	}
	return cellMatrix
}
