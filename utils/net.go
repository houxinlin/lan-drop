package utils

import (
	"fmt"
	"net"
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

func GetOutBoundIP() string {
	netInterfaceAddresses, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, netInterfaceAddress := range netInterfaceAddresses {
		networkIp, ok := netInterfaceAddress.(*net.IPNet)
		if ok && !networkIp.IP.IsLoopback() && networkIp.IP.To4() != nil {
			ip := networkIp.IP.String()
			return ip
		}
	}
	return ""
}
