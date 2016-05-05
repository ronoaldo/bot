// Copyright 2015 Ronoaldo JLP <ronoaldo@gmail.com>
// Licensed under the Apache License, Version 2.0

package bot

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var (
	// ErrTooManyRedirects is returned when the bot reaches more than 10 redirects.
	ErrTooManyRedirects = errors.New("bot: too many redirects")
)

// Bot implements a statefull HTTP client for interacting with websites.
type Bot struct {
	b     string
	c     *http.Client
	debug bool

	// lastURL records the last seen URL using the CheckRedirect function.
	// TODO(ronoaldo): change to a history of recent URLs.
	history *History
}

// New initializes a new Bot with an in-memory cookie management.
func New() *Bot {
	return ReuseClient(&http.Client{})
}

func ReuseClient(c *http.Client) *Bot {
	jar, err := cookiejar.New(nil)
	if err != nil {
		// Currently, cookiejar.Nil never returns an error
		panic(err)
	}
	c.Jar = jar
	bot := &Bot{
		c: c,
		history: &History{},
	}
	t := &transport{
		t: http.DefaultTransport,
		b: bot,
	}
	bot.c.Transport = t
	bot.c.CheckRedirect = bot.checkRedirect
	return bot
}

// GET performs the HTTP GET to the provided URL and returns a Page.
// It returns a nil page if there is a network error.
// It will also return an error if the response is not 2xx,
// but the returned page is non-nil, and you can parse the error body.
func (bot *Bot) GET(url string) (*Page, error) {
	bot.history.Add(bot.b + url)
	resp, err := bot.c.Get(bot.b + url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("bot: non 2xx response code: %d: %s", resp.StatusCode, resp.Status)
	}
	return &Page{resp: resp}, nil
}

// POST performs an HTTP POST to the provided URL,
// using the form as a payload, and returns a Page.
// It returns a nil page if there is a network error.
// It will also return an error if the response is not 2xx,
// but the returned page is non-nil, and you can parse the error body.
func (bot *Bot) POST(url string, form url.Values) (*Page, error) {
	bot.history.Add(bot.b + url)
	resp, err := bot.c.PostForm(bot.b+url, form)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("bot: non 2xx response code: %d: %s", resp.StatusCode, resp.Status)
	}
	return &Page{resp: resp}, nil
}

// Debug enables debugging messages to standard error stream.
func (bot *Bot) Debug(enabled bool) *Bot {
	bot.debug = enabled
	return bot
}

// SetUA allows one to change the default user agent used by the Bot.
func (bot *Bot) SetUA(userAgent string) *Bot {
	bot.c.Transport.(*transport).ua = userAgent
	return bot
}

// BaseURL can be used to setup Bot base URL,
// that will then be a prefix used by Get and Post methods.
func (bot *Bot) BaseURL(baseURL string) *Bot {
	bot.b = baseURL
	return bot
}

func (bot *Bot) History() *History {
	return bot.history
}

func (bot *Bot) checkRedirect(req *http.Request, via []*http.Request) error {
	log.Printf("Redirecting to: %v (via %v)", req, via)
	bot.history.Add(req.URL.String())
	if len(via) > 10 {
		return ErrTooManyRedirects
	}
	return nil
}
