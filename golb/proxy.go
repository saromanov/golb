package golb

import (
	"fmt"
	"net/http"

	"github.com/saromanov/golb/server"
)

// HTTPProxy defines main struct for the proxy
type HTTPProxy struct {
	serv *server.Server
}

// Do provides executing of the proxy
func (p *HTTPProxy) Do(w http.ResponseWriter, r *http.Request) error {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, v []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	r.Host = p.serv.Host
	r.RequestURI = ""
	resp, err := client.Do(r)
	if err != nil {
		return err
	}

	fmt.Println(resp)
	return nil
}
