package cool

import (
	"bytes"
	"cool-transmission/utils"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type IntCallback func(message string)

// SendTaskContent 发送文件上下文
type SendTaskContent struct {
	Progress     func(size int, current int, progress float64) //发送进度回调
	doneCallback func(data int)                                //发送完成回调
	failCallback func(err error)                               //发送失败回调
	FilePath     string
	Ip           string
	Port         int
}

func SendFileToTarget(content SendTaskContent) error {
	err := doSendTo(content)
	if err != nil {
		content.failCallback(err)
		return err
	}
	return nil
}

type progressWriter struct {
	progress  chan<- int64
	sentBytes *int64
}

func (pw progressWriter) Write(p []byte) (int, error) {
	n := len(p)
	*pw.sentBytes += int64(n)
	pw.progress <- *pw.sentBytes
	return n, nil
}
func doSendTo(content SendTaskContent) error {
	parentDir := filepath.Dir(content.FilePath)
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", content.Ip, content.Port))
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	file, err := os.Open(content.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		files, err := utils.GetAllFiles(content.FilePath)
		if err != nil {
			return err
		}
		var dirCount, fileCount int
		for _, f := range files {
			info, err := os.Stat(f)
			if err != nil {
				return err
			}
			if !info.IsDir() {
				fileCount++
			}
		}
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.BigEndian, uint32(fileCount+dirCount))
		conn.Write(buf.Bytes())

		curIndex := 0
		for _, f := range files {
			buf.Reset()
			info, err := os.Stat(f)
			if err != nil {
				return err
			}
			if info.IsDir() {
				err = binary.Write(buf, binary.BigEndian, byte(0))
			} else {
				err = binary.Write(buf, binary.BigEndian, byte(1))
			}
			if err != nil {
				return err
			}
			fileItemName := strings.TrimPrefix(f, parentDir)
			err = binary.Write(buf, binary.BigEndian, uint32(len(fileItemName)))
			if err != nil {
				return err
			}
			// 写入文件名
			_, err = buf.Write([]byte(fileItemName))
			if err != nil {
				return err
			}

			if !info.IsDir() {
				err = binary.Write(buf, binary.BigEndian, uint32(info.Size()))
				if err != nil {
					return err
				}
				conn.Write(buf.Bytes()) //发送头部信息
				fileItem, _ := os.Open(f)
				// 读取文件内容并发送
				curIndex += 1
				writeFile(len(files), curIndex, info, conn, fileItem, content)
			}
		}
		content.doneCallback(0)
	} else {
		// 写入文件数量
		buf := new(bytes.Buffer)

		err = binary.Write(buf, binary.BigEndian, uint32(1))
		if err != nil {
			return err
		}

		// 写入文件类型和文件名长度
		err = binary.Write(buf, binary.BigEndian, byte(1))
		if err != nil {
			return err
		}
		err = binary.Write(buf, binary.BigEndian, uint32(len(fileInfo.Name())))
		if err != nil {
			return err
		}

		// 写入文件名
		_, err = buf.Write([]byte(fileInfo.Name()))
		if err != nil {
			fmt.Println(err)
			return err
		}

		// 写入文件大小并发送文件内容
		err = binary.Write(buf, binary.BigEndian, uint32(fileInfo.Size()))
		if err != nil {
			fmt.Println(err)
			return err
		}
		conn.Write(buf.Bytes())
		err := writeFile(1, 1, fileInfo, conn, file, content)
		if err != nil {
			fmt.Println(err)
			return err
		}
		content.doneCallback(0)
	}
	return nil
}

func writeFile(total int, cur int, fileInfo os.FileInfo, conn *net.TCPConn, file *os.File, content SendTaskContent) error {
	bufLen := fileInfo.Size()
	var sentBytes int64
	progress := make(chan int64, 100)
	go func() {
		for p := range progress {
			value := float64(p) / float64(bufLen) * 100
			content.Progress(total, cur, value)
			if total == cur && value == 100 {
				content.doneCallback(0)
			}
		}
	}()
	_, err := io.Copy(io.MultiWriter(conn, progressWriter{progress, &sentBytes}), file)
	return err
}
