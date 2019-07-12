package network

import "net"

type WebSocketHub struct {
	ClientConn net.Conn

	DockerConn net.Conn
}

// Pipe creates a full-duplex pipe between the two sockets and transfers data from one to the other.
func (wsh *WebSocketHub) StartPiping() {
	chan1 := wsh.chanFromConn(wsh.ClientConn)
	chan2 := wsh.chanFromConn(wsh.DockerConn)

	for {
		select {
		case b1 := <-chan1:
			if b1 == nil {
				return
			} else {
				wsh.ClientConn.Write(b1)
			}
		case b2 := <-chan2:
			if b2 == nil {
				return
			} else {
				wsh.DockerConn.Write(b2)
			}
		}
	}
}

// chanFromConn creates a channel from a Conn object, and sends everything it
//  Read()s from the socket to the channel.
func (wsh *WebSocketHub) chanFromConn(conn net.Conn) chan []byte {
	c := make(chan []byte)

	go func() {
		b := make([]byte, 1024)

		for {
			n, err := conn.Read(b)
			if n > 0 {
				res := make([]byte, n)
				// Copy the buffer so it doesn't get changed while read by the recipient.
				copy(res, b[:n])
				c <- res
			}
			if err != nil {
				c <- nil
				break
			}
		}
	}()

	return c
}
