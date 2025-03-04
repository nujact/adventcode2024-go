package Day

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

// Day3 is the entry point for Day 3 solution
// It reads input, processes it, and outputs results
func Day3() {
	fmt.Println("2024 Day 3 start")

	inputStr, done := GetInput()
	if done {
		return
	}

	fmt.Printf("inputMemory size: %d\n", len(inputStr))
	fmt.Printf("inputMemory: %s\n", inputStr)

	// Uncomment to run part 1
	totalPt1 := getTotalPt1(inputStr)
	fmt.Printf("total Pt1: %d\n", totalPt1)

	// Run part 2
	totalPt2 := getTotalPt2(inputStr)
	fmt.Printf("total Pt2: %d\n", totalPt2)

	fmt.Println("2024 Day 3 end")
}

func GetInput() (string, bool) {
	// Read the entire input file into memory
	var inputFileNameBegins = "test" // "input" or "test"
	file, err := os.Open(fmt.Sprintf("inputs/%s3.txt", inputFileNameBegins))
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return "", true
	}
	defer file.Close()

	var inputStr = ""

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		inputStr += line
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return "", true
	}
	return inputStr, false
}

// getTotalPt2 processes input by matching patterns and calculating the sum
// It looks for mul(num1,num2) patterns and adds their products to a total
// The doFlag state can be toggled by do() and don't() tokens to control whether calculations are performed
func getTotalPt2(inputStr string) int64 {
	regex := regexp.MustCompile(`do\(\)|don't\(\)|mul\((\d+),(\d+)\)`)
	matches := regex.FindAllStringSubmatch(inputStr, -1)

	var total int64 = 0
	doFlag := true // true = do, false = don't

	for _, match := range matches {
		fmt.Printf("match: %s\n", match[0])

		// Handle control flags
		if match[0] == "do()" {
			doFlag = true
			continue
		} else if match[0] == "don't()" {
			doFlag = false
			continue
		}

		// Skip multiplication if doFlag is false
		if !doFlag {
			continue
		}

		// Process multiplication
		// In Go regex, the whole match is at index 0, and capture groups start at index 1
		num1, err1 := strconv.ParseInt(match[1], 10, 64)
		num2, err2 := strconv.ParseInt(match[2], 10, 64)

		if err1 == nil && err2 == nil {
			total += (num1 * num2)
			// Uncomment for debugging
			// fmt.Printf("num1: %d num2: %d total: %d\n", num1, num2, total)
		} else {
			if err1 != nil {
				fmt.Printf("Error parsing num1: %v\n", err1)
			}
			if err2 != nil {
				fmt.Printf("Error parsing num2: %v\n", err2)
			}
		}
	}

	return total
}

// getTotalPt1 processes input by matching mul(num1,num2) patterns
// and calculating the sum of all products
func getTotalPt1(inputStr string) int64 {
	regex := regexp.MustCompile(`mul\((\d+),(\d+)\)`)
	matches := regex.FindAllStringSubmatch(inputStr, -1)

	var total int64 = 0

	for _, match := range matches {
		fmt.Printf("match: %s\n", match[0])

		// Process multiplication
		num1, err1 := strconv.ParseInt(match[1], 10, 64)
		num2, err2 := strconv.ParseInt(match[2], 10, 64)

		if err1 == nil && err2 == nil {
			total += (num1 * num2)
			// Uncomment for debugging
			// fmt.Printf("num1: %d num2: %d total: %d\n", num1, num2, total)
		} else {
			if err1 != nil {
				fmt.Printf("Error parsing num1: %v\n", err1)
			}
			if err2 != nil {
				fmt.Printf("Error parsing num2: %v\n", err2)
			}
		}
	}

	return total
}
