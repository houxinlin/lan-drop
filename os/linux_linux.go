package os

import (
	"fmt"
	"os/exec"
)

func SetAutoRun(bool2 bool) {

}

func ShowFile(dir string) {
	cmd := exec.Command("xdg-open", dir)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
