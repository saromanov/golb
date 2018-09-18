package balancer

import "github.com/saromanov/golb/golb"

// Servers provides definition of the list of servers
type Servers []*golb.Server

func (s Servers) Len() int { return len(s) }

func (s Servers) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Servers) Less(i, j int) bool {
	return s[i].GetActiveConnections() < s[j].GetActiveConnections()
}

// Balancer defines basic interface for
// load balancing algorithms
type Balancer interface {
	Do() (*golb.Server, error)
}
