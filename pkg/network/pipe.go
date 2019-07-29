package network

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

type Pipe struct {
	ClientConn net.Conn
	DockerConn net.Conn
}

// Pipe creates a full-duplex pipe between the two sockets and transfers data from one to the other.
func (p *Pipe) StartPiping() {
	go p.pipeItFrom(p.ClientConn, p.DockerConn)
	go p.pipeItFrom(p.DockerConn, p.ClientConn)
}

func (p *Pipe) pipeItFrom(conn net.Conn, conn2 net.Conn)  {
	buf := bufio.NewReader(conn)
	byteArray := make([]byte, 128)

	for {
		n, err := buf.Read(byteArray)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Some err?:", err)
			}
			p.DockerConn.Close()
			p.ClientConn.Close()
			break
		}

		conn2.Write(byteArray[:n])
	}
}
