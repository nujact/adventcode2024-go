package Day

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Day1 solves the Advent of Code 2024 Day 1 puzzle
// The puzzle involves processing two lists of numbers and finding relationships between them
// Part A: For each pair of numbers (one from each list), calculate the absolute difference
//
//	and sum all differences
//
// Part B: For each number in the left list, count how many times it appears in the right list
//
//	and add the product of the number and its count to the total
func Day1() {
	fmt.Println("Day 1")

	// Solve part A and print result
	solutionA := solutionA()
	fmt.Println("Solution A:", solutionA)

	// Solve part B and print result
	solutionB := solutionB()
	fmt.Println("Solution B:", solutionB)
}

// solutionA solves part A of the puzzle
// For each pair of numbers (one from each list), calculates the absolute difference
// and sums all differences
// Returns:
//   - The sum of all absolute differences between paired numbers
//   - -1 if there was an error reading the input
func solutionA() int {
	// Get input data
	leftList, rightList, shouldReturn := GetInputs()
	if shouldReturn {
		return -1
	}

	// Sort both lists to ensure proper pairing
	sort.Ints(leftList)
	sort.Ints(rightList)

	// Calculate sum of absolute differences between paired numbers
	var solution int = 0
	for i, leftItem := range leftList {
		rightItem := rightList[i]
		dist := int(math.Abs(float64(leftItem) - float64(rightItem)))
		solution += dist
	}

	return solution
}

// GetInputs reads the input file and parses it into two lists of integers
// The input file contains pairs of numbers separated by three spaces
// Returns:
//   - left: List of numbers from the left column
//   - right: List of numbers from the right column
//   - shouldReturn: true if there was an error reading the file
func GetInputs() ([]int, []int, bool) {
	var left []int
	var right []int

	// Open file and gather raw inputs
	var inputFileNameBegins = "input" // "input" or "test"
	file, err := os.Open(fmt.Sprintf("inputs/%s1.txt", inputFileNameBegins))
	if err != nil {
		fmt.Println(err)
		return nil, nil, true
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
		return nil, nil, true
	}

	// split inputMemory into lines
	lines := strings.Split(inputMemory, "\n")

	// Parse each line into two numbers
	for _, line := range lines {
		slice := strings.Split(line, "   ") // Split by three spaces

		// Parse left number
		t, err := strconv.Atoi(slice[0])
		if err == nil {
			left = append(left, t)
		}

		// Parse right number
		t, err = strconv.Atoi(slice[1])
		if err == nil {
			right = append(right, t)
		}
	}
	return left, right, false
}

// solutionB solves part B of the puzzle
// For each number in the left list, counts how many times it appears in the right list
// and adds the product of the number and its count to the total
// Returns:
//   - The sum of all products (number × count)
//   - -1 if there was an error reading the input
func solutionB() int {
	// Get input data
	leftList, rightList, shouldReturn := GetInputs()
	if shouldReturn {
		return -1
	}

	// Calculate sum of products
	var solution int = 0
	for _, leftItem := range leftList {
		// Count occurrences of leftItem in rightList
		count := 0
		for _, rightItem := range rightList {
			if leftItem == rightItem {
				count++
			}
		}

		// Add product to total
		solution += leftItem * count
	}

	return solution
}
