package golb

import (
	"fmt"
	"net/http"
)

// urlParseError retruns error in the case
// of url parsing
type urlParseError struct {
	err error
}

func (e urlParseError) Error() string {
	return fmt.Sprintf("unable to parse url: %v", e.err)
}

type httpRequestError struct {
	err error
	req *http.Request
}

func (e httpRequestError) Error() string {
	return fmt.Sprintf("unable to make http request: %s %v", e.req.Host, e.err)
}
