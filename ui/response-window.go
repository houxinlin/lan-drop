package ui

import (
	"cool-transmission/common"
	"cool-transmission/utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"strings"
)

func isEmptyString(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

func CreateResponseUI(name string, callback func(int, string)) fyne.Window {
	window := common.ApplicationContext.Application.NewWindow("是否接受")
	label := widget.NewLabel(name + "对您扔了一个文件，是否接收?")
	rejectButton := widget.NewButton("拒绝", func() {
		callback(0, "")
	})
	acceptButton := widget.NewButton("接收", func() {
		directory := utils.GetDefaultSaveDirectory()
		//如果是空并且可以创建一个不存在的目录是报错
		if isEmptyString(directory) && utils.CreateDirectories(directory) != nil {
			path := ShowDirectorySelectDialog()
			if path == "" {
				return
			}
			callback(1, path)
			return
		}
		callback(1, directory)
	})

	buttons := container.New(layout.NewHBoxLayout(), rejectButton, acceptButton)
	window.Resize(fyne.NewSize(300, 100))
	window.SetContent(container.New(layout.NewVBoxLayout(), label, buttons))
	window.CenterOnScreen()
	return window
}
