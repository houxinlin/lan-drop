package cool

import (
	"cool-transmission/common"
	coolOs "cool-transmission/os"
	"cool-transmission/utils"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

func CreateNewReceiveTask(savePath string, callback common.ReceiveCallback) (int, net.Listener) {
	port, _ := utils.FindAvailablePort()
	listener, _ := net.Listen("tcp", ":"+strconv.Itoa(port))
	go begin(listener, savePath, callback)
	return port, listener
}
func begin(listener net.Listener, path string, callback common.ReceiveCallback) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		callback.ClientTcpConnCallback(conn) //用于断开链接
		go func() {
			if err := handleClient(conn, path, callback); err != nil {
				listener.Close() //有错误后断开链接
				callback.StatusCallback(0, err.Error())
			}
		}()
	}
}
func handleClient(conn net.Conn, path string, callback common.ReceiveCallback) error {
	defer conn.Close()
	var numFiles uint32
	err := binary.Read(conn, binary.BigEndian, &numFiles)
	if err != nil {
		return err
	}
	for i := uint32(0); i < numFiles; i++ {
		var fileType byte
		err = binary.Read(conn, binary.BigEndian, &fileType)
		if err != nil {
			return err
		}

		var nameLen uint32
		err = binary.Read(conn, binary.BigEndian, &nameLen)
		if err != nil {
			return err
		}

		nameBuf := make([]byte, nameLen)
		_, err = io.ReadFull(conn, nameBuf)
		if err != nil {
			return err
		}

		name := string(nameBuf)
		if fileType == 0 {
			err = os.MkdirAll(filepath.Join(path, name), os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			var fileSize uint32
			err = binary.Read(conn, binary.BigEndian, &fileSize)
			if err != nil {
				return err
			}
			receiveBuf := make([]byte, 4096)
			fullPath := filepath.Join(path, name)
			utils.CreateDirectories(fullPath)
			targetFile, err := os.Create(fullPath)
			if err != nil {
				return err
			}
			var remaining = int(fileSize)
			for {
				if remaining == 0 {
					break
				}
				if remaining < len(receiveBuf) {
					receiveBuf = make([]byte, remaining)
				}
				n, err := conn.Read(receiveBuf)
				if err != nil && err != io.EOF {
					return err
				}
				if n == 0 {
					break
				}
				_, err = targetFile.Write(receiveBuf[:n])
				if err != nil {
					return err
				}
				remaining -= n
				callback.Progress(float64(int(fileSize)-remaining)/float64(fileSize)*100, fmt.Sprintf("第%d /%d个文件", i+1, numFiles))
			}
			targetFile.Close()
		}
	}
	callback.Progress(100, fmt.Sprintf("共%d个文件接收完毕", numFiles))
	if utils.IsAutoOpen() {
		coolOs.ShowFile(path)
	}
	return nil
}
