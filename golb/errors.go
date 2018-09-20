package golb

import "fmt"

type urlParseError struct {
	err error
}

// Error retruns message of the openStateError
func (e *urlParseError) Error() string {
	return fmt.Sprintf("unable to parse url: %v", e.err)
}
