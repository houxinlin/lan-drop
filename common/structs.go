package common

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"net"
)

type BroadcastInfo struct {
	Name          string `json:"name"`
	Version       string `json:"version"`
	SourceAddress string
	ProtocolPort  int `json:"protocolAddress"`
}

type FileData struct {
	FileName  string `json:"filename"`
	FileCount int    `json:"filecount"`
}

type ProtocolResponse struct {
	Status  int    //状态 1 表示同意 0拒绝
	TaskId  string //任务id
	TcpPort int    //对方收到后创建一个tcp，接收文件
}

type TaskInfo struct {
	Id            string //任务Id
	ToAddress     string //目的地址
	FileName      string //文件名
	IsFile        bool
	FileCount     int //文件数量
	SourceAddress int //源目的地址
	SendUserName  string
}

type Context struct {
	Application fyne.App
}

var ApplicationContext Context

func InitContext() {
	ApplicationContext = Context{
		Application: app.New(),
	}
}

type ReceiveCallback struct {
	Progress              func(value int, msg string)
	StatusCallback        func(status int, msg string)
	ClientTcpConnCallback func(conn net.Conn)
}
