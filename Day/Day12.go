// Package Day implements solutions for the Advent of Code 2024 challenges.
package Day

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Plot represents a single plot in the garden
type day12Plot struct {
	plant    string
	x, y     int
	regionID int
}

// Region represents a connected group of plots with the same plant type
type day12Region struct {
	id        int
	plots     []*day12Plot
	perimeter int
}

// Day12 solves the Day 12 puzzle of Advent of Code 2024
func Day12() {
	fmt.Println("2024 Day 12 start")

	// Read input file
	file, err := os.Open("../inputs/input12.txt")
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

	fmt.Printf("inputMemory size: %d\n", inputMemory.Len())

	// Parse input into plots
	lines := strings.Split(strings.TrimSpace(inputMemory.String()), "\n")
	plots := make([][]*day12Plot, len(lines))
	for i, line := range lines {
		plots[i] = make([]*day12Plot, len(line))
		for j, char := range line {
			plots[i][j] = &day12Plot{
				plant:    string(char),
				x:        i,
				y:        j,
				regionID: -1,
			}
		}
	}

	// Print plot
	fmt.Printf("plot: %dx%d\n", len(plots), len(plots[0]))
	for _, row := range plots {
		for _, cell := range row {
			fmt.Print(cell.plant)
		}
		fmt.Println()
	}

	// Get unique plants
	plants := getPlants(plots)
	fmt.Printf("plants: %d\n", len(plants))
	for _, plant := range plants {
		fmt.Println(plant)
	}

	// Get regions from plots
	regions := getRegionsFromPlots(plots)

	// Calculate area and perimeter for each region
	var totalPrice int64
	for _, region := range regions {
		region.perimeter = calcPerimeter(region, len(plots), len(plots[0]))
		price := int64(len(region.plots)) * int64(region.perimeter)
		totalPrice += price
	}

	// Print regions
	fmt.Printf("regions: %d\n", len(regions))
	for _, region := range regions {
		fmt.Printf("Region %d area x peri %d x %d plots(%d):",
			region.id, len(region.plots), region.perimeter, len(region.plots))
		for _, plot := range region.plots {
			fmt.Printf("  %d,%d", plot.x, plot.y)
		}
		fmt.Println()
	}

	fmt.Printf("total price: %d\n", totalPrice)
	fmt.Println("2024 Day 12 end")
}

// getPlants returns a slice of unique plant types
func getPlants(plots [][]*day12Plot) []string {
	plantsMap := make(map[string]bool)
	for _, row := range plots {
		for _, plot := range row {
			plantsMap[plot.plant] = true
		}
	}

	plants := make([]string, 0, len(plantsMap))
	for plant := range plantsMap {
		plants = append(plants, plant)
	}
	return plants
}

// getRegionsFromPlots identifies and returns connected regions of the same plant type
func getRegionsFromPlots(plots [][]*day12Plot) map[int]*day12Region {
	regions := make(map[int]*day12Region)
	regionID := -1

	for i := range plots {
		for j := range plots[i] {
			plot := plots[i][j]

			// Get valid neighbor coordinates
			neighbors := [][2]int{
				{i - 1, j}, // top
				{i, j + 1}, // right
				{i + 1, j}, // bottom
				{i, j - 1}, // left
			}

			// Filter valid neighbors
			validNeighbors := make([][2]int, 0)
			for _, addr := range neighbors {
				if addr[0] >= 0 && addr[0] < len(plots) &&
					addr[1] >= 0 && addr[1] < len(plots[0]) {
					validNeighbors = append(validNeighbors, addr)
				}
			}

			updatedRegion := -1

			// Check if plot is regionless
			if plot.regionID == -1 {
				// Try to join existing region
				for _, addr := range validNeighbors {
					neighbor := plots[addr[0]][addr[1]]
					if neighbor.plant == plot.plant && neighbor.regionID != -1 {
						plot.regionID = neighbor.regionID
						regions[neighbor.regionID].plots = append(regions[neighbor.regionID].plots, plot)
						updatedRegion = neighbor.regionID
						break
					}
				}

				// Create new region if still regionless
				if plot.regionID == -1 {
					regionID++
					plot.regionID = regionID
					newRegion := &day12Region{
						id:    regionID,
						plots: []*day12Plot{plot},
					}
					regions[regionID] = newRegion
					updatedRegion = regionID
				}
			}

			// Spread to neighbors
			for _, addr := range validNeighbors {
				neighbor := plots[addr[0]][addr[1]]
				if neighbor.plant == plot.plant && neighbor.regionID == -1 {
					neighbor.regionID = plot.regionID
					regions[plot.regionID].plots = append(regions[plot.regionID].plots, neighbor)
					updatedRegion = plot.regionID
				}
			}

			if updatedRegion != -1 {
				regionSpread(regions, plots, updatedRegion)
			}
		}
	}

	return regions
}

// regionSpread recursively spreads a region to connected plots of the same type
func regionSpread(regions map[int]*day12Region, plots [][]*day12Plot, regionID int) {
	region := regions[regionID]
	for spreading := true; spreading; {
		spreading = false
		for _, plot := range region.plots {
			neighbors := [][2]int{
				{plot.x - 1, plot.y}, // top
				{plot.x, plot.y + 1}, // right
				{plot.x + 1, plot.y}, // bottom
				{plot.x, plot.y - 1}, // left
			}

			for _, addr := range neighbors {
				if addr[0] >= 0 && addr[0] < len(plots) &&
					addr[1] >= 0 && addr[1] < len(plots[0]) {
					neighbor := plots[addr[0]][addr[1]]
					if neighbor.plant == plot.plant && neighbor.regionID == -1 {
						neighbor.regionID = regionID
						region.plots = append(region.plots, neighbor)
						spreading = true
						regionSpread(regions, plots, regionID)
					}
				}
			}
			if spreading {
				break
			}
		}
	}
}

// calcPerimeter calculates the perimeter of a region
func calcPerimeter(region *day12Region, matrixRows, matrixCols int) int {
	perimeter := 0
	for _, plot := range region.plots {
		neighbors := [][2]int{
			{plot.x - 1, plot.y}, // top
			{plot.x, plot.y + 1}, // right
			{plot.x + 1, plot.y}, // bottom
			{plot.x, plot.y - 1}, // left
		}

		neighborCount := 0
		for _, addr := range neighbors {
			if addr[0] >= 0 && addr[0] < matrixRows &&
				addr[1] >= 0 && addr[1] < matrixCols {
				// Check if neighbor is in region
				for _, neighborPlot := range region.plots {
					if neighborPlot.x == addr[0] && neighborPlot.y == addr[1] {
						neighborCount++
						break
					}
				}
			}
		}
		perimeter += 4 - neighborCount
	}
	return perimeter
}
