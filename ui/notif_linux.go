package ui

import (
	"cool-transmission/common"
	"fyne.io/fyne/v2"
)

func NotifyMessage(msg string) {
	common.ApplicationContext.Application.SendNotification(fyne.NewNotification("lan-drop", msg))
}
