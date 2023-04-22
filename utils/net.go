package utils

import (
	"fmt"
	"net"
	"strings"
)

func FindAvailableUDPPort() (int, error) {
	ip := net.ParseIP("0.0.0.0").To4()
	addr := &net.UDPAddr{IP: ip, Port: 0}
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return -1, fmt.Errorf("Failed to listen: %v", err)
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).Port, nil
}

func FindAvailablePort() (int, error) {
	l, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 0})
	if err != nil {
		return 0, err
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port, nil
}

func GetOutBoundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}
