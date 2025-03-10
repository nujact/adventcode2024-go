// Package Day implements solutions for the Advent of Code 2024 challenges.
package Day

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// day13Machine represents a machine with two buttons and a prize location.
// Each button press moves the player in a specific X,Y direction.
// The goal is to reach the prize location with the minimum cost.
type day13Machine struct {
	buttonAX     int         // X movement when pressing button A
	buttonAY     int         // Y movement when pressing button A
	buttonBX     int         // X movement when pressing button B
	buttonBY     int         // Y movement when pressing button B
	prizeX       int         // X coordinate of the prize
	prizeY       int         // Y coordinate of the prize
	possibleRuns []*day13Run // All possible combinations of button presses that reach the prize
}

// day13Run represents a possible solution for reaching the prize.
// Each run tracks how many times each button was pressed and the total cost.
type day13Run struct {
	buttonAPresses int   // Number of times button A was pressed
	buttonBPresses int   // Number of times button B was pressed
	totalCost      int64 // Total cost (button A costs 3, button B costs 1)
}

// newMachine creates a new machine from a configuration string.
// The configuration string format is:
// Button A: X+n1, Y+n2
// Button B: X+n3, Y+n4
// Prize: X=n5, Y=n6
func newMachine(configStanza string) *day13Machine {
	m := &day13Machine{}
	lines := strings.Split(configStanza, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "Button A:") {
			parts := strings.Split(line, " ")
			m.buttonAX = parseInt(parts[2][2 : len(parts[2])-1])
			m.buttonAY = parseInt(parts[3][2:])
		} else if strings.HasPrefix(line, "Button B:") {
			parts := strings.Split(line, " ")
			m.buttonBX = parseInt(parts[2][2 : len(parts[2])-1])
			m.buttonBY = parseInt(parts[3][2:])
		} else if strings.HasPrefix(line, "Prize:") {
			parts := strings.Split(line, " ")
			m.prizeX = parseInt(parts[1][2 : len(parts[1])-1])
			m.prizeY = parseInt(parts[2][2:])
		}
	}
	m.possibleRuns = make([]*day13Run, 0)
	return m
}

// parseInt converts a string to an integer, ignoring errors.
// Used for parsing configuration values.
func parseInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

// Day13 solves the Day 13 puzzle of Advent of Code 2024.
// The puzzle involves finding the optimal way to reach a prize location
// by pressing two buttons (A and B) that move in different directions.
// Button A costs 3 units and Button B costs 1 unit.
// For each machine configuration, we need to:
// 1. Find all possible combinations of button presses that reach the prize
// 2. Find the combination with the lowest total cost
// 3. Sum up the lowest costs across all machines
func Day13() {
	fmt.Println("2024 Day 13 start")

	// Read input file
	file, err := os.Open("../inputs/input13.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var inputMemory strings.Builder
	for scanner.Scan() {
		inputMemory.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Test data
	inputMemory.Reset()
	inputMemory.WriteString("Button A: X+94, Y+34\n" +
		"Button B: X+22, Y+67\n" +
		"Prize: X=8400, Y=5400\n" +
		"\n" +
		"Button A: X+26, Y+66\n" +
		"Button B: X+67, Y+21\n" +
		"Prize: X=12748, Y=12176\n" +
		"\n" +
		"Button A: X+17, Y+86\n" +
		"Button B: X+84, Y+37\n" +
		"Prize: X=7870, Y=6450\n" +
		"\n" +
		"Button A: X+69, Y+23\n" +
		"Button B: X+27, Y+71\n" +
		"Prize: X=18641, Y=10279")

	fmt.Printf("inputMemory size: %d\n", inputMemory.Len())

	// Split input into stanzas, each representing a machine configuration
	inputStanzas := strings.Split(strings.TrimSpace(inputMemory.String()), "\n\n")

	// Print inputStanzas for verification
	fmt.Println("\ninputStanzas:")
	for _, stanza := range inputStanzas {
		fmt.Printf("%s\n\n", stanza)
	}

	// Parse inputStanzas into machines
	machines := make([]*day13Machine, 0)
	for _, stanza := range inputStanzas {
		machine := newMachine(stanza)
		machines = append(machines, machine)
	}

	// Print machines for verification
	fmt.Println("\nmachines:")
	for _, machine := range machines {
		fmt.Printf("Button A: %d, %d\n", machine.buttonAX, machine.buttonAY)
		fmt.Printf("Button B: %d, %d\n", machine.buttonBX, machine.buttonBY)
		fmt.Printf("Prize: %d, %d\n\n", machine.prizeX, machine.prizeY)
	}

	// Calculate possible runs for each machine
	// For each possible number of A button presses:
	// 1. Calculate required B button presses to reach prize X coordinate
	// 2. Verify if those button presses also reach prize Y coordinate
	// 3. If valid, calculate total cost and add to possible runs
	for _, machine := range machines {
		for buttonAPresses := 0; buttonAPresses < machine.prizeX/machine.buttonAX; buttonAPresses++ {
			xPos := buttonAPresses * machine.buttonAX
			buttonBPresses := (machine.prizeX - xPos) / machine.buttonBX

			// Verify this combination reaches both X and Y coordinates
			if (buttonAPresses*machine.buttonAX)+(buttonBPresses*machine.buttonBX) != machine.prizeX {
				continue
			}
			if (buttonAPresses*machine.buttonAY)+(buttonBPresses*machine.buttonBY) != machine.prizeY {
				continue
			}

			// Valid combination found, calculate cost and add to possible runs
			machine.possibleRuns = append(machine.possibleRuns, &day13Run{
				buttonAPresses: buttonAPresses,
				buttonBPresses: buttonBPresses,
				totalCost:      int64(buttonAPresses)*3 + int64(buttonBPresses),
			})
		}
	}

	// Print possible runs for each machine
	fmt.Println("\npossible runs:")
	for _, machine := range machines {
		fmt.Printf("Machine AX,Y BX,Y PX,Y: A %d,%d B %d,%d P %d,%d\n",
			machine.buttonAX, machine.buttonAY,
			machine.buttonBX, machine.buttonBY,
			machine.prizeX, machine.prizeY)

		for _, run := range machine.possibleRuns {
			fmt.Printf("Button A: %d, Button B: %d, Total Cost: %d\n",
				run.buttonAPresses, run.buttonBPresses, run.totalCost)
		}
	}

	// Find winning run (lowest cost) for each machine and calculate total
	machineWithWinner := 0
	var allCosts int64
	for _, machine := range machines {
		var winner *day13Run
		for _, run := range machine.possibleRuns {
			if winner == nil || run.totalCost < winner.totalCost {
				winner = run
			}
		}
		if winner != nil {
			machineWithWinner++
			allCosts += winner.totalCost
		}
	}
	fmt.Printf("\nMachine with winner: %d, Total Cost: %d\n", machineWithWinner, allCosts)

	fmt.Println("2024 Day 13 end")
}
