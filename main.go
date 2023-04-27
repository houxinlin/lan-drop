package main

import (
	"cool-transmission/common"
	"cool-transmission/cool"
	coolOs "cool-transmission/os"
	"cool-transmission/ui"
	"cool-transmission/utils"
	"embed"
	"io/ioutil"
	"os"
	"path/filepath"
)

//go:embed res/font/*
var fonts embed.FS

func notifyIfUpdate() {
	if len(os.Args) == 2 {
		if os.Args[1] == "for-update" {
			ui.NotifyMessage("程序已自动从局域网中其他电脑更新更新")
		}
	}
}

func main() {
	coolOs.EnvInit()
	fontData, err := fonts.ReadFile("res/font/font.ttf")
	if err != nil {
		panic(err)
	}
	dir, _ := utils.GetExecutableDir()
	fontPath := filepath.Join(dir, "font", "font.ttf")
	utils.CreateDirectories(fontPath)
	err = ioutil.WriteFile(fontPath, fontData, os.ModePerm)
	if err != nil {
		panic(err)
	}
	os.Setenv("FYNE_FONT", fontPath)
	os.Setenv("FYNE_THEME", "light")

	common.InitContext()

	var startMainServer = len(os.Args) <= 2 //参数是for-date auto
	var startSender = len(os.Args) == 4     //发送者提供的必须参数
	if startMainServer {
		notifyIfUpdate()
		utils.ClearFileMenu()
		ui.RunMain()
		cool.StartMainServer()
		common.ApplicationContext.Application.Run()
	}

	if startSender {
		cool.StartSender()
	}

}
