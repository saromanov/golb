package server

// Server defines struct for the server definition
type Server struct {
	Host              string
	ActiveConnections uint32
}

// GetActiveConnections return current number of
// active connections on the server
func (s *Server) GetActiveConnections() uint32 {
	return s.ActiveConnections
}
