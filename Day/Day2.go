package Day

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

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

func processInputs2(rowList []int64) int {
	isSafeRisingOrFalling := 0

	// Validate that rowList is rising across the row
	//risingValid := true
	for i := 0; i < len(rowList)-1; i++ {
		if rowList[i] < rowList[i+1] {
			isSafeRisingOrFalling = 1
		} else {
			isSafeRisingOrFalling = 0
			//risingValid = false
			break
		}
	}

	// Validate only if rising is false
	if isSafeRisingOrFalling == 0 {
		// Validate that rowList is falling across the row
		for i := 0; i < len(rowList)-1; i++ {
			if rowList[i] > rowList[i+1] {
				isSafeRisingOrFalling = 1
			} else {
				isSafeRisingOrFalling = 0
				break
			}
		}
	}

	// Validate that diff between each is less than or equal to 3
	isSafeGradual := 0
	if isSafeRisingOrFalling == 1 {
		//gradualValid := true
		for i := 0; i < len(rowList)-1; i++ {
			diff := int64(math.Abs(float64(rowList[i+1] - rowList[i])))
			if diff <= 3 {
				isSafeGradual = 1
			} else {
				isSafeGradual = 0
				//gradualValid = false
				break
			}
		}
	}

	return isSafeRisingOrFalling * isSafeGradual
}

func processInputs2pt2(rowList []int64, allowForgive bool) int {
	//fmt.Println("rowList:", rowList)

	forgiveRowNum := -1
	isSafeRisingOrFalling := 0

	// Validate that rowList is rising across the row
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

	// Validate only if rising is false
	if isSafeRisingOrFalling == 0 {
		// Validate that rowList is falling across the row
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

	// Validate that diff between each is less than or equal to 3
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

	if forgiveRowNum != -1 {
		// Remove element at forgiveRowNum
		rowList = append(rowList[:forgiveRowNum], rowList[forgiveRowNum+1:]...)
	}

	if isSafeRisingOrFalling*isSafeGradual == 0 {
		fmt.Printf("rowList: %v isSafeRisingOrFalling: %d isSafeGradual: %d forgiveRowNum: %d final: %d\n",
			rowList, isSafeRisingOrFalling, isSafeGradual, forgiveRowNum, isSafeRisingOrFalling*isSafeGradual)
	}

	return isSafeRisingOrFalling * isSafeGradual
}

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
