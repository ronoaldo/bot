package bot

import (
	"fmt"
	"log"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
	"strconv"
)

func TestBotCookieJar(t *testing.T) {
	s := httptest.NewServer(&TestServer{})
	defer s.Close()

	bot := New()
	p, err := bot.Get(s.URL + "/private/")
	if err != nil {
		t.Error(err)
	}
	resp := p.Raw()
	defer resp.Body.Close()
	checkStatus(t, "for unauthorized response", resp.StatusCode, 403)

	// Login and create a cookie
	p, err = bot.Post(s.URL + "/login/", make(url.Values))
	if err != nil {
		t.Error(err)
	}

	resp = p.Raw()
	defer resp.Body.Close()
	checkStatus(t, "for login", resp.StatusCode, 200)
	checkBody(t, resp.Body, "OK")

	// We should be able to see private data
	p, err = bot.Get(s.URL + "/private/")
	if err != nil {
		t.Error(err)
	}

	resp = p.Raw()
	defer resp.Body.Close()
	checkStatus(t, "for authorized request", resp.StatusCode, 200)
	checkBody(t, resp.Body, "PRIVATE")
}

func checkStatus(t *testing.T, when string, got, expected int) {
	if got != expected {
		t.Errorf("Unexpected status code %s: %d, expected %d", when, got, expected)
	}
}

func checkBody(t *testing.T, b io.Reader, expected string) {
	if body, err := ioutil.ReadAll(b); err != nil {
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
			Name: "TSID",
			Value: t.session,
			Path: "/",
		})
		fmt.Fprintf(w, "OK")
	case "/private/":
		var (
			c *http.Cookie
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
