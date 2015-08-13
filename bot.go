package bot

import (
	"fmt"
	"net/url"
	"net/http"
	"net/http/cookiejar"
)

type Bot struct {
	b string
	c http.Client
	d bool
}

func New() (*Bot) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		// Currently, cookiejar.Nil never returns an error
		panic(err)
	}

	return &Bot{
		c: http.Client {
			Transport: &transport{
				t: http.DefaultTransport},
			Jar: jar,
		},
	}
}

func (bot *Bot) Get(url string) (*Page, error) {
	resp, err := bot.c.Get(bot.b + url)
	if err != nil {
		return nil, err
	}
	return &Page{resp: resp}, nil
}

func (bot *Bot) Post(url string, form url.Values) (*Page, error) {
	resp, err := bot.c.PostForm(bot.b + url, form)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("bot: non 2xx response code: %d: %s", resp.StatusCode, resp.Status)
	}
	return &Page{resp: resp}, nil
}

func (bot *Bot) Do(r *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("bot.Do: not implemented")
}

// Debug enables debugging messages to standard error stream.
func (bot *Bot) Debug(enabled bool) (*Bot) {
	bot.d = enabled
	return bot
}

func (bot *Bot) SetUA(userAgent string) (*Bot) {
	bot.c.Transport.(*transport).ua = userAgent
	return bot
}

func (bot *Bot) BaseURL(baseURL string) (*Bot) {
	bot.b = baseURL
	return bot
}