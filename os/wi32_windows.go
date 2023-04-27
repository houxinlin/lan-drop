package os

import (
	"cool-transmission/asset"
	_ "embed"
	"github.com/TheTitanrain/w32"
	"golang.org/x/sys/windows/registry"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func EnvInit() {
	open, _ := asset.Resource.ReadFile("bin/update.exe")
	err := os.WriteFile("update.exe", open, 0755)
	if err != nil {
		return
	}
}
func RunUpdate() {
	EnvInit()
	executable, _ := os.Executable()
	parentDir := filepath.Dir(executable)
	exec.Command(filepath.Join(parentDir, "update.exe")).Start()

}
func RunMainBin(root string) {
	cmd := exec.Command(filepath.Join(root, "lad-drop.exe"), "for-update")
	err := cmd.Start()
	fmt.Println(err)
}

func CopyTo(src *os.File, dest string) bool {
	dst, err := os.Create(dest + ".exe")
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
func SetAutoRun(auto bool) {
	const registryPath = `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`
	k, _, _ := registry.CreateKey(registry.CURRENT_USER, registryPath, registry.SET_VALUE)
	defer k.Close()
	executable, _ := os.Executable()
	if auto {
		k.SetStringValue("CoolTransmission", executable+" auto")
	} else {
		k.SetStringValue("CoolTransmission", "")

	}
}
func ShowFile(dir string) {
	w32.ShellExecute(w32.HWND(0), "open", dir, "", "", w32.SW_SHOW)
}
