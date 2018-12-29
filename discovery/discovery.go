package discovery

import "github.com/saromanov/golb/server"

// Discovery defines interface for several ways for discovery
type Discovery interface {
	Search() error
	GetServers() []*server.Server
}
