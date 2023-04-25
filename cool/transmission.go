package cool

import (
	"cool-transmission/common"
	"cool-transmission/config"
	"cool-transmission/ui"
	"cool-transmission/utils"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
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

var broadcastInfo []byte

// 广播自己的信息
func startBroadcastInfo(port int) {
	for {
		broadcast(port)
		time.Sleep(5 * time.Second)
	}
}
func broadcast(port int) {
	var info common.BroadcastInfo
	info.ProtocolPort = selfProtocolPort
	info.Name = utils.GetCoolUserName()
	info.HttpPort = port
	info.Version = config.AppVersion
	info.SourceAddress = utils.GetOutBoundIP()

	broadcastAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:12556")
	if err != nil {
		log.Fatalf("Error resolving UDP address: %v", err)
	}
	ip := utils.GetOutBoundIP()
	laddr := net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: 3000,
	}
	conn, err := net.DialUDP("udp", &laddr, broadcastAddr)
	if err != nil {
		fmt.Printf("err %s", err)
		return
	}
	defer conn.Close()
	jsonData, err := json.Marshal(info)
	broadcastInfo = jsonData
	conn.Write(broadcastInfo)
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
		fmt.Println(message)
		//如果已经存在
		if !ok || lanUserInfoMap[info.SourceAddress].ProtocolPort != info.ProtocolPort ||
			lanUserInfoMap[info.SourceAddress].Name != info.Name {
			fmt.Println("online:" + info.Name)
			lanUserInfoMap[info.SourceAddress] = info
			menuInfo := utils.FileMenuInfo{InfoMap: lanUserInfoMap}
			menuInfo.WriteServiceFile()
		}
		//如果对方版本大于自己，则更新一下子
		if utils.CompareVersions(info.Version, config.AppVersion) >= 1 {
			go Update("http://" + info.SourceAddress + ":" + strconv.Itoa(info.HttpPort) + "/cool-transmission")
		}
		expireName[info.SourceAddress] = time.Now().Unix()
	}
}
func startExpireListener() {
	interval := 5 * time.Second
	ticker := time.Tick(interval)
	for range ticker {
		for ip, v := range expireName {
			if time.Now().Unix()-v > 7 {
				info := lanUserInfoMap[ip]
				delete(expireName, ip)
				delete(lanUserInfoMap, ip)
				fmt.Printf("%s下线\n", ip)
				utils.ClearSpecifiedUser(info.Name, lanUserInfoMap)
			}
		}
	}
}

// 处理文件传输协议,被接收文件着接收信息
func startProtocolListener() {
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(selfProtocolPort))
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
		createResponse(addr, info)
	}
}

func response(addr *net.UDPAddr, status int, info common.TaskInfo, savePath string) {
	var taskPort int
	taskPort = 0
	if status == 1 { //同意接收文件，窗口一个进度窗口
		var clientConn []net.Conn
		receiveWindowCallback := ui.NewProgressWindow("接收中", func() {
			for _, num := range clientConn {
				num.Close()
			}
		})
		receiveWindowCallback.ClientTcpConnCallback = func(conn net.Conn) {
			clientConn = append(clientConn, conn)
		}
		receiveWindowCallback.Window.Show()
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
func waitTargetResponse(port int, callback common.ReceiveCallback) {
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
		//生产文件发送
		content := common.FileSendContent{FilePath: taskInfo.FileName, DoneCallback: func(data int) {
			callback.Progress(100, "发送成功")
		}, FailCallback: func(err error) {
			callback.StatusCallback(0, "发送失败"+err.Error())
		}, Ip: taskInfo.ToAddress, Port: response.TcpPort,
			Progress: func(size int, current int, progress float64) {
				callback.Progress(progress, fmt.Sprintf("%d/%d", current, size))
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
	httpPort, _ := utils.FindAvailablePort()
	selfProtocolPort = port

	go StartHttpServer(httpPort)    //http服务，用于更新局域网中程序
	go startProtocolListener()      //文件传输服务监听
	go startBroadcastInfo(httpPort) //广播自己的信息
	go startBroadcastListener()     // 接受广播
	go startExpireListener()        //过期监听

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
		progressWindow := ui.NewProgressWindow("发送中", func() {

		})
		progressWindow.Progress(0, "正在等待客户端接受请求")
		go func() {
			port, _ := utils.FindAvailableUDPPort() //找到一个可用的端口号，等待
			selfProtocolPort = port
			taskId := uuid.New()
			taskInfo := common.TaskInfo{SourceAddress: port, FileName: file, ToAddress: args[2], Id: taskId.String()}
			RegisterTask(taskId.String(), taskInfo) //注册一个任务，等待回调
			go waitTargetResponse(port, progressWindow)
			notifySend(taskInfo, args[2], args[3]) //向目标发送通知

		}()
		progressWindow.Window.ShowAndRun()
		protocolUdp.Close()
		os.Exit(0)
	}
}
