package Day

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// Day2 solves the Advent of Code 2024 Day 2 puzzle
// The puzzle involves analyzing sequences of numbers to determine if they are "safe"
// Part 1: A sequence is safe if it's either strictly increasing or decreasing AND
//
//	the difference between consecutive numbers is ≤ 3
//
// Part 2: Similar to part 1 but allows one number to be removed to make the sequence safe
func Day2() {
	fmt.Println("2024 Day2 start")

	// Open file and gather raw inputs
	inputArray := [][]int64{}
	var inputFileNameBegins = "input" // "input" or "test"
	file, err := os.Open(fmt.Sprintf("inputs/%s2.txt", inputFileNameBegins))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Read input file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		gatherInputs2(line, &inputArray)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("inputArray size:", len(inputArray))

	// Process the inputs
	safeCount := 0
	safeCount2 := 0
	for _, row := range inputArray {
		// Uncomment to use part 1
		// safeCount += processInputs2(row)

		// For part 2, try both with and without removing a number
		row2 := make([]int64, len(row))
		copy(row2, row)
		safe2 := processInputs2pt2(row, true)
		if safe2 == 0 && len(row) < len(row2) {
			safe2 = processInputs2pt2(row, false)
		}
		safeCount2 += safe2
	}

	fmt.Println("safeCount:", safeCount)
	fmt.Println("safeCount2:", safeCount2)

	fmt.Println("2024 Day2 end")
}

// processInputs2 checks if a sequence of numbers is "safe" according to part 1 rules
// A sequence is safe if:
// 1. It's either strictly increasing or decreasing
// 2. The difference between consecutive numbers is ≤ 3
// Returns:
//   - 1 if the sequence is safe
//   - 0 if the sequence is not safe
func processInputs2(rowList []int64) int {
	isSafeRisingOrFalling := 0

	// First check if the sequence is strictly increasing
	for i := 0; i < len(rowList)-1; i++ {
		if rowList[i] < rowList[i+1] {
			isSafeRisingOrFalling = 1
		} else {
			isSafeRisingOrFalling = 0
			break
		}
	}

	// If not increasing, check if it's strictly decreasing
	if isSafeRisingOrFalling == 0 {
		for i := 0; i < len(rowList)-1; i++ {
			if rowList[i] > rowList[i+1] {
				isSafeRisingOrFalling = 1
			} else {
				isSafeRisingOrFalling = 0
				break
			}
		}
	}

	// Check if differences between consecutive numbers are ≤ 3
	isSafeGradual := 0
	if isSafeRisingOrFalling == 1 {
		for i := 0; i < len(rowList)-1; i++ {
			diff := int64(math.Abs(float64(rowList[i+1] - rowList[i])))
			if diff <= 3 {
				isSafeGradual = 1
			} else {
				isSafeGradual = 0
				break
			}
		}
	}

	return isSafeRisingOrFalling * isSafeGradual
}

// processInputs2pt2 checks if a sequence of numbers is "safe" according to part 2 rules
// Similar to part 1 but allows one number to be removed to make the sequence safe
// Parameters:
//   - rowList: The sequence of numbers to check
//   - allowForgive: Whether to allow removing one number to make the sequence safe
//
// Returns:
//   - 1 if the sequence is safe (or can be made safe by removing one number)
//   - 0 if the sequence is not safe
func processInputs2pt2(rowList []int64, allowForgive bool) int {
	forgiveRowNum := -1
	isSafeRisingOrFalling := 0

	// Check if sequence is strictly increasing
	for i := 0; i < len(rowList)-1; i++ {
		if rowList[i] < rowList[i+1] {
			isSafeRisingOrFalling = 1
		} else {
			if allowForgive {
				forgiveRowNum = i + 1
				allowForgive = false
				isSafeRisingOrFalling = 1
			} else {
				isSafeRisingOrFalling = 0
				break
			}
		}
	}

	// If not increasing, check if it's strictly decreasing
	if isSafeRisingOrFalling == 0 {
		for i := 0; i < len(rowList)-1; i++ {
			if i == forgiveRowNum {
				continue
			}
			if rowList[i] > rowList[i+1] {
				isSafeRisingOrFalling = 1
			} else {
				if allowForgive {
					forgiveRowNum = i + 1
					allowForgive = false
					isSafeRisingOrFalling = 1
				} else {
					isSafeRisingOrFalling = 0
					break
				}
			}
		}
	}

	// Check if differences between consecutive numbers are ≤ 3
	isSafeGradual := 0
	if isSafeRisingOrFalling == 1 {
		for i := 0; i < len(rowList)-1; i++ {
			if i == forgiveRowNum {
				continue
			}
			diff := int64(math.Abs(float64(rowList[i+1] - rowList[i])))
			if diff <= 3 {
				isSafeGradual = 1
			} else {
				if allowForgive {
					forgiveRowNum = i + 1
					isSafeGradual = 1
				} else {
					isSafeGradual = 0
					break
				}
			}
		}
	}

	// Remove the forgiven number if one was found
	if forgiveRowNum != -1 {
		rowList = append(rowList[:forgiveRowNum], rowList[forgiveRowNum+1:]...)
	}

	// Debug output for unsafe sequences
	if isSafeRisingOrFalling*isSafeGradual == 0 {
		fmt.Printf("rowList: %v isSafeRisingOrFalling: %d isSafeGradual: %d forgiveRowNum: %d final: %d\n",
			rowList, isSafeRisingOrFalling, isSafeGradual, forgiveRowNum, isSafeRisingOrFalling*isSafeGradual)
	}

	return isSafeRisingOrFalling * isSafeGradual
}

// gatherInputs2 parses a line of input into a slice of integers
// Parameters:
//   - line: The input line to parse
//   - inputArray: Pointer to the array where the parsed numbers will be stored
func gatherInputs2(line string, inputArray *[][]int64) {
	// Split the line by space
	parts := strings.Fields(line)
	rowList := []int64{}
	for _, part := range parts {
		num, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			fmt.Println("Error parsing:", err)
			continue
		}
		rowList = append(rowList, num)
	}
	*inputArray = append(*inputArray, rowList)
}
