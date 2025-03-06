package Day

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Day5 solves the Advent of Code 2024 Day 5 puzzle
// The puzzle involves processing page order rules and updates
// Part 1: Find the sum of middle pages from updates that are already valid
// Part 2: Find the sum of middle pages from updates that can be made valid by reordering
func Day5() {
	fmt.Println("2024 Day5 start")

	// Open file and gather raw inputs
	inputMemory, shouldReturn := day5GetInput()
	if shouldReturn {
		return
	}

	fmt.Println("inputMemory size:", len(inputMemory))

	// Parse input into rules and updates
	inputRules := getPageOrderRules(inputMemory)
	inputUpdates := getUpdates(inputMemory)

	// // Print page order rules
	// fmt.Printf("inputRules: %d\n", len(inputRules))
	// for _, rule := range inputRules {
	// 	fmt.Println(rule)
	// }

	// // Print updates
	// fmt.Printf("\ninputUpdates: %d\n", len(inputUpdates))
	// for _, update := range inputUpdates {
	// 	fmt.Println(update)
	// }

	// Process parts
	day5part1(inputUpdates, inputRules)
	day5part2(inputUpdates, inputRules)

	fmt.Println("2024 Day5 end")
}

// day5part1 processes part 1 of the puzzle
// Finds the sum of middle pages from updates that are already valid
// An update is valid if all its pages are in the correct order according to the rules
func day5part1(inputUpdates [][]string, inputRules []string) {
	allValid := true
	middleOfTruth := 0

	for _, update := range inputUpdates {
		allValid = true
		for i := 0; i < len(update); i++ {
			// Validate each page of this update
			if !validateOrder(inputRules, i, update) {
				allValid = false
			}
		}

		if allValid {
			middlePage := update[getMiddlePage(update)]
			val, _ := strconv.Atoi(middlePage)
			middleOfTruth += val
		}

		// fmt.Printf("# %v %v\n", update, allValid)
	}

	fmt.Printf("middleSum %d\n", middleOfTruth)
}

// day5part2 processes part 2 of the puzzle
// Finds the sum of middle pages from updates that can be made valid by reordering
// An update can be made valid by moving pages to positions that satisfy the rules
func day5part2(inputUpdates [][]string, inputRules []string) {
	// Build list of invalid updates
	invalidUpdates := make([]int, 0)
	for j, update := range inputUpdates {
		allValid := true
		for i := 0; i < len(update); i++ {
			if !validateOrder(inputRules, i, update) {
				allValid = false
			}
		}
		if !allValid {
			invalidUpdates = append(invalidUpdates, j)
		}
	}

	// Try to make invalid updates valid by reordering
	allValid := false
	for !allValid {
		allValid = true
		for i := 0; i < len(invalidUpdates); i++ {
			update := inputUpdates[invalidUpdates[i]]
			// Loop through pages of this update
			for j := 0; j < len(update); j++ {
				if !validateOrder(inputRules, j, update) {
					allValid = false
					update = moveOrder(inputRules, j, update)
					inputUpdates[invalidUpdates[i]] = update
				}
			}
		}
	}

	// Calculate sum of middle pages from invalid updates
	middleOfTruth := 0
	for _, invalidUpdate := range invalidUpdates {
		update := inputUpdates[invalidUpdate]
		// fmt.Printf("update %v", update)
		middlePage := update[getMiddlePage(update)]
		// fmt.Printf("  mid %s\n", middlePage)
		val, _ := strconv.Atoi(middlePage)
		middleOfTruth += val
	}

	fmt.Printf("middleSum %d\n", middleOfTruth)
}

// getMiddlePage returns the index of the middle page in an update
// For even-length updates, returns the index of the first middle page
func getMiddlePage(update []string) int {
	return (len(update) / 2)
}

// moveOrder attempts to move a page to a valid position based on rules
// If a page should come before another page according to the rules,
// it moves the page to a position before that other page
func moveOrder(inputRules []string, thisIndex int, update []string) []string {
	// Find pages that should occur before updatePage
	mustFollow := make([]string, 0)
	updatePage := update[thisIndex]

	for _, rule := range inputRules {
		if strings.Contains(rule, "|"+updatePage) {
			parts := strings.Split(rule, "|")
			mustFollow = append(mustFollow, parts[0])
		}
	}

	// Check if current entry has pages AFTER this that are in the mustFollow list
	validOrder := true
	for i := thisIndex + 1; i < len(update); i++ {
		pg := update[i]
		for _, mustFollowPg := range mustFollow {
			if mustFollowPg == pg {
				validOrder = false
				update = listAddAt(update, i, updatePage)
				update = listRemAt(update, thisIndex, updatePage)
				break
			}
		}
		if !validOrder {
			break
		}
	}

	return update
}

// listRemAt removes an element at a specific position from a string slice
// Returns a new slice with the element removed
func listRemAt(orgList []string, posRem int, remElement string) []string {
	newList := make([]string, 0)
	for i, pg := range orgList {
		if i == posRem && pg == remElement {
			continue
		}
		newList = append(newList, pg)
	}
	return newList
}

// listAddAt adds an element at a specific position in a string slice
// Returns a new slice with the element added after the specified position
func listAddAt(orgList []string, posAddAfter int, newElement string) []string {
	newList := make([]string, 0)
	for i, pg := range orgList {
		newList = append(newList, pg)
		if i == posAddAfter {
			newList = append(newList, newElement)
		}
	}
	return newList
}

// validateOrder checks if a page is in a valid position according to the rules
// A page is valid if all pages that should come before it according to the rules
// are actually before it in the update
func validateOrder(inputRules []string, thisIndex int, update []string) bool {
	// Find pages that should occur before updatePage
	mustFollow := make([]string, 0)
	updatePage := update[thisIndex]

	for _, rule := range inputRules {
		if strings.Contains(rule, "|"+updatePage) {
			parts := strings.Split(rule, "|")
			mustFollow = append(mustFollow, parts[0])
		}
	}

	// Check if current entry has pages AFTER this that are in the mustFollow list
	validOrder := true
	for i := thisIndex + 1; i < len(update); i++ {
		if contains(mustFollow, update[i]) {
			validOrder = false
		}
	}

	return validOrder
}

// getPageOrderRules extracts the page order rules from the input
// Rules are in the format "page1|page2" where page1 must come before page2
// Returns a slice of rules, stopping at the first empty line
func getPageOrderRules(inputMemory string) []string {
	// Split by new line into array
	inputRows := strings.Split(inputMemory, "\n")
	pageOrderRulesList := make([]string, 0)
	for _, row := range inputRows {
		// Add rows till empty line
		if row == "" {
			break
		}
		pageOrderRulesList = append(pageOrderRulesList, row)
	}
	return pageOrderRulesList
}

// getUpdates extracts the updates from the input
// Updates are comma-separated numbers representing page numbers
// Returns a slice of updates, starting after the first empty line
func getUpdates(inputMemory string) [][]string {
	// Split by new line into array
	inputRows := strings.Split(inputMemory, "\n")
	updateList := make([][]string, 0)
	seenEmptyLine := false
	for _, row := range inputRows {
		// Skip rows till empty line
		if !seenEmptyLine {
			if row == "" {
				seenEmptyLine = true
			}
			continue
		}
		update := strings.Split(row, ",")
		updateList = append(updateList, update)
	}
	return updateList
}

// contains checks if a string slice contains a specific string
// Returns true if the string is found in the slice
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// day5GetInput reads the input file and returns its contents as a string
// Returns the input string and a boolean indicating if there was an error
func day5GetInput() (string, bool) {
	var inputFileNameBegins = "input" // "input" or "test"
	file, err := os.Open(fmt.Sprintf("../inputs/%s5.txt", inputFileNameBegins))
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
