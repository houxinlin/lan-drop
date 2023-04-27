package ui

func NotifyMessage(msg string) {
	common.ApplicationContext.Application.SendNotification(fyne.NewNotification("lan-drop", msg))
}
