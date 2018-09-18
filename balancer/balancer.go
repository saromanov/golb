package balancer

import (
	"github.com/saromanov/golb/server"
)

// Balancer defines basic interface for
// load balancing algorithms
type Balancer interface {
	Do() (*server.Server, error)
}
