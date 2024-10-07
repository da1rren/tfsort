package main

import (
	"path/filepath"
	"runtime"
)

func main() {
	modulePath := filepath.Join(getAssetPath(), "tests/simple")
	files := LoadModule(modulePath)
	for _, file := range files {
		file.Sort()
	}
}

func getAssetPath() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(b, "../../../assets")
}
