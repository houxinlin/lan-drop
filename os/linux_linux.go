package os

import (
	"cool-transmission/asset"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
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
func SetAutoRun(auto bool) {
	currentUser, err := user.Current()
	targetDir := currentUser.HomeDir + "/.config/autostart"
	targetFile := targetDir + "/lad-drop.desktop"
	if !auto {
		os.Remove(targetFile)
		return
	}
	if err != nil {
		return
	}

	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		return
	}
	file, err := os.Create(targetFile)
	if err != nil {
		return
	}
	defer file.Close()
	executable, err := os.Executable()
	content := fmt.Sprintf(`[Desktop Entry]
Encoding=UTF-8
Name=lad-drop
Comment=lad-drop
Comment[zh_CN]=lad-drop
Exec=%s
Icon=%s
Keywords=lad-drop
Type=Application
Terminal=false`, executable, executable)
	_, err = file.WriteString(content)
	if err != nil {
		return
	}
	os.Chmod(targetFile, 0755)
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
