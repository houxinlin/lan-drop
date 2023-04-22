package ui

import (
	"cool-transmission/common"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type MyWindowCreatedListener struct{}

func NewClientWindow(bind binding.String) fyne.Window {
	var infoLabel = widget.NewLabelWithData(bind)
	var myWindow = common.ApplicationContext.Application.NewWindow("Cool文件传输")
	centerContent := container.New(layout.NewCenterLayout(), infoLabel)
	myWindow.SetContent(centerContent)
	myWindow.Resize(fyne.NewSize(300, 200))
	myWindow.SetFixedSize(true)
	myWindow.CenterOnScreen()
	return myWindow

}
