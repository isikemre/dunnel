package network

import (
	"net"
)

func StartTCPProxy() {
	//server := tcp_server.New("localhost:9999")
}

func handleTCP(conn net.Conn) {
	//unix := createSocketConnection("/var/run/docker.sock")
	//tcpReader(unix)
}

func createSocketConnection(socketPath string) net.Conn {
	c, err := net.Dial("unix", socketPath)
	if err != nil {
		panic(err)
	}
	return c
}
