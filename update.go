package main

import (
	os2 "cool-transmission/os"
	"os"
	"path/filepath"
	"time"
)

func main() {
	time.Sleep(3 * time.Second)
	executable, _ := os.Executable()
	rootPath := filepath.Dir(executable)

	srcFile, err := os.Open(filepath.Join(rootPath, "download/cool.download"))
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()
	dstFile, err := os.Create("lad-drop")
	if err != nil {
		panic(err)
	}
	defer dstFile.Close()
	os2.CopyTo(srcFile, filepath.Join(rootPath, "lad-drop"))
	os2.RunMainBin(rootPath)
}
