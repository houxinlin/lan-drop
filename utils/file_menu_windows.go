package utils

import (
	"cool-transmission/common"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"os"
	"strconv"
	"strings"
)

func deleteKey(root string, subKey string) {
	key, err := registry.OpenKey(registry.CURRENT_USER, root+subKey, registry.ALL_ACCESS)
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
	fmt.Println("delete key" + name)
	var userName = filterUserName(name)
	deleteKey(`SOFTWARE\Classes\*\shell\Cool文件传输\shell\`, userName+"\\"+"command") //先删除command
	deleteKey(`SOFTWARE\Classes\*\shell\Cool文件传输\shell\`, userName)                //后删除姓名

	deleteKey(`SOFTWARE\Classes\Directory\shell\Cool文件传输\shell\`, userName+"\\"+"command") //先删除command
	deleteKey(`SOFTWARE\Classes\Directory\shell\Cool文件传输\shell\`, userName)
}
func ClearSpecifiedUser(userName string) {
	deleteMenu(userName)
}
func ClearFileMenu() {
	fmt.Println("删除")
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
	for _, value := range info.InfoMap {
		userName := filterUserName(value.Name)
		createSubMenu(userName, value)
		createDirectory(userName, value)
	}
	return nil
}
func filterUserName(username string) string {
	userName := strings.ReplaceAll(username, "\\", "")
	userName = strings.ReplaceAll(userName, "/", "")
	return userName
}

// 创建子菜单项
func createSubMenu(key string, value common.BroadcastInfo) {
	keyPath := "SOFTWARE\\Classes\\*\\shell\\Cool文件传输\\shell\\" + key + "\\command"
	createKey(keyPath, value)
}
func createKey(key string, value common.BroadcastInfo) {
	k, _, err := registry.CreateKey(registry.CURRENT_USER, key, registry.ALL_ACCESS)
	if err != nil {
		panic(err)
	}
	defer k.Close()
	executable, _ := os.Executable()
	if err := k.SetStringValue("", executable+"  %1 "+value.SourceAddress+" "+strconv.Itoa(value.ProtocolPort)); err != nil {
	}
}

func createDirectory(key string, value common.BroadcastInfo) {
	keyPath := "SOFTWARE\\Classes\\Directory\\shell\\Cool文件传输\\shell\\" + key + "\\command"
	createKey(keyPath, value)

}
