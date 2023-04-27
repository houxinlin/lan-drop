package cool

import (
	"cool-transmission/common"
	coolOs "cool-transmission/os"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func doUpdate(url string) bool {
	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}
	parentDir := filepath.Dir(executable)
	filename := filepath.Join(parentDir, "download", "cool.download")
	os.MkdirAll(filepath.Dir(filename), 0755)
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer file.Close()
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return false
	}
	return true
}
func Update(url string) {
	if common.UpdateIng {
		return
	}
	common.UpdateIng = true
	ok := doUpdate(url)
	if ok {
		time.Sleep(3 * time.Second)
		coolOs.RunUpdate()
		fmt.Println("run update")

		os.Exit(0)
	}
	common.UpdateIng = false
}
