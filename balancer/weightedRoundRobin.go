package balancer

import (
	"sync"

	"github.com/saromanov/golb/server"
)

// WeightedRoundRobin defines round robin with weights
type WeightedRoundRobin struct {
	sync.Mutex
	Connections uint32
	Servers     server.Servers
	serverNum   uint32
}

func (rr *WeightedRoundRobin) Do() (*server.Server, error) {
	return rr.getServer(), nil
}

func (rr *WeightedRoundRobin) getServer() *server.Server {
	i := -1
	numServers := len(rr.Servers)
	for {
		i = (i + 1) % numServers
		cw := gcd(rr.Servers)
		if cw <= 0 {
			cw = getMaxWeight(rr.Servers)
			if cw == 0 {
				return nil
			}
		}
		if rr.Servers[i].Weight >= cw {
			return rr.Servers[i]
		}
	}
	return nil
}

// getMaxWeight return max weight from all servers
func getMaxWeight(servers server.Servers) int32 {

	var maxWeight int32
	for _, s := range servers {
		if s.Weight > maxWeight {
			maxWeight = s.Weight
		}
	}
	return maxWeight
}

func gcd(servers server.Servers) int32 {
	return 1
}
