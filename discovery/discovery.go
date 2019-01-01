package discovery

import "github.com/saromanov/golb/server"

// Discovery defines interface for several ways for discovery
type Discovery interface {
	Search() error
	GetServers() []*server.Server
}

// Config provides definition of discovery config
type Config struct {
	// Filters provides definition of restrictions
	// for find services
	Filters string

	// DockerEndpoint provides definition for endpoint
	DockerEndpoint string
}
