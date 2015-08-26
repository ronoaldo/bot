package bot

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
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
	if page, err = bot.Get(s.URL + "/private/"); err != nil {
		t.Error(err)
	}
	if resp, err = page.Raw(); err != nil {
		t.Error(err)
	}
	checkStatus(t, "for unauthorized response", resp.StatusCode, 403)

	// Login and create a cookie
	if page, err = bot.Post(s.URL+"/login/", make(url.Values)); err != nil {
		t.Error(err)
	}

	if resp, err = page.Raw(); err != nil {
		t.Error(err)
	}
	checkStatus(t, "for login", resp.StatusCode, 200)
	checkBody(t, page, "OK")

	// We should be able to see private data
	if page, err = bot.Get(s.URL + "/private/"); err != nil {
		t.Error(err)
	}

	if resp, err = page.Raw(); err != nil {
		t.Error(err)
	}
	checkStatus(t, "for authorized request", resp.StatusCode, 200)
	checkBody(t, page, "PRIVATE")
}

func checkStatus(t *testing.T, when string, got, expected int) {
	if got != expected {
		t.Errorf("Unexpected status code %s: %d, expected %d", when, got, expected)
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
