package network

type Tunnel struct {
	destinationIp   string
	destinationPort int
}

func CreateTunnel() *Tunnel {
	t := Tunnel{
		destinationIp: "192.168.1.1",
		destinationPort: 9011,
	}
	return &t
}