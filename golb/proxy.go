package golb

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/saromanov/golb/server"
)

// HTTPProxy defines main struct for the proxy
type HTTPProxy struct {
	serv *server.Server
}

// Do provides executing of the proxy
func (p *HTTPProxy) Do(w http.ResponseWriter, r *http.Request) error {
	u, err := url.Parse("http://" + p.serv.Host + r.RequestURI)
	if err != nil {
		return err
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, v []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	r.Header.Set("X-Forwarded-Host", "127.0.0.1")
	r.Host = "127.0.0.1"
	r.URL = u
	r.RequestURI = ""
	resp, err := client.Do(r)
	if err != nil {
		return err
	}

	fmt.Println(resp)
	return nil
}
