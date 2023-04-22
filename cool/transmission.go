package cool

import (
	"cool-transmission/common"
	"cool-transmission/ui"
	"cool-transmission/utils"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var protocolUdp *net.UDPConn

func notifySend(info common.TaskInfo, host string, port string) error {
	info.SendUserName = utils.GetCoolUserName()
	jsonData, err := json.Marshal(info)
	if err != nil {
		return err
	}
	udpAddr, err := net.ResolveUDPAddr("udp", host+":"+port)
	if err != nil {
		return err
	}
	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}
	defer udpConn.Close()
	_, err = udpConn.Write(jsonData)
	if err != nil {
		return err
	}
	fmt.Printf("文件发送Request: %s\n", string(jsonData))
	return nil
}

var lanUserInfoMap = make(map[string]common.BroadcastInfo)
var expireName = make(map[string]int64)
var selfProtocolPort int

// 广播自己的信息
func startBroadcastInfo(port int) {
	var info common.BroadcastInfo
	info.ProtocolPort = port
	info.Name = utils.GetCoolUserName()
	broadcastAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:12556")
	if err != nil {
		log.Fatalf("Error resolving UDP address: %v", err)
	}
	ip, _ := utils.GetOutBoundIP()

	laddr := net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: 3000,
	}
	conn, err := net.DialUDP("udp", &laddr, broadcastAddr)
	if err != nil {
	}
	defer conn.Close()

	for {
		jsonData, err := json.Marshal(info)
		if err != nil {
		}

		_, err = conn.Write(jsonData)
		if err != nil {
		}
		time.Sleep(5 * time.Second)
	}
}

// 接收局域网广播的用户信息
func startBroadcastListener() {
	listenAddr, err := net.ResolveUDPAddr("udp", ":12556")
	if err != nil {
		log.Fatalf("Error resolving UDP address: %v", err)
	}
	conn, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}
	for {
		buffer := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			continue
		}
		message := string(buffer[:n])
		var info common.BroadcastInfo
		err = json.Unmarshal([]byte(message), &info)
		info.SourceAddress = addr.IP.String()

		_, ok := lanUserInfoMap[info.SourceAddress]
		//如果已经存在
		if !ok && lanUserInfoMap[info.Name].ProtocolPort != info.ProtocolPort {
			fmt.Println("上线" + info.Name)
			lanUserInfoMap[info.SourceAddress] = info
			menuInfo := utils.FileMenuInfo{InfoMap: lanUserInfoMap}
			menuInfo.WriteServiceFile()
		}
		expireName[info.Name] = time.Now().Unix()
	}
}
func startExpireListener() {
	interval := 5 * time.Second
	ticker := time.Tick(interval)
	for range ticker {
		for k, v := range expireName {
			if time.Now().Unix()-v > 7 {
				delete(expireName, k)
				delete(lanUserInfoMap, k)
				fmt.Printf("%s下线", k)
			}
		}
	}
}

// 处理文件传输协议,被接收文件着接收信息
func startProtocolListener(port int) {
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	if err != nil {
		return
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return
	}
	defer udpConn.Close()
	buffer := make([]byte, 4096)
	for {
		n, addr, err := udpConn.ReadFromUDP(buffer)
		if err != nil {
			continue
		}
		var info common.TaskInfo
		err = json.Unmarshal(buffer[:n], &info)
		if err != nil {
			continue
		}
		//创建一个响应的提示窗口
		fmt.Println("收到文件发送请求")
		createResponse(addr, info)
	}
}
func createReceiveWindow(close func()) common.ReceiveCallback {
	window := common.ApplicationContext.Application.NewWindow("接收中")
	label := widget.NewLabel("Loading...")
	progressBar := widget.NewProgressBar()
	progressBar.Max = 100
	progressBar.Min = 0
	window.SetContent(container.NewVBox(label, progressBar))
	window.Resize(fyne.NewSize(300, 100))
	window.SetOnClosed(close)
	window.CenterOnScreen()
	window.Show()
	return common.ReceiveCallback{Progress: func(i int, msg string) {
		if i >= 100 {
			window.SetTitle("接收完毕")
		}
		progressBar.SetValue(float64(i))
		label.SetText(msg)
	}, StatusCallback: func(status int, msg string) {
		ui.ShowMessageDialog(msg, "提示")
		window.Close()
	}}

}
func response(addr *net.UDPAddr, status int, info common.TaskInfo, savePath string) {
	var taskPort int
	taskPort = 0
	if status == 1 { //同意接收文件，窗口一个进度窗口
		var clientConn []net.Conn
		receiveWindowCallback := createReceiveWindow(func() {
			for _, num := range clientConn {
				num.Close()
			}
		})
		receiveWindowCallback.ClientTcpConnCallback = func(conn net.Conn) {
			clientConn = append(clientConn, conn)
		}
		port, _ := CreateNewReceiveTask(savePath, receiveWindowCallback)
		taskPort = port
	}
	sendResponseTo(common.ProtocolResponse{TaskId: info.Id, Status: status, TcpPort: taskPort}, addr.IP.String(), strconv.Itoa(info.SourceAddress))

}
func createResponse(addr *net.UDPAddr, info common.TaskInfo) {
	var window fyne.Window
	//自动接收
	if utils.IsAutoReceive() {
		response(addr, 1, info, utils.GetDefaultSaveDirectory())
		return
	}
	//确认接收
	window = ui.CreateResponseUI(info.SendUserName, func(status int, savePath string) {
		window.Close()
		response(addr, status, info, savePath)
	})
	window.Show()
}



// 等待客户端响应
func waitTargetResponse(port int, callback func(msg string)) {
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	if err != nil {
		return
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return
	}
	protocolUdp = udpConn
	defer udpConn.Close()
	buffer := make([]byte, 1024)
	n, _, err := udpConn.ReadFromUDP(buffer)
	if err != nil {
		return
	}
	var response common.ProtocolResponse
	err = json.Unmarshal(buffer[:n], &response)
	if err != nil {
		return
	}
	//客户同意接收文件
	if response.Status == 1 {
		taskInfo := GetTaskInfo(response.TaskId)
		content := SendTaskContent{FilePath: taskInfo.FileName, doneCallback: func(data int) {
			callback("发送成功")
		}, failCallback: func(err error) {
			callback("发送失败" + err.Error())
		}, Ip: taskInfo.ToAddress, Port: response.TcpPort, Progress: func(size int, current int, progress float64) {
			callback(fmt.Sprintf("%d/%d %.2f%%", current, size, progress))
		}}
		go SendFileToTarget(content)
		return
	}
	ui.ShowMessageDialog("对方拒绝接收文件", "提示")
	os.Exit(0)
	return

}

func sendResponseTo(res common.ProtocolResponse, host string, port string) error {
	jsonData, err := json.Marshal(res)
	if err != nil {
		return err
	}
	udpAddr, err := net.ResolveUDPAddr("udp", host+":"+port)
	if err != nil {
		return err
	}
	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}
	defer udpConn.Close()
	_, err = udpConn.Write(jsonData)
	if err != nil {
		return err
	}
	return nil

}
func startServer() {
	port, _ := utils.FindAvailableUDPPort()
	selfProtocolPort = port

	go startProtocolListener(port)
	go startBroadcastInfo(port)
	go startBroadcastListener()
	go startExpireListener()

	go func() {
	}()
}
// StartMainServer 启动主服务
func StartMainServer() {
	go startServer()
}

func StartSender() {
	args := os.Args
	file := args[1]
	_, err := os.Stat(file)
	if err == nil {
		bindMessage := binding.NewString()
		window := ui.NewClientWindow(bindMessage)
		bindMessage.Set("等待客户端响应")
		go func() {
			port, _ := utils.FindAvailableUDPPort() //找到一个可用的端口号，等待
			selfProtocolPort = port
			taskId := uuid.New()
			taskInfo := common.TaskInfo{SourceAddress: port, FileName: file, ToAddress: args[2], Id: taskId.String()}
			RegisterTask(taskId.String(), taskInfo) //注册一个任务，等待回调
			go waitTargetResponse(port, func(msg string) {
				bindMessage.Set(msg)
			})
			notifySend(taskInfo, args[2], args[3]) //向目标发送通知

		}()

		window.ShowAndRun()
		protocolUdp.Close()
		os.Exit(0)
	}
}
