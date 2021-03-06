package server

// Server defines struct for the server definition
type Server struct {
	ID 				  string
	Host              string
	Port              uint32
	ActiveConnections uint32
	Requests          uint32

	// After n failed requests,
	// server removing from the list
	FailedRequests  uint32
	SuccessRequests uint32
	RemovedFromList bool
	Weight          int32
}

// Servers provides definition of the list of servers
type Servers []*Server

func (s Servers) Len() int { return len(s) }

func (s Servers) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Servers) Less(i, j int) bool {
	return s[i].GetActiveConnections() < s[j].GetActiveConnections()
}

// GetActiveConnections return current number of
// active connections on the server
func (s *Server) GetActiveConnections() uint32 {
	return s.ActiveConnections
}

// StartServer provides starting of the server
func (s *Server) StartServer() error {
	return nil
}

// IncRequests provides increment of requests
func (s *Server) IncSuccessRequests() {
	s.Requests++
	s.SuccessRequests++
}
