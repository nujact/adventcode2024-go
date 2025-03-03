package Day

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"utils"
)

func Day1() {
	fmt.Println("Day 1")

	solutionA := solutionA()
	fmt.Println("Solution A:", solutionA)

	solutionB := solutionB()
	fmt.Println("Solution B:", solutionB)
}

func solutionA() int {
	leftList, rightList, shouldReturn := GetInputs()
	if shouldReturn {
		return -1
	}

	sort.Ints(leftList)
	sort.Ints(rightList)

	var solution int = 0
	for i, leftItem := range leftList {
		rightItem := rightList[i]
		dist := int(math.Abs(float64(leftItem) - float64(rightItem)))
		solution += dist
	}

	return solution
}

func GetInputs() ([]int, []int, bool) {
	var left []int
	var right []int
	var day = 1
	var inputFileNameBegins = "input" // "input" or "test"
	lines, err := utils.ReadLinesFromFile(fmt.Sprintf("inputs/%s%d.txt", inputFileNameBegins, day))
	if err != nil {
		fmt.Println("Error:", err)
		return nil, nil, true
	}

	for _, line := range lines {
		slice := strings.Split(line, "   ")

		t, err := strconv.Atoi(slice[0])
		if err == nil {
			left = append(left, t)
		}

		t, err = strconv.Atoi(slice[1])
		if err == nil {
			right = append(right, t)
		}
	}
	return left, right, false
}

func solutionB() int {
	leftList, rightList, shouldReturn := GetInputs()
	if shouldReturn {
		return -1
	}

	var solution int = 0
	for _, leftItem := range leftList {
		count := 0
		for _, rightItem := range rightList {
			if leftItem == rightItem {
				count++
			}
		}

		solution += leftItem * count
	}

	return solution
}
