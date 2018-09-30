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

// http://kb.linuxvirtualserver.org/wiki/Weighted_Round-Robin_Scheduling
func (rr *WeightedRoundRobin) getServer() *server.Server {
	i := 0
	var cw int32
	numServers := len(rr.Servers)
	for {
		i = (i + 1) % numServers
		if i == 0 {
			cw = gcd(rr.Servers[i].Weight, rr.Servers[i+1].Weight)
			if cw <= 0 {
				cw = getMaxWeight(rr.Servers)
				if cw == 0 {
					return nil
				}
			}
		}
		if rr.Servers[i].Weight >= cw {
			return rr.Servers[i]
		}
	}
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

func gcd(x, y int32) int32 {
	var n int32
	for {
		n = x % y
		if n <= 0 {
			return y
		}
		x = y
		y = n
	}
}
