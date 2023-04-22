package utils

import (
	"cool-transmission/common"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"os"
	"strconv"
	"strings"
)

func deleteKey(subKey string) {
	path := `SOFTWARE\Classes\*\shell\Cool文件传输\shell\`
	key, err := registry.OpenKey(registry.CURRENT_USER, path+subKey, registry.ALL_ACCESS)
	if err != nil {
		fmt.Println("OpenKey error:", err)
		return
	}
	defer key.Close()
	if err := registry.DeleteKey(key, ""); err != nil {
		fmt.Println("DeleteKey error:", err)
		return
	}
}

func deleteMenu(name string) {
	deleteKey(name + "\\" + "command") //先删除command
	deleteKey(name)                    //后删除姓名
}
func ClearSpecifiedUser(userName string) {
	deleteMenu(userName)
}
func ClearFileMenu() {
	path := `SOFTWARE\Classes\*\shell\Cool文件传输\shell\`
	rooKey, err := registry.OpenKey(registry.CURRENT_USER, path, registry.ALL_ACCESS)
	defer rooKey.Close()
	if err != nil {
		fmt.Println("OpenKey error:", err)
		return
	}
	names, err := rooKey.ReadSubKeyNames(-1)
	if err != nil {
		return
	}
	for _, name := range names {
		deleteMenu(name)
	}

}
func (info *FileMenuInfo) WriteServiceFile() error {
	appPath, _ := os.Executable()
	for _, value := range info.InfoMap {
		userName := strings.ReplaceAll(value.Name, "\\", "")
		userName = strings.ReplaceAll(userName, "/", "")
		createSubMenu(userName, appPath, value)
		createDirectory(userName, appPath, value)
	}
	return nil
}

// 创建子菜单项
func createSubMenu(key string, appPath string, value common.BroadcastInfo) {
	keyPath := "SOFTWARE\\Classes\\*\\shell\\Cool文件传输\\shell\\" + key + "\\command"
	k, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.ALL_ACCESS)
	if err != nil {
		panic(err)
	}
	defer k.Close()
	if err := k.SetStringValue("", appPath+"  %1 "+value.SourceAddress+" "+strconv.Itoa(value.ProtocolPort)); err != nil {
		panic(err)
	}

}

func createDirectory(key string, filename string, value common.BroadcastInfo) {
	k, _, err := registry.CreateKey(registry.CURRENT_USER, `SOFTWARE\Classes\Directory\Background\Shell\Cool文件传输`, registry.ALL_ACCESS)
	if err != nil {
		panic(err)
	}
	defer k.Close()
	if err := k.SetStringValue("", key); err != nil {
		panic(err)
	}

}
