package network

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

func handleWebSocket(w http.ResponseWriter, r *http.Request, session *DunnelSession) {
	h, ok := w.(http.Hijacker)
	if !ok {
		//return u.returnError(w, r, http.StatusInternalServerError, "websocket: response does not implement http.Hijacker")
		return
	}

	conn, _, err := h.Hijack()

	if err != nil {
		fmt.Println("Error while hijacking connection", err)
	}

	conn.SetDeadline(time.Time{})

	dockerConn, err := net.Dial("unix", "/var/run/docker.sock")

	if err != nil {
		fmt.Println("Error while dialing to unix socket '/var/run/docker.sock'", err)
	}

	session.WebSocketHub.ClientConn = conn
	session.WebSocketHub.DockerConn = dockerConn

	session.WebSocketHub.StartPiping()
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
