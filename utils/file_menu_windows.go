package utils

import (
	"cool-transmission/common"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"os"
	"strconv"
	"strings"
)

const MenuName = "发送文件到"
const MenuKey = "cool_transmission"

const FileMenuRootKey = "SOFTWARE\\Classes\\*\\shell\\" + MenuKey + "\\shell\\"
const DirectoryMenuKey = "SOFTWARE\\Classes\\Directory\\shell\\" + MenuKey + "\\shell\\"

func initKey() {
	var keyPath = `SOFTWARE\Classes\Directory\shell\` + MenuKey + "\\"
	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.ALL_ACCESS)
	if err == nil {
		defer key.Close()
		key.SetStringValue("SubCommands", "")
		key.SetStringValue("MUIVerb", MenuName)
	}

	keyPath = "SOFTWARE\\Classes\\*\\shell\\" + MenuKey
	key, _, err = registry.CreateKey(registry.CURRENT_USER, keyPath, registry.ALL_ACCESS)
	if err == nil {
		defer key.Close()
		key.SetStringValue("SubCommands", "")
		key.SetStringValue("MUIVerb", MenuName)
	}

}

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
	deleteKey(FileMenuRootKey, userName+"\\"+"command") //先删除command
	deleteKey(FileMenuRootKey, userName)
	//后删除姓名
	deleteKey(DirectoryMenuKey, userName+"\\"+"command") //先删除command
	deleteKey(DirectoryMenuKey, userName)
}
func ClearSpecifiedUser(userName string, infoMap map[string]common.BroadcastInfo) {
	deleteMenu(userName)
}
func ClearFileMenu() {

	rooKey, _, err := registry.CreateKey(registry.CURRENT_USER, FileMenuRootKey, registry.ALL_ACCESS)
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
	initKey()

	for _, value := range info.InfoMap {
		userName := filterUserName(value.Name)
		createAllFileMenu(userName, value)
		createDirectory(userName, value)
	}
	return nil
}
func filterUserName(username string) string {
	userName := strings.ReplaceAll(username, "\\", "")
	userName = strings.ReplaceAll(userName, "/", "")
	return userName
}

// 创建所有文件的菜单项
func createAllFileMenu(key string, value common.BroadcastInfo) {
	keyPath := FileMenuRootKey + key + "\\command"
	createKey(keyPath, value)
}

// 创建目录的菜单项
func createDirectory(key string, value common.BroadcastInfo) {
	keyPath := DirectoryMenuKey + key + "\\command"
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
