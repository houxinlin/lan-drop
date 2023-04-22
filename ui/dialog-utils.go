package ui

import (
	"github.com/sqweek/dialog"
)

func ShowMessageDialog(msg string, title string) {
	dialog.Message("%s", msg).Title(title).Info()
}
func ShowDirectorySelectDialog() string {
	directory, _ := dialog.Directory().Title("选择文件夹").Browse()
	return directory
}
