package utils

import (
	"cool-transmission/common"
	"strings"
)

func IsAutoReceive() bool {
	config := GetConfig()
	return strings.Compare("false", strings.ToLower(GetValueOrDefault(config, common.AutoReceive, "false"))) != 0
}
func IsAutoOpen() bool {
	config := GetConfig()
	return strings.Compare("false", strings.ToLower(GetValueOrDefault(config, common.AutoOpenFolder, "false"))) != 0
}
func IsAutoRun() bool {
	config := GetConfig()
	return strings.Compare("false", strings.ToLower(GetValueOrDefault(config, common.AutoRun, "false"))) != 0
}

func GetDefaultSaveDirectory() string {
	config := GetConfig()
	return GetValueOrDefault(config, common.DefaultSaveDir, "")
}
func GetCoolUserName() string {
	config := GetConfig()
	return GetValueOrDefault(config, common.UserName, GetUserName())
}
