package main

import (
	"cool-transmission/common"
	"cool-transmission/cool"
	"cool-transmission/ui"
	"cool-transmission/utils"
	"embed"
	"io/ioutil"
	"os"
	"path/filepath"
)

//go:embed res/font/*
var fonts embed.FS

func main() {
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
	var startMainServer = len(os.Args) <= 2
	var startSender = len(os.Args) == 4
	if startMainServer {
		utils.ClearFileMenu()
		cool.StartMainServer()
		ui.RunMain()
	}

	if startSender {
		cool.StartSender()
	}

}
