package main

import (
	"fmt"
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

func MoveBlocks(blocks []Block) []Block {
	lastBlockIdx := len(blocks) - 1
	firstFreeBlockIdx := -1
	for true {
		// Find the first block from the right side
		for i := lastBlockIdx; i > -1; i-- {
			if !blocks[i].isFile {
				continue
			}
			lastBlockIdx = i
			break
		}

		for i := firstFreeBlockIdx + 1; i < len(blocks); i++ {
			if blocks[i].isFile {
				continue
			}
			firstFreeBlockIdx = i
			break
		}

		if firstFreeBlockIdx >= lastBlockIdx || firstFreeBlockIdx == -1 {
			break
		}

		blocks[firstFreeBlockIdx] = blocks[lastBlockIdx]
		blocks[lastBlockIdx] = Block{false, 0}
		fmt.Printf("Moved block of fileID %d from %d to %d\n", blocks[firstFreeBlockIdx].fileID, lastBlockIdx, firstFreeBlockIdx)
		// fmt.Printf("Current diskmap: %s\n", blocks2str(blocks))
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
	fmt.Printf("Decompressed diskmap: %s\n", blocks2str(blocks))
	movedBlocks := MoveBlocks(blocks)
	fmt.Printf("Moved diskmap: %s\n", blocks2str(movedBlocks))
	fmt.Printf("Checksum: %d\n", Checksum(movedBlocks))
}
