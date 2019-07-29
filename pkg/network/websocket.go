package network

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

func handleWebSocket(wr http.ResponseWriter, r *http.Request, session *DunnelSession) {
	dockerConn, err := net.Dial("unix", "/var/run/docker.sock")
	if err != nil {
		panic(err)
	}

	requestString := fmt.Sprintf("%s %s HTTP/1.1\nHost: %s\n", r.Method, r.URL, r.Host)

	for name, value := range r.Header {
		for _, headerValue := range value {
			requestString += fmt.Sprintf("%s: %s\n", name, headerValue)
		}
	}

	// End HTTP Header
	requestString += "\n\n"

	if r.Method != "GET" && r.Method != "HEAD" {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		requestString += string(bodyBytes)
		requestString += "\n"
	}

	fmt.Printf("%s", requestString)

	clientConn, _ := hijack(wr)
	pipe := Pipe{DockerConn: dockerConn, ClientConn: clientConn}
	pipe.StartPiping()
	dockerConn.Write([]byte(requestString))
	//clientWriter.Write([]byte("HTTP/1.1 101 Switching Protocols\nConnection: Upgrade\nContent-Type: application/vnd.docker.raw-stream\nUpgrade: tcp\nDate: Fri, 19 Jul 2019 14:46:44 GMT"))
	//clientWriter.Flush()
}

func isWebSocket(r *http.Request) bool {
	if strings.ToLower(r.Header.Get("Connection")) != "upgrade" {
		return false
	}
	if strings.ToLower(r.Header.Get("Upgrade")) != "tcp" && strings.ToLower(r.Header.Get("Upgrade")) != "websocket" {
		return false
	}
	fmt.Println("This connection is a WEBSOCKET")
	return true
}


func hijack(w http.ResponseWriter) (net.Conn, *bufio.ReadWriter) {
	hj, _ := w.(http.Hijacker)
	conn, buf, _ := hj.Hijack()
	return conn, buf
}
