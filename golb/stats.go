package golb

// Stats defines structure for the golb stats
type Stats struct {
	Requests         uint32 	`json:"requests"`
	StatusCodes      map[int]uint32 `json:"status_codes"`
	CompleteRequests uint32 `json:"complete_requests"`
	Servers          uint32 `json:"servers"`
}
