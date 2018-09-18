package balancer

import (
	"sync"

	"github.com/saromanov/golb/golb"
)

type RoundRobin struct {
	*sync.Mutex
	Connections uint32
	Servers     Servers
	serverNum   uint32
}

func (rr *RoundRobin) Do() (*golb.Server, error) {
	idx := rr.serverNum
	rr.updateNum()
	return rr.Servers[idx], nil
}

func (rr *RoundRobin) updateNum() {
	rr.Lock()
	defer rr.Unlock()
	rr.serverNum++
}
