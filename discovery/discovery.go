package discovery

// Discovery defines interface for several ways for discovery
type Discovery interface {
	Search() error
	Stop()
}
