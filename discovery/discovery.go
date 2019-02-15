package discovery

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/saromanov/golb/server"
)

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

// GenID provides generation of ID for server
func GenID(s *server.Server) string {
	hasher := md5.New()
	port := fmt.Sprintf("%d", s.Port)
	hasher.Write([]byte(fmt.Sprintf("%s%s", s.Host, port))) // nolint
	return hex.EncodeToString(hasher.Sum(nil))
}
