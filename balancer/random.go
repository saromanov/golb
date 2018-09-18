package balancer

import (
	"math/rand"
	"sync"

	"github.com/saromanov/golb/server"
)

type Random struct {
	*sync.Mutex
	Connections uint32
	Servers     server.Servers
	serverNum   uint32
}

// Do provides sorting of the servers by active connections
func (rr *Random) Do() (*server.Server, error) {
	n := rand.Intn(len(rr.Servers))
	return rr.Servers[n], nil
}
