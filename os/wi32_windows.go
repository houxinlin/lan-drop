package os

import (
	"github.com/TheTitanrain/w32"
	"golang.org/x/sys/windows/registry"
	"os"
)

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
