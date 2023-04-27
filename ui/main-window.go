package ui

import (
	"cool-transmission/common"
	coolOS "cool-transmission/os"
	"cool-transmission/utils"
	_ "embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"os"
	"strconv"
)

//go:embed  icon.png
var iconData []byte

func init() {

}

func RunMain() {
	var resourceIconPng = &fyne.StaticResource{
		StaticName:    "Icon.png",
		StaticContent: iconData,
	}
	common.ApplicationContext.Application.SetIcon(resourceIconPng)

	config := utils.GetConfig()

	var mainWindow = common.ApplicationContext.Application.NewWindow("Cool文件传输")
	mainWindow.Resize(fyne.NewSize(500, 300))
	autoRun := widget.NewCheck("开机自启", func(b bool) {
		config[common.AutoRun] = strconv.FormatBool(b)
		coolOS.SetAutoRun(b)
		utils.SaveProperties(config)
	})

	autoReceiveCheckbox := widget.NewCheck("是否自动接收文件", func(b bool) {
		config[common.AutoReceive] = strconv.FormatBool(b)
		utils.SaveProperties(config)
	})
	autoOpenFolderCheckbox := widget.NewCheck("接受完毕后自动打开文件夹", func(b bool) {
		config[common.AutoOpenFolder] = strconv.FormatBool(b)
		utils.SaveProperties(config)
	})
	autoRun.SetChecked(utils.IsAutoRun())
	autoOpenFolderCheckbox.SetChecked(utils.IsAutoOpen())
	autoReceiveCheckbox.SetChecked(utils.IsAutoReceive())

	defaultSavePathInput := widget.NewEntry()
	defaultSavePathInput.SetText(utils.GetDefaultSaveDirectory())
	defaultSavePathInput.SetPlaceHolder("默认保存路径")
	defaultSavePathInput.OnChanged = func(path string) {
		config[common.DefaultSaveDir] = path
		defaultSavePathInput.SetText(path)
		utils.SaveProperties(config)
	}
	changeDirectoryButton := widget.NewButton("更改目录", func() {
		path := ShowDirectorySelectDialog()
		config[common.DefaultSaveDir] = path
		defaultSavePathInput.SetText(path)
		utils.SaveProperties(config)
	})
	var top float32 = 20
	defaultSavePathContainer := container.NewWithoutLayout(defaultSavePathInput, changeDirectoryButton)

	autoRun.Move(fyne.NewPos(0, top))
	autoRun.Resize(fyne.NewSize(100, 10))

	top += 30

	autoReceiveCheckbox.Move(fyne.NewPos(0, top))
	autoReceiveCheckbox.Resize(fyne.NewSize(100, 10))

	top += 30
	autoOpenFolderCheckbox.Move(fyne.NewPos(0, top))
	autoOpenFolderCheckbox.Resize(fyne.NewSize(100, 10))

	defaultSavePathContainer.Move(fyne.NewPos(0, top))
	defaultSavePathInput.Resize(fyne.NewSize(300, 35))
	defaultSavePathInput.Move(fyne.NewPos(0, 40))

	changeDirectoryButton.Resize(fyne.NewSize(100, 40))
	changeDirectoryButton.Move(fyne.NewPos(310, 40))
	defaultSavePathContainer.Refresh()

	firstTabContent := container.NewWithoutLayout(autoRun, autoReceiveCheckbox, autoOpenFolderCheckbox, defaultSavePathContainer)

	userNameLabel := widget.NewLabel("用户名")
	defaultUserInput := widget.NewEntry()
	defaultUserInput.SetText(utils.GetCoolUserName())
	defaultUserInput.SetPlaceHolder("用户名")

	changeUserNameButton := widget.NewButton("保存", func() {
		config[common.UserName] = defaultUserInput.Text
		utils.SaveProperties(config)
	})

	userNameLabel.Move(fyne.NewPos(0, 0))
	defaultUserInput.Move(fyne.NewPos(7, 25))
	defaultUserInput.Resize(fyne.NewSize(300, 35))
	changeUserNameButton.Move(fyne.NewPos(310, 25))
	changeUserNameButton.Resize(fyne.NewSize(100, 35))
	userContainer := container.NewWithoutLayout(userNameLabel, defaultUserInput, changeUserNameButton)

	userTabContent := container.NewWithoutLayout(userContainer)

	tabs := container.NewAppTabs(
		container.NewTabItem("基本设置", firstTabContent),
		container.NewTabItem("用户设置", userTabContent),
	)

	if desk, ok := common.ApplicationContext.Application.(desktop.App); ok {
		m := fyne.NewMenu("App",
			fyne.NewMenuItem("设置", func() {
				mainWindow.Show()
			}))
		desk.SetSystemTrayMenu(m)
	}

	mainWindow.SetCloseIntercept(func() {
		mainWindow.Hide()
	})
	mainWindow.CenterOnScreen()
	mainWindow.SetContent(tabs)
	if len(os.Args) == 1 {
		mainWindow.Show()
	}
}
