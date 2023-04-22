package ui

import (
	"cool-transmission/common"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewProgressWindow(title string, close func()) common.ReceiveCallback {
	window := common.ApplicationContext.Application.NewWindow(title)
	label := widget.NewLabel("Loading...")
	progressBar := widget.NewProgressBar()
	progressBar.Max = 100
	progressBar.Min = 0
	window.SetContent(container.NewVBox(label, progressBar))
	window.Resize(fyne.NewSize(300, 100))
	window.SetOnClosed(close)
	window.CenterOnScreen()
	return common.ReceiveCallback{Window: window, Progress: func(i float64, msg string) {
		progressBar.SetValue(i)
		label.SetText(msg)
	}, StatusCallback: func(status int, msg string) {
		ShowMessageDialog(msg, "提示")
		window.Close()
	}}

}
