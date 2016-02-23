// Copyright 2015 Ronoaldo JLP <ronoaldo@gmail.com>
// Licensed under the Apache License, Version 2.0

package bot

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestBotCookieJar(t *testing.T) {
	var (
		resp *http.Response
		page *Page
		err  error
	)
	s := httptest.NewServer(&TestServer{})
	defer s.Close()

	bot := New()
	page, err = bot.GET(s.URL + "/private/")

	if err == nil {
		t.Errorf("Expected forbidden error, got nil")
	} else if !strings.Contains(err.Error(), "403") {
		t.Errorf("Expected 403 forbidden in error message, got %s", err.Error())
	}

	// Login and create a cookie
	if page, err = bot.POST(s.URL+"/login/", make(url.Values)); err != nil {
		t.Errorf("Error performing test authentication: %v", err)
	}

	if resp, err = page.Raw(); err != nil {
		t.Errorf("Unable to fetch raw page: %v", err)
		return
	}
	checkStatus(t, "for login", resp, 200)
	checkBody(t, page, "OK")

	// We should be able to see private data
	if page, err = bot.GET(s.URL + "/private/"); err != nil {
		t.Error(err)
	}

	if resp, err = page.Raw(); err != nil {
		t.Error(err)
	}
	checkStatus(t, "for authorized request", resp, 200)
	checkBody(t, page, "PRIVATE")
}

func checkStatus(t *testing.T, when string, resp *http.Response, expected int) {
	if resp == nil {
		t.Errorf("Response is nil")
		return
	}
	if resp.StatusCode != expected {
		t.Errorf("Unexpected status code %s: %d, expected %d", when, resp.StatusCode, expected)
	}
}

func checkBody(t *testing.T, page *Page, expected string) {
	if body, err := page.Bytes(); err != nil {
		t.Error(err)
	} else {
		// Body should be "OK"
		if string(body) != expected {
			t.Errorf("Unexpected body, expected '%s', got '%s'", expected, string(body))
		}
	}
}

// TestServer implements an http.Hander that handle /login/ and /private/
type TestServer struct {
	session string
}

func (t *TestServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/login/":
		t.session = strconv.FormatInt(time.Now().Unix(), 16)
		http.SetCookie(w, &http.Cookie{
			Name:  "TSID",
			Value: t.session,
			Path:  "/",
		})
		fmt.Fprintf(w, "OK")
	case "/private/":
		var (
			c   *http.Cookie
			err error
		)
		if c, err = r.Cookie("TSID"); err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if t.session != c.Value {
			log.Printf("E: Invalid session value %s -> expected %s", c.Value, t.session)
			// User is not logged in
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		fmt.Fprintf(w, "PRIVATE")
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}
