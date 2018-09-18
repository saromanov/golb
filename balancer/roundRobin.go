package balancer

import (
	"sync"

	"github.com/saromanov/golb/server"
)

type RoundRobin struct {
	sync.Mutex
	Connections uint32
	Servers     server.Servers
	serverNum   uint32
}

func (rr *RoundRobin) Do() (*server.Server, error) {
	idx := rr.serverNum
	rr.updateNum()
	return rr.Servers[idx], nil
}

func (rr *RoundRobin) updateNum() {
	rr.Lock()
	defer rr.Unlock()
	if rr.serverNum == uint32(len(rr.Servers)-1) {
		rr.serverNum = 0
		return
	}
	rr.serverNum++
}
