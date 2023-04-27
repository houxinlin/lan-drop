package main

import (
	os2 "cool-transmission/os"
	"fmt"
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
		fmt.Println("file open error")
		panic(err)
	}
	defer srcFile.Close()
	destPath := filepath.Join(rootPath, "lad-drop")
	os2.CopyTo(srcFile, destPath)
	os.Chmod(destPath, 0755)
	os2.RunMainBin(rootPath)
}
