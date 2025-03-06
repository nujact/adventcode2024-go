package Day

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Equation represents a mathematical equation with target result and input values
// The puzzle involves finding ways to combine input values using operators to reach the target result
type Equation struct {
	TargetResult    int64   // The target result to achieve
	InputValues     []int   // List of input values to use in calculations
	PossibleResults []int64 // List of all possible results from different operator combinations
}

// NewEquation creates a new Equation from an input string
// Format: "target: value1 value2 value3 ..."
// Example: "190: 10 19" means target is 190, input values are [10, 19]
func NewEquation(inputString string) *Equation {
	splitString := strings.Split(inputString, ": ")
	targetResult, _ := strconv.ParseInt(splitString[0], 10, 64)

	// Parse input values from space-separated string
	inputValues := make([]int, 0)
	for _, input := range strings.Split(splitString[1], " ") {
		value, _ := strconv.Atoi(input)
		inputValues = append(inputValues, value)
	}

	return &Equation{
		TargetResult:    targetResult,
		InputValues:     inputValues,
		PossibleResults: make([]int64, 0),
	}
}

// calcNextInput recursively calculates all possible results using different operators
// Parameters:
//   - inputList: List of input values to combine
//   - index: Current position in inputList
//   - result: Current accumulated result
//   - operator: Current operator to apply (+, *, |)
//   - allResults: Pointer to slice storing all possible results
//   - isPart2: Whether to include concatenation operator (|) for part 2
//
// Operators:
//   - + : Addition (e.g., 10 + 19 = 29)
//   - * : Multiplication (e.g., 10 * 19 = 190)
//   - | : Concatenation (e.g., 10 | 19 = 1019) - only in part 2
func calcNextInput(inputList []int, index int, result int64, operator string, allResults *[]int64, isPart2 bool) {
	value := inputList[index]

	// Apply operator to current result and value
	switch operator {
	case "+":
		result += int64(value)
	case "*":
		result *= int64(value)
	case "|":
		// Concatenate numbers as strings and convert back to int64
		res := fmt.Sprintf("%d%d", result, value)
		result, _ = strconv.ParseInt(res, 10, 64)
	}

	// If we've used all input values, add result to allResults
	if index == len(inputList)-1 {
		*allResults = append(*allResults, result)
	} else {
		// Try all operators for next value
		calcNextInput(inputList, index+1, result, "+", allResults, isPart2)
		calcNextInput(inputList, index+1, result, "*", allResults, isPart2)
		if isPart2 {
			calcNextInput(inputList, index+1, result, "|", allResults, isPart2)
		}
	}
}

// Day7 solves the Advent of Code 2024 Day 7 puzzle
// The puzzle involves finding equations where the target result can be achieved
// using different combinations of operators on the input values
//
// Part 1: Use addition (+) and multiplication (*) operators
// Part 2: Also use concatenation (|) operator
func Day7() {
	fmt.Println("2024 Day7 start")

	// Open file and gather raw inputs
	inputMemory, shouldReturn := day7GetInput()
	if shouldReturn {
		return
	}

	fmt.Println("inputMemory size:", len(inputMemory))

	// Process parts
	day7part1(inputMemory, false) // Part 1: Only + and * operators
	day7part2(inputMemory)        // Part 2: Also includes | operator

	fmt.Println("2024 Day7 end")
}

// day7part1 processes part 1 of the puzzle
// Finds equations where the target result can be achieved using the input values
// Parameters:
//   - inputMemory: Raw input string containing equations
//   - isPart2: Whether to include concatenation operator (|)
func day7part1(inputMemory string, isPart2 bool) {
	// Parse input into equations
	equations := make([]*Equation, 0)
	inputStrings := strings.Split(inputMemory, "\n")
	for _, inputString := range inputStrings {
		equation := NewEquation(inputString)
		equations = append(equations, equation)
	}

	// Print equations for debugging
	// for _, equation := range equations {
	// 	fmt.Printf("Equation %d %v\n", equation.TargetResult, equation.InputValues)
	// }
	// fmt.Println()

	// Calculate possible results for each equation
	for _, equation := range equations {
		allResults := make([]int64, 0)
		// Try each operator starting with first value
		calcNextInput(equation.InputValues, 1, int64(equation.InputValues[0]), "+", &allResults, isPart2)
		calcNextInput(equation.InputValues, 1, int64(equation.InputValues[0]), "*", &allResults, isPart2)
		if isPart2 {
			calcNextInput(equation.InputValues, 1, int64(equation.InputValues[0]), "|", &allResults, isPart2)
		}
		equation.PossibleResults = allResults
	}

	// Find equations with target result
	var sumSuccessTargets int64
	for _, equation := range equations {
		// Check if target result is in possible results
		for _, result := range equation.PossibleResults {
			if result == equation.TargetResult {
				//fmt.Printf("success Equation %d %v %v\n", equation.TargetResult, equation.InputValues, equation.PossibleResults)
				sumSuccessTargets += equation.TargetResult
				break
			}
		}
	}
	if isPart2 {
		fmt.Printf("\nsumSuccessTargets part2: %d\n\n", sumSuccessTargets)
	} else {
		fmt.Printf("\nsumSuccessTargets part1: %d\n\n", sumSuccessTargets)
	}
}

// day7part2 processes part 2 of the puzzle
// Uses the same logic as part 1 but includes the concatenation operator (|)
func day7part2(inputMemory string) {
	// call part1 with isPart2 = true
	day7part1(inputMemory, true)
}

// day7GetInput reads the input file and returns its contents as a string
// Returns the input string and a boolean indicating if there was an error
//
// The input file should contain equations in the format:
// target: value1 value2 value3 ...
// Example:
// 190: 10 19
// 3267: 81 40 27
func day7GetInput() (string, bool) {
	var inputFileNameBegins = "input" // "input" or "test"
	file, err := os.Open(fmt.Sprintf("../inputs/%s7.txt", inputFileNameBegins))
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
