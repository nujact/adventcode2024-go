// Package Day implements solutions for Advent of Code 2024 challenges.
// Day9 specifically solves a puzzle involving disk defragmentation.
// The puzzle involves a disk storage system where files need to be rearranged
// to optimize their positions. Each file has a unique ID and can be moved to
// any contiguous empty space that can fit its length. The goal is to minimize
// fragmentation by moving files towards the beginning of the disk.
package Day

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// DiskMap represents the disk storage with file positions and empty spaces.
// The disk is represented as a linear array where:
//   - Positive integers represent file IDs
//   - -1 represents empty space
type DiskMap struct {
	Map []int // Slice representing disk positions, -1 for empty space, fileId for files
}

// day9NewDiskMap creates a new DiskMap from the input string.
// It parses alternating file size and empty space values to create the disk layout.
// The input string format alternates between:
//   - Even indices: File sizes
//   - Odd indices: Empty space sizes
//
// Parameters:
//   - inputString: String containing alternating file and space sizes
//
// Returns:
//   - *DiskMap: Initialized disk map with files and empty spaces
func day9NewDiskMap(inputString string) *DiskMap {
	diskMap := &DiskMap{
		Map: make([]int, 0),
	}

	fileID := 0
	needNewFileID := false

	for i := 0; i < len(inputString); i++ {
		currentInt, _ := strconv.Atoi(string(inputString[i]))
		if i%2 == 0 {
			// Odd, file size
			if needNewFileID {
				fileID++
				needNewFileID = false
			}
			for takeCtr := 0; takeCtr < currentInt; takeCtr++ {
				diskMap.Map = append(diskMap.Map, fileID)
			}
		} else {
			// Even, empty space
			for takeCtr := 0; takeCtr < currentInt; takeCtr++ {
				diskMap.Map = append(diskMap.Map, -1)
			}
			needNewFileID = true
		}
	}

	return diskMap
}

// Print displays the current state of the disk map to standard output.
// The display format is:
//   - "." for empty spaces
//   - File ID number for occupied spaces
//
// This method is useful for visualizing the disk state during defragmentation
// and debugging the file movement process.
func (d *DiskMap) Print() {
	for i := 0; i < len(d.Map); i++ {
		output := d.Map[i]
		if output == -1 {
			fmt.Print(". ")
		} else {
			fmt.Printf("%d ", output)
		}
	}
	fmt.Println()
}

// getLastFileID returns the highest fileID present in the disk map.
// This is used to determine the range of files that need to be processed
// during defragmentation.
//
// Returns:
//   - int: The highest file ID found, or -1 if no files exist
func (d *DiskMap) getLastFileID() int {
	lastFileID := -1
	for i := len(d.Map) - 1; i >= 0; i-- {
		if d.Map[i] != -1 {
			lastFileID = d.Map[i]
			break
		}
	}
	return lastFileID
}

// getFileLength returns the length of a file with the given fileID.
// The length is calculated by counting consecutive positions containing
// the same file ID.
//
// Parameters:
//   - fileID: The ID of the file to measure
//
// Returns:
//   - int: The length of the file in disk positions
func (d *DiskMap) getFileLength(fileID int) int {
	length := 0
	foundStart := false
	for i := 0; i < len(d.Map); i++ {
		if d.Map[i] == fileID {
			length++
			foundStart = true
		}
		if foundStart && d.Map[i] != fileID {
			break
		}
	}
	return length
}

// getFirstPositionByFileID returns the first position where the given fileID appears.
// This is used to locate files before moving them during defragmentation.
//
// Parameters:
//   - fileID: The ID of the file to locate
//
// Returns:
//   - int: The first position where the file appears, or -1 if not found
func (d *DiskMap) getFirstPositionByFileID(fileID int) int {
	for i := 0; i < len(d.Map); i++ {
		if d.Map[i] == fileID {
			return i
		}
	}
	return -1
}

// getFirstEmptyPositionByLength finds the first empty position that can fit a file of given length.
// The method searches for a contiguous sequence of empty spaces that can accommodate the file.
//
// Parameters:
//   - fileLength: The length of the file that needs to be placed
//
// Returns:
//   - int: The starting position of suitable empty space, or -1 if no suitable space found
func (d *DiskMap) getFirstEmptyPositionByLength(fileLength int) int {
	pos := -1
	foundEmpty := false
	emptyCtr := 0

	for i := 0; i < len(d.Map); i++ {
		if d.Map[i] == -1 {
			if !foundEmpty {
				pos = i
			}
			foundEmpty = true
			emptyCtr++
		} else {
			if foundEmpty && emptyCtr >= fileLength {
				break
			}
			foundEmpty = false
			emptyCtr = 0
			pos = -1
		}
	}
	return pos
}

// DefragmentWholeFilesOnce performs one pass of defragmentation on whole files.
// The method processes files from highest ID to lowest, attempting to move each file
// to the earliest possible position that can accommodate its entire length.
//
// The defragmentation process:
// 1. Identifies each file by ID
// 2. Determines the file's current position and length
// 3. Finds the earliest empty space that can fit the entire file
// 4. Moves the file if a better position is found
//
// Returns:
//   - bool: true if any files were moved during this pass, false otherwise
func (d *DiskMap) DefragmentWholeFilesOnce() bool {
	maxFileID := d.getLastFileID()
	moved := false

	for fileIDCtr := maxFileID; fileIDCtr > -1; fileIDCtr-- {
		fileIDCtrLength := d.getFileLength(fileIDCtr)
		positionFileCtr := d.getFirstPositionByFileID(fileIDCtr)
		positionFirstEmptyThatFitsCtr := d.getFirstEmptyPositionByLength(fileIDCtrLength)

		if positionFirstEmptyThatFitsCtr < positionFileCtr && positionFirstEmptyThatFitsCtr != -1 {
			// Empty space that fits this file prior to file pos, so move it
			for i := 0; i < fileIDCtrLength; i++ {
				d.Map[positionFirstEmptyThatFitsCtr+i] = d.Map[positionFileCtr+i]
				d.Map[positionFileCtr+i] = -1
			}
			moved = true
		}
	}
	return moved
}

// CalculateChecksum computes the checksum of the disk map based on file positions.
// The checksum is calculated by:
// 1. For each file position, multiply the file ID by its position
// 2. Sum all these products
//
// Returns:
//   - string: The calculated checksum as a string
func (d *DiskMap) CalculateChecksum() string {
	var checksum int64
	for i := 0; i < len(d.Map); i++ {
		if d.Map[i] != -1 {
			checksum += int64(d.Map[i]) * int64(i)
		}
	}
	return strconv.FormatInt(checksum, 10)
}

// Day9 solves the Advent of Code 2024 Day 9 puzzle.
// The puzzle involves optimizing file storage on a disk by moving files
// to reduce fragmentation. The solution follows these steps:
// 1. Read the input describing initial file and space positions
// 2. Create a disk map representation
// 3. Perform defragmentation by moving files towards the start
// 4. Calculate a checksum based on final file positions
//
// The final answer is the checksum value after defragmentation.
func Day9() {
	fmt.Println("2024 Day9 start")

	// Open file and gather raw inputs
	inputMemory, shouldReturn := day9GetInput()
	if shouldReturn {
		return
	}

	fmt.Printf("inputMemory size: %d\n", len(inputMemory))

	// Create disk map from input
	diskMap := day9NewDiskMap(inputMemory)

	// Print initial disk map
	diskMap.Print()

	// Defragment disk map
	diskMap.DefragmentWholeFilesOnce()

	// Print defragmented disk map
	fmt.Println("defragmented disk map")
	diskMap.Print()

	// Calculate checksum
	checksum := diskMap.CalculateChecksum()
	fmt.Printf("checksum: %s\n", checksum)

	fmt.Println("2024 Day9 end")
}

// day9GetInput reads the puzzle input file and returns its contents.
// The input file contains a string of digits where:
//   - Even-indexed digits represent file sizes
//   - Odd-indexed digits represent empty space sizes
//
// Returns:
//   - string: The contents of the input file
//   - bool: true if an error occurred, false otherwise
//
// The function handles file operations and error checking, returning appropriate values
// to indicate success or failure of the input reading process.
func day9GetInput() (string, bool) {
	var inputFileNameBegins = "input" // "input" or "test"
	file, err := os.Open(fmt.Sprintf("../inputs/%s9.txt", inputFileNameBegins))
	if err != nil {
		fmt.Println(err)
		return "", true
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var inputMemory string
	for scanner.Scan() {
		inputMemory += scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return "", true
	}
	return inputMemory, false
}
