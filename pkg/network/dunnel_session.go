package network

import "github.com/google/uuid"

type DunnelSession struct {

	Id string

}

func newDunnelSession(id string) (DunnelSession) {
	if id == "" {
		id = uuid.New().String()
	}
	return DunnelSession{ Id: id }
}