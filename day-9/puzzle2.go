package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Block struct {
	isFile bool
	fileID int
}

func Checksum(blocks []Block) int {
	result := 0
	for i, block := range blocks {
		result += i * block.fileID // empty space has fileID 0, so they don't matter
	}
	return result
}

func DecompressDiskmap(diskmap string) []Block {
	output := []Block{}
	isFile := true
	nextFileID := 0
	for _, c := range diskmap {
		occupiedSpace, err := strconv.Atoi(string(c))
		if err != nil {
			panic(err)
		}

		block := Block{
			isFile,
			0,
		}
		if isFile {
			block.fileID = nextFileID
			nextFileID++
		}

		for range occupiedSpace {
			output = append(output, block)
		}

		isFile = !isFile
	}

	return output
}

func blocks2str(blocks []Block) string {
	var out strings.Builder
	for _, block := range blocks {
		if block.isFile {
			_, err := out.WriteString(strconv.Itoa(block.fileID))
			if err != nil {
				panic(err)
			}
		} else {
			_, err := out.WriteRune('.')
			if err != nil {
				panic(err)
			}
		}
	}
	return out.String()
}

func MoveFiles(blocks []Block) []Block {
	lastFileRange := [2]int{0, len(blocks) - 1} // start, end block index
	lastFileID := math.MaxInt
	for lastFileID > 0 {

		// fmt.Printf("Current diskmap: %s\n", blocks2str(blocks))
		// Find the first file from the right side (highest fileID)
		for i := lastFileRange[1]; i > -1; i-- {
			if !blocks[i].isFile || lastFileID <= blocks[i].fileID {
				continue
			}
			lastFileRange[1] = i
			lastFileID = blocks[i].fileID
			foundStart := false
			for j := i; j > -1; j-- {
				if blocks[j].fileID != lastFileID {
					lastFileRange[0] = j + 1
					foundStart = true
					break
				}
			}
			// Reached index 0
			if !foundStart {
				lastFileRange[0] = 0
			}
			break
		}

		requiredSpace := lastFileRange[1] - lastFileRange[0] + 1
		// fmt.Printf("Last file range: %v, required space: %d\n", lastFileRange, requiredSpace)

		// find required Space
		freeSpaceSpan := 0
		freeSpaceStartIdx := 0
		for i := 0; i < len(blocks); i++ {
			if blocks[i].isFile {
				freeSpaceSpan = 0
				freeSpaceStartIdx = 0
				continue
			}
			freeSpaceSpan += 1
			if freeSpaceStartIdx == 0 {
				freeSpaceStartIdx = i
			}

			if freeSpaceSpan == requiredSpace {
				break
			}
		}
		// fmt.Printf("Free space span: %d, start index: %d\n", freeSpaceSpan, freeSpaceStartIdx)

		if freeSpaceSpan < requiredSpace {
			continue
		}

		if freeSpaceStartIdx > lastFileRange[1] {
			continue
		}

		for i := range requiredSpace {
			blocks[freeSpaceStartIdx+i] = blocks[lastFileRange[0]+i]
			blocks[lastFileRange[0]+i] = Block{false, 0}
			// fmt.Printf("Moved block of fileID %d from %d to %d\n", blocks[freeSpaceStartIdx+i].fileID, lastFileRange[0]+i, freeSpaceStartIdx+i)
		}

	}
	return blocks
}

func main() {
	inputB, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	diskmap := strings.Trim(string(inputB), "\n")

	// fmt.Printf("Encoded diskmap: %s\n", diskmap)
	blocks := DecompressDiskmap(diskmap)
	// fmt.Printf("Decompressed diskmap: %s\n", blocks2str(blocks))
	movedBlocks := MoveFiles(blocks)
	// fmt.Printf("Moved diskmap: %s\n", blocks2str(movedBlocks))
	fmt.Printf("Checksum: %d\n", Checksum(movedBlocks))
}
