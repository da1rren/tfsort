package main

import (
	"bufio"
	"cmp"
	"fmt"
	"os"
	"slices"
	"strings"
)

type TfVariableMetadata struct {
	Path         string
	Name         string
	StartingLine int
}

type TfFile struct {
	Path      string
	Variables []TfVariableMetadata
}

type VariablePos struct {
	startingLine int
	endingLine   int
}

func (re TfFile) FilePos() []VariablePos {
	lines := readFile(re.Path)

	slices.SortFunc(re.Variables, func(a TfVariableMetadata, b TfVariableMetadata) int {
		return cmp.Compare(a.Name, b.Name)
	})

	variablePos := make([]VariablePos, 0)

	for _, variable := range re.Variables {
		startingIndex := variable.StartingLine - 1

		openBrackets := 0

		for index, line := range lines[startingIndex:] {
			openBrackets += strings.Count(line, "{")
			openBrackets -= strings.Count(line, "}")

			fmt.Printf("%d: %s\n", index+startingIndex, line)

			if openBrackets == 0 {
				variablePos = append(variablePos, VariablePos{
					startingLine: startingIndex,
					endingLine:   index + startingIndex,
				})

				break
			}
		}
	}

	return variablePos
}

func RewriteFile(pos VariablePos[]) {
	
}

func readFile(path string) []string {
	file, _ := os.Open(path)
	defer file.Close()
	lines := make([]string, 0)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text := scanner.Text()
		lines = append(lines, text)
	}

	return lines
}
