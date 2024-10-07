package main

import (
	"bufio"
	"cmp"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/samber/lo"
)

type TfFile struct {
	Path      string
	Contents  []string
	Variables []TfVariableMetadata
}

type TfVariableMetadata struct {
	Path             string
	Name             string
	StartingLine     int
	EndingLine       int
	StartingPosition int
	TargetPosition   int
	Contents         []string
}

func LoadModule(path string) []TfFile {
	module, _ := tfconfig.LoadModule(path)
	position := 0

	tfVariableMetadata := lo.MapToSlice(module.Variables, func(key string, value *tfconfig.Variable) TfVariableMetadata {
		metadata := TfVariableMetadata{
			Path:             value.Pos.Filename,
			Name:             value.Name,
			StartingLine:     value.Pos.Line - 1,
			StartingPosition: position,
		}

		position++
		return metadata
	})

	slices.SortFunc(tfVariableMetadata, func(a TfVariableMetadata, b TfVariableMetadata) int {
		return cmp.Compare(a.Name, b.Name)
	})

	for i := 0; i < len(tfVariableMetadata); i++ {
		tfVariableMetadata[i].TargetPosition = i
	}

	groupedMetadata := lo.GroupBy(tfVariableMetadata, func(metadata TfVariableMetadata) string {
		return metadata.Path
	})

	files := make([]TfFile, 0)

	for path, metadata := range groupedMetadata {
		files = append(files, LoadFile(path, metadata))
	}

	return files
}

func LoadFile(path string, metadata []TfVariableMetadata) TfFile {
	contents := readFile(path)
	file := TfFile{
		Path:      path,
		Variables: metadata,
		Contents:  contents,
	}
	file.Tokenize()
	return file
}

func (re TfFile) Sort() {
	lookup := lo.SliceToMap(re.Variables, func(metadata TfVariableMetadata) (int, TfVariableMetadata) {
		return metadata.TargetPosition, metadata
	})

	for _, variable := range lo.Reverse(re.Variables) {
		//Already in the right place
		if variable.StartingPosition == variable.TargetPosition {
			continue
		}

		//todo fix this deleting array contents
		slices.Delete(re.Contents, variable.StartingLine, variable.EndingLine)
		slices.Insert(re.Contents, variable.StartingLine, lookup[variable.TargetPosition].Contents...)
	}

	fmt.Println(re.Contents)
}

func (re TfFile) Tokenize() {
	for i := 0; i < len(re.Variables); i++ {
		startingIndex := re.Variables[i].StartingLine

		openBrackets := 0

		for index, line := range re.Contents[startingIndex:] {
			openBrackets += strings.Count(line, "{")
			openBrackets -= strings.Count(line, "}")

			if openBrackets == 0 {
				re.Variables[i].EndingLine = index + startingIndex + 1
				starting := re.Variables[i].StartingLine
				length := re.Variables[i].EndingLine - starting
				re.Variables[i].Contents = lo.Subset(re.Contents, starting, uint(length))
				break
			}
		}
	}
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
