package golb

// Stats defines structure for the golb stats
type Stats struct {
	Requests         uint32
	StatusCodes      map[uint32]uint32
	CompleteRequests uint32
	Servers          uint32
}
