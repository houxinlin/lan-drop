package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const serviceContent = `[Desktop Entry]
Type=Service
X-KDE-ServiceTypes=KonqPopupMenu/Plugin
MimeType=all/all;
Actions=%s;
X-KDE-Submenu=Cool文件传输
Icon=plasma
`

func ClearFileMenu() {
	filename := "cool-transmission.desktop"
	path := filepath.Join(os.Getenv("HOME"), ".local", "share", "kservices5", "ServiceMenus", filename)
	file, _ := os.Create(path)
	file.WriteString("")
	file.Close()
}
func (info *FileMenuInfo) WriteServiceFile() error {
	filename := "cool-transmission.desktop"
	path := filepath.Join(os.Getenv("HOME"), ".local", "share", "kservices5", "ServiceMenus", filename)
	CreateDirectories(path)
	file, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return err
	}
	executablePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	defer file.Close()
	actions := make([]string, len(info.InfoMap))
	for i := 0; i < len(info.InfoMap); i++ {
		actions[i] = "child" + strconv.Itoa(i)
	}

	file.WriteString(fmt.Sprintf(serviceContent, strings.Join(actions, ";")))
	var count = 0
	for _, value := range info.InfoMap {
		// 写入文件内容
		item := fmt.Sprintf("[Desktop Action %s]\nName=%s\nIcon=preferences-desktop-plasma\nExec=%s %%u %s %s \n",
			actions[count], value.Name, executablePath, value.SourceAddress, strconv.Itoa(value.ProtocolPort))
		_, err = file.WriteString(item)
		count++
	}
	fmt.Println("KO")
	return nil
}
