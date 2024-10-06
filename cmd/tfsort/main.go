package main

import (
	"cmp"
	"fmt"
	"path/filepath"
	"runtime"
	"slices"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/samber/lo"
)

func main() {
	modulePath := filepath.Join(getAssetPath(), "tests/simple")
	module, _ := tfconfig.LoadModule(modulePath)

	tfVariableMetadata := lo.MapToSlice(module.Variables, func(key string, value *tfconfig.Variable) TfVariableMetadata {
		return TfVariableMetadata{
			Path:         value.Pos.Filename,
			Name:         value.Name,
			StartingLine: value.Pos.Line,
		}
	})

	slices.SortFunc(tfVariableMetadata, func(a, b TfVariableMetadata) int {
		return cmp.Or(
			cmp.Compare(a.Path, b.Path),
			cmp.Compare(a.Name, b.Name),
		)
	})

	groupedMetadata := lo.GroupBy(tfVariableMetadata, func(metadata TfVariableMetadata) string {
		return metadata.Path
	})

	for path, metadata := range groupedMetadata {
		file := TfFile{
			Path:      path,
			Variables: metadata,
		}

		pos := file.FilePos()
		fmt.Println(pos)
	}
}

func getAssetPath() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(b, "../../../assets")
}
