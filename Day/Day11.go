// Package Day implements solutions for the Advent of Code 2024 challenges.
package Day

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// blinkCache stores the cached results of blink operations
type blinkCache struct {
	stone      int64
	blinkCount int
}

// Global cache map to store blink results
var cachedBlinks = make(map[blinkCache]int64)

// Day11 solves the Day 11 puzzle of Advent of Code 2024.
// The puzzle involves "blinking" stones according to specific rules:
// - Rule 1: flip 0 to 1
// - Rule 2: even length numbers are split into halves
// - Rule 3: odd length numbers are multiplied by 2024
func Day11() {
	fmt.Println("2024 Day 11 start")

	// Read input file
	file, err := os.Open("./inputs/input11.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var inputMemory strings.Builder
	for scanner.Scan() {
		inputMemory.WriteString(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Printf("inputMemory size: %d\n", inputMemory.Len())

	// Parse input string into array of int64
	var stones []int64
	for _, numStr := range strings.Split(inputMemory.String(), " ") {
		num, err := strconv.ParseInt(numStr, 10, 64)
		if err != nil {
			fmt.Printf("Error parsing number %s: %v\n", numStr, err)
			continue
		}
		stones = append(stones, num)
	}

	// Print starting stones
	fmt.Print("starting stones: ")
	for _, stone := range stones {
		fmt.Printf("%d ", stone)
	}
	fmt.Println()

	blinkCount := 75
	var totalStoneCount int64 = 0

	fmt.Printf("Blinking %d times, final row\n", blinkCount)
	for _, stone := range stones {
		blinkRecurseCount := blinkRecurse(stone, blinkCount)
		totalStoneCount += blinkRecurseCount
		fmt.Printf("blinkRecurseCount: %d stoneCount: %d\n", blinkRecurseCount, totalStoneCount)
	}

	fmt.Printf("\nstoneCount: %d\n", totalStoneCount)
	fmt.Println("2024 Day 11 end")
}

// blinkRecurse implements the recursive blinking logic
func blinkRecurse(stone int64, blinkCount int) int64 {
	if blinkCount <= 0 {
		return 1
	}

	// Check cache
	cacheKey := blinkCache{stone: stone, blinkCount: blinkCount}
	if cachedResult, exists := cachedBlinks[cacheKey]; exists {
		return cachedResult
	}

	var stoneCount int64

	// Implement the three rules
	if stone == 0 {
		// Rule 1: flip 0 to 1
		stoneCount = blinkRecurse(1, blinkCount-1)
	} else {
		// Convert stone to string to check length
		stoneStr := strconv.FormatInt(stone, 10)
		if len(stoneStr)%2 == 0 {
			// Rule 2: split even length numbers
			mid := len(stoneStr) / 2
			leftStone, _ := strconv.ParseInt(stoneStr[:mid], 10, 64)
			rightStone, _ := strconv.ParseInt(stoneStr[mid:], 10, 64)
			stoneCount = blinkRecurse(leftStone, blinkCount-1) + blinkRecurse(rightStone, blinkCount-1)
		} else {
			// Rule 3: multiply odd length numbers by 2024
			stoneCount = blinkRecurse(stone*2024, blinkCount-1)
		}
	}

	// Update cache
	cachedBlinks[cacheKey] = stoneCount
	return stoneCount
}
