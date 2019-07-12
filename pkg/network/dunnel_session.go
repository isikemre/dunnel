package network

import "github.com/google/uuid"

type DunnelSession struct {

	Id string

	WebSocketHub WebSocketHub

}

func newDunnelSession(id string) (DunnelSession) {
	if id == "" {
		id = uuid.New().String()
	}
	return DunnelSession{ Id: id, WebSocketHub: WebSocketHub{} }
}