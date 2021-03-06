package golb

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/saromanov/golb/server"
)

// HTTPProxy defines main struct for the proxy
type HTTPProxy struct {
	serv    *server.Server
	Scheme  string
	Headers map[string]string
}

// HTTPProxyResponse returns data from response
type HTTPProxyResponse struct {
	statusCode int
	duration   time.Duration
	request    string
	date       time.Time
}

// Do provides executing of the proxy
func (p *HTTPProxy) Do(w http.ResponseWriter, r *http.Request) (*HTTPProxyResponse, error) {
	u, err := url.Parse(fmt.Sprintf("%s://%s:%d", p.Scheme, p.serv.Host, p.serv.Port) + r.RequestURI)
	if err != nil {
		return nil, urlParseError{err: err}
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, v []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	r.Header.Set("X-Forwarded-Host", p.serv.Host)
	if len(p.Headers) > 0 {
		for key, value := range p.Headers {
			r.Header.Set(key, value)
		}
	}
	startTime := time.Now()
	r.Host = p.serv.Host
	r.URL = u
	r.RequestURI = ""
	resp, err := client.Do(r)
	if err != nil {
		return nil, httpRequestError{err: err, req: r}
	}

	response := &HTTPProxyResponse{
		statusCode: resp.StatusCode,
		duration:   time.Since(startTime),
		request:    u.Path,
		date:       time.Now().UTC(),
	}
	defer resp.Body.Close()
	return response, nil
}
