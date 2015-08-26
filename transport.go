package bot

import (
	"net/http"
)

// request is a http.Request wrapper to add some helper functions
type request struct {
	*http.Request
}

func (r *request) header() http.Header {
	if r.Request.Header == nil {
		r.Request.Header = make(http.Header)
	}
	return r.Request.Header
}

func (r *request) setUserAgent(ua string) {
	r.header().Set("User-Agent", ua)
}

// transport type implements http.RoundTripper in order to allow
// doing some magic in the Bot requests.
type transport struct {
	t  http.RoundTripper
	ua string
}

func (t *transport) userAgent() string {
	if t.ua == "" {
		return "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"
	}
	return t.ua
}

// RoundTrip implements the http.RoundTripper interface.
func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	req := &request{Request: r}
	req.setUserAgent(t.userAgent())
	return t.t.RoundTrip(req.Request)
}
