package os

import (
	"cool-transmission/asset"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func EnvInit() {
	open, _ := asset.Resource.ReadFile("bin/update")
	err := os.WriteFile("update", open, 0755)
	if err != nil {
		return
	}
}
func RunUpdate() {
	EnvInit()
	executable, _ := os.Executable()
	parentDir := filepath.Dir(executable)
	exec.Command(filepath.Join(parentDir, "update")).Start()
}
func RunMainBin(root string) {
	cmd := exec.Command(filepath.Join(root, "lad-drop"), "for-update")
	err := cmd.Start()
	fmt.Println(err)
}
func SetAutoRun(bool2 bool) {

}

func CopyTo(src *os.File, dest string) bool {
	dst, err := os.Create(dest)
	if err != nil {
		return false
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		return false
	}
	return true
}
func ShowFile(dir string) {
	cmd := exec.Command("xdg-open", dir)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
