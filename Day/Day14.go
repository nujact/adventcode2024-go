// Package Day implements solutions for the Advent of Code 2024 challenges.
package Day

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// day14Robot represents a robot with position and velocity.
// Each robot moves in a fixed direction and wraps around the room boundaries.
// Multiple robots can occupy the same position.
type day14Robot struct {
	px, py int // Current position (x,y)
	vx, vy int // Velocity vector (x,y)
}

// newRobot creates a new robot with the specified position and velocity.
// Parameters:
//
//	px, py: Initial position coordinates
//	vx, vy: Velocity vector components
func newRobot(px, py, vx, vy int) *day14Robot {
	return &day14Robot{
		px: px,
		py: py,
		vx: vx,
		vy: vy,
	}
}

// printRoom prints the current state of the room.
// The room is displayed as a grid where:
// - '.' represents an empty cell
// - Numbers represent how many robots are in that cell
func printRoom(roomHeight, roomWidth int, room [][]int) {
	fmt.Println("Room:")
	for j := 0; j < roomHeight; j++ {
		for i := 0; i < roomWidth; i++ {
			if room[i][j] == 0 {
				fmt.Print(". ")
			} else {
				fmt.Printf("%d ", room[i][j])
			}
		}
		fmt.Println()
	}
}

// Day14 solves the Day 14 puzzle of Advent of Code 2024.
// The puzzle involves simulating robots moving in a room:
// 1. Each robot has a fixed velocity and wraps around room boundaries
// 2. Multiple robots can occupy the same position
// 3. After simulation, the room is divided into quadrants
// 4. The answer is the product of robot counts in each quadrant
func Day14() {
	fmt.Println("2024 Day 14 start")

	// Read input file
	file, err := os.Open("../inputs/input14.txt")
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

	// Test data - format: "p=x,y v=vx,vy" for each robot
	inputMemory.Reset()
	inputMemory.WriteString("p=0,4 v=3,-3\n" +
		"p=6,3 v=-1,-3\n" +
		"p=10,3 v=-1,2\n" +
		"p=2,0 v=2,-1\n" +
		"p=0,0 v=1,3\n" +
		"p=3,0 v=-2,-2\n" +
		"p=7,6 v=-1,-3\n" +
		"p=3,0 v=-1,-2\n" +
		"p=9,3 v=2,3\n" +
		"p=7,3 v=-1,2\n" +
		"p=2,4 v=2,-3\n" +
		"p=9,5 v=-3,-3")

	fmt.Printf("inputMemory size: %d\n", inputMemory.Len())

	// Parse input into robots
	// Format: "p=x,y v=vx,vy" where:
	// - x,y is the initial position
	// - vx,vy is the velocity vector
	robots := make([]*day14Robot, 0)
	for _, line := range strings.Split(strings.TrimSpace(inputMemory.String()), "\n") {
		parts := strings.Split(line, " ")
		positionParts := strings.Split(parts[0], ",")
		velocityParts := strings.Split(parts[1], ",")

		px := day14ParseInt(positionParts[0][2:])
		py := day14ParseInt(positionParts[1])
		vx := day14ParseInt(velocityParts[0][2:])
		vy := day14ParseInt(velocityParts[1])

		robots = append(robots, newRobot(px, py, vx, vy))
	}

	// Print initial robot positions and velocities
	fmt.Println("Robots:")
	for _, robot := range robots {
		fmt.Printf("P %d,%d  V %d,%d\n", robot.px, robot.py, robot.vx, robot.vy)
	}

	// Define room dimensions
	roomWidth := 11
	roomHeight := 7

	// Create empty room grid
	rooms := make([][]int, roomWidth)
	for i := range rooms {
		rooms[i] = make([]int, roomHeight)
	}

	// Mark initial robot positions in room
	for _, robot := range robots {
		rooms[robot.px][robot.py]++
	}

	// Print initial room state
	printRoom(roomHeight, roomWidth, rooms)

	// Simulate robot movement for specified number of steps
	// Each step:
	// 1. Update position based on velocity
	// 2. Handle wrapping around room boundaries
	steps := 100
	for _, robot := range robots {
		for i := 1; i <= steps; i++ {
			robot.px += robot.vx
			robot.py += robot.vy

			// Handle wrapping around room boundaries
			if robot.px < 0 {
				robot.px = roomWidth + robot.px
			}
			if robot.px >= roomWidth {
				robot.px = robot.px - roomWidth
			}
			if robot.py < 0 {
				robot.py = roomHeight + robot.py
			}
			if robot.py >= roomHeight {
				robot.py = robot.py - roomHeight
			}
		}
	}

	// Reset room and update with final robot positions
	for i := range rooms {
		for j := range rooms[i] {
			rooms[i][j] = 0
		}
	}
	for _, robot := range robots {
		rooms[robot.px][robot.py]++
	}

	// Calculate robot count in each quadrant
	// Room is divided into four quadrants by the center point
	quadrantRobotCount := make([]int, 4)
	for i := 0; i < roomWidth; i++ {
		for j := 0; j < roomHeight; j++ {
			if rooms[i][j] > 0 {
				if i < roomWidth/2 && j < roomHeight/2 {
					quadrantRobotCount[0] += rooms[i][j] // Top-left quadrant
				} else if i > roomWidth/2 && j < roomHeight/2 {
					quadrantRobotCount[1] += rooms[i][j] // Top-right quadrant
				} else if i < roomWidth/2 && j > roomHeight/2 {
					quadrantRobotCount[2] += rooms[i][j] // Bottom-left quadrant
				} else if i > roomWidth/2 && j > roomHeight/2 {
					quadrantRobotCount[3] += rooms[i][j] // Bottom-right quadrant
				}
			}
		}
	}

	// Calculate final answer as product of quadrant counts
	answer := 1
	for _, count := range quadrantRobotCount {
		answer *= count
	}

	// Print final room state and answer
	printRoom(roomHeight, roomWidth, rooms)
	fmt.Printf("Answer: %d\n", answer)
	fmt.Println("2024 Day 14 end")
}

// day14ParseInt converts a string to an integer, ignoring errors.
// Used for parsing robot position and velocity values.
func day14ParseInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
