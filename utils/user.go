package utils

import (
	"os/user"
	"path/filepath"
)

func GetUserName() string {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	return currentUser.Username
}
func GetUserHome() string {
	currentUser, err := user.Current()
	if err != nil {
		return ""
	}
	homeDir := currentUser.HomeDir
	return homeDir
}
func GetConfigPath() string {
	return filepath.Join(GetUserHome(), ".config", "cool", "transmission", "config.properties")
}
func GetConfig() map[string]string {
	configPath := GetConfigPath()
	CreateDirectories(configPath)
	properties, _ := ReadProperties(configPath)
	return properties
}

func GetValueOrDefault(m map[string]string, key string, defaultValue string) string {
	value, ok := m[key]
	if !ok {
		value = defaultValue
	}
	return value
}
