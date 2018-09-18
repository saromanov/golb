package balancer

import (
	"sort"
	"sync"

	"github.com/saromanov/golb/golb"
)

type LeastConnect struct {
	*sync.Mutex
	Connections uint32
	Servers     Servers
	serverNum   uint32
}

// Do provides sorting of the servers by active connections
func (rr *LeastConnect) Do() (*golb.Server, error) {
	sort.Sort(rr.Servers)
	return rr.Servers[0], nil
}
