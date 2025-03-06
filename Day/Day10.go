// Package Day implements solutions for the Advent of Code 2024 challenges.
package Day

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// Day10Position represents a coordinate position in the game map.
// It stores both row and column indices.
type Day10Position struct {
	row int
	col int
}

// Day10Trail represents a complete trail from a starting position (head)
// to an ending position (tail). Each trail must start at elevation 0
// and end at elevation 9.
type Day10Trail struct {
	headRow, headCol int // Starting position coordinates
	tailRow, tailCol int // Ending position coordinates
}

// Day10 solves the Day 10 puzzle of Advent of Code 2024.
// The puzzle involves finding valid trails in a map where:
// - Each position has an elevation from 0-9
// - Trails must start at elevation 0
// - Trails must end at elevation 9
// - Each step must increase elevation by exactly 1
func Day10() {
	fmt.Println("2024 Day 10 start")

	// Read input file
	file, err := os.Open("../inputs/input10.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var inputLines []string
	for scanner.Scan() {
		inputLines = append(inputLines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Build map of inputs
	rows := len(inputLines)
	if rows == 0 {
		return
	}
	cols := len(inputLines[0])

	// Create 2D slice for the map
	gameMap := make([][]int, rows)
	for i := range gameMap {
		gameMap[i] = make([]int, cols)
		for j, ch := range inputLines[i] {
			num, _ := strconv.Atoi(string(ch))
			gameMap[i][j] = num
		}
	}

	// Print map
	fmt.Printf("map size: %dx%d\n", rows, cols)
	// for _, row := range gameMap {
	// 	for _, cell := range row {
	// 		fmt.Print(cell)
	// 	}
	// 	fmt.Println()
	// }

	// Find all trail heads (cells with value 0)
	var trailHeads []Day10Position
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if gameMap[i][j] == 0 {
				trailHeads = append(trailHeads, Day10Position{row: i, col: j})
			}
		}
	}

	fmt.Printf("trailHeads size: %d\n", len(trailHeads))
	// for _, head := range trailHeads {
	// 	fmt.Printf("%d, %d\n", head.row, head.col)
	// }

	// Find all trails
	var trails []Day10Trail

	for _, head := range trailHeads {
		if day10TakeNextStep(head, gameMap, head, &trails) {
			fmt.Printf("trail found for trailHead: %d,%d\n", head.row, head.col)
		}
	}

	// Print trails
	// fmt.Printf("trails size: %d\n", len(trails))
	// for _, trail := range trails {
	// 	fmt.Printf("head %d,%d tail %d,%d\n", trail.headRow, trail.headCol, trail.tailRow, trail.tailCol)
	// }

	// Remove duplicate trails
	var destTrails []Day10Trail
	for _, trail := range trails {
		seen := false
		for _, destTrail := range destTrails {
			if destTrail == trail {
				seen = true
				break
			}
		}
		if !seen {
			destTrails = append(destTrails, trail)
		}
	}

	// Print destTrails
	// fmt.Printf("destTrails size: %d\n", len(destTrails))
	// for _, trail := range destTrails {
	// 	fmt.Printf("head %d,%d tail %d,%d\n", trail.headRow, trail.headCol, trail.tailRow, trail.tailCol)
	// }

	// Calculate and print scores
	totalScore := 0
	for _, head := range trailHeads {
		trailScore := 0
		for _, trail := range destTrails {
			if trail.headRow == head.row && trail.headCol == head.col {
				trailScore++
			}
		}
		//fmt.Printf("trailHead %d,%d score %d\n", head.row, head.col, trailScore)
		totalScore += trailScore
	}
	fmt.Printf("Part 1 totalScore: %d\n", totalScore)

	// Calculate Part 2 score
	totalScore = 0
	for _, head := range trailHeads {
		trailScore := 0
		for _, trail := range trails {
			if trail.headRow == head.row && trail.headCol == head.col {
				trailScore++
			}
		}
		//fmt.Printf("trailHead %d,%d score %d\n", head.row, head.col, trailScore)
		totalScore += trailScore
	}
	fmt.Printf("Part 2 totalScore: %d\n", totalScore)

	fmt.Println("2024 Day 10 end")
}

// day10TakeNextStep recursively explores possible paths from the current position.
// It returns true if a valid trail ending at elevation 9 is found.
// Parameters:
//   - currentPos: The current position being explored
//   - gameMap: The 2D grid containing elevation values
//   - startingPos: The original starting position (elevation 0)
//   - trails: Pointer to slice storing all valid trails found
func day10TakeNextStep(currentPos Day10Position, gameMap [][]int, startingPos Day10Position, trails *[]Day10Trail) bool {
	foundTrail := false
	currentElevation := gameMap[currentPos.row][currentPos.col]

	for _, step := range day10GetPossibleSteps(currentPos, gameMap) {
		possibleElevation := gameMap[step.row][step.col]
		if possibleElevation == currentElevation+1 {
			if possibleElevation == 9 {
				*trails = append(*trails, Day10Trail{
					headRow: startingPos.row,
					headCol: startingPos.col,
					tailRow: step.row,
					tailCol: step.col,
				})
				foundTrail = true
			} else {
				foundTrail = day10TakeNextStep(step, gameMap, startingPos, trails)
			}
		}
	}
	return foundTrail
}

// day10GetPossibleSteps returns all valid adjacent positions from the current position.
// A valid position must be within the game map boundaries.
// Parameters:
//   - currentPos: The position to check adjacent squares from
//   - gameMap: The 2D grid containing elevation values
//
// Returns:
//   - A slice of valid Day10Position that can be moved to
func day10GetPossibleSteps(currentPos Day10Position, gameMap [][]int) []Day10Position {
	possibleSteps := []Day10Position{
		{row: currentPos.row - 1, col: currentPos.col}, // up
		{row: currentPos.row + 1, col: currentPos.col}, // down
		{row: currentPos.row, col: currentPos.col - 1}, // left
		{row: currentPos.row, col: currentPos.col + 1}, // right
	}

	var validSteps []Day10Position
	rows, cols := len(gameMap), len(gameMap[0])

	for _, step := range possibleSteps {
		if step.row >= 0 && step.row < rows && step.col >= 0 && step.col < cols {
			validSteps = append(validSteps, step)
		}
	}

	return validSteps
}
