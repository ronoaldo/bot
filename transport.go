// Copyright 2015 Ronoaldo JLP <ronoaldo@gmail.com>
// Licensed under the Apache License, Version 2.0

package bot

import (
	"log"
	"net/http"
	"net/http/httputil"
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
	b  *Bot
}

func (t *transport) userAgent() string {
	if t.ua == "" {
		return "Mozilla/5.0 (compatible)"
	}
	return t.ua
}

// RoundTrip implements the http.RoundTripper interface.
func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	if t.b.debug {
		b, _ := httputil.DumpRequest(r, true)
		log.Printf("> Dumped request: \n>>>\n%s\n>>>\n", string(b))
	}
	req := &request{Request: r}
	req.setUserAgent(t.userAgent())
	resp, err := t.t.RoundTrip(req.Request)
	if t.b.debug {
		if resp != nil {
			b, _ := httputil.DumpResponse(resp, false)
			log.Printf("Dumped response: \n<<<\n%s\n<<<\n", string(b))
		} else {
			log.Printf("Dumped response: \n<<<\n(nil): err=%v\n<<<\n", err)
		}
	}
	return resp, err
}
