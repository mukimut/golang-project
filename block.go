package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type ngninxBlock struct {
	Start        int
	End          int
	AllLines     *[]*string
	NestedBlocks []*ngninxBlock
	totalBlocks  int
}

type ngnixBlocks struct {
	blocks   *[]*ngninxBlock
	allLines *[]*string
}

func main() {
	file, _ := os.Open("ngnix.conf")
	scanner := bufio.NewScanner(file)
	var lines []*string
	var blockList []*ngninxBlock

	for scanner.Scan() {
		currentLine := scanner.Text()
		lines = append(lines, &currentLine)
	}
	totalLines := len(lines)

	for i := 0; i < totalLines; {
		if strings.Contains(*lines[i], "{") {
			currentBlock := getAllBlocks(lines, i, totalLines, 0, 1)
			blockList = append(blockList, currentBlock)
			i = currentBlock.End - 1
			continue
		}

		i++
	}

	searchResult := getNgnixBlocks(blockList, "80")

	for _, line := range *searchResult.allLines {
		var lineString string = *line
		fmt.Println(lineString)
	}
}

func getNgnixBlock(lines *[]*string, startIndex, endIndex, recursionMax int) *ngninxBlock {
	return getAllBlocks(*lines, startIndex, endIndex, 0, recursionMax)
}

func getAllBlocks(lines []*string, start int, lineCount int, currentRecursion, maxRecursion int) *ngninxBlock {
	var blockLines []*string
	var nestedBlocks []*ngninxBlock
	var end int
	totalBlocks := 0

	blockLines = append(blockLines, lines[start])
	start++

	for i := start; i < lineCount; {
		blockLines = append(blockLines, lines[i])
		if strings.Contains(*lines[i], "}") {
			end = i + 1
			break
		}

		if strings.Contains(*lines[i], "{") {
			newBlock := getAllBlocks(lines, i, lineCount, currentRecursion+1, maxRecursion)
			i = newBlock.End

			if maxRecursion == -1 || currentRecursion <= maxRecursion {
				nestedBlocks = append(nestedBlocks, newBlock)
				totalBlocks++
			}
			continue
		}

		i++
	}

	// if maxRecursion == -1 || currentRecursion <= maxRecursion {
	// 	testPrinter(*lines[start-1], start)
	// 	testPrinter(*lines[end-1], end)
	// 	fmt.Println(currentRecursion)
	// 	fmt.Println(maxRecursion)
	// }

	return &ngninxBlock{start, end, &blockLines, nestedBlocks, totalBlocks}
}

func getNgnixBlocks(allBlocks []*ngninxBlock, config string) *ngnixBlocks {
	var foundBlocks []*ngninxBlock
	var foundLines []*string

	for _, block := range allBlocks {
		if block.totalBlocks > 0 {
			nestedSearchResult := getNgnixBlocks(block.NestedBlocks, config)
			foundBlocks = append(foundBlocks, *nestedSearchResult.blocks...)
			foundLines = append(foundLines, *nestedSearchResult.allLines...)
		}

		for _, line := range *block.AllLines {
			if strings.Contains(*line, config) {
				foundBlocks = append(foundBlocks, block)
				foundLines = append(foundLines, *block.AllLines...)
				break
			}
		}
	}

	return &ngnixBlocks{&foundBlocks, &foundLines}
}

func testPrinter(line string, value int) string {
	fmt.Print("Text is " + line + "|| ")
	fmt.Print("Value is: ")
	fmt.Println(value)
	return line
}
