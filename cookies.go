package bot

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

type CookieJar struct {
	Data map[string][]*http.Cookie
}

func (bot *Bot) Cookies() ([]byte, error) {
	history := bot.History().Entries()
	jar := &CookieJar{
		Data: make(map[string][]*http.Cookie),
	}
	for i := range history {
		u, err := url.Parse(history[i])
		if err != nil {
			log.Printf("Invalid URL in history! Skipping cookies from it: " + history[i])
			continue
		}
		jar.Data[u.String()] = bot.j.Cookies(u)
	}
	return json.MarshalIndent(jar, "", "  ")
}

func (bot *Bot) SetCookies(cookies []byte) error {
	jar := &CookieJar{
		Data: make(map[string][]*http.Cookie),
	}
	if err := json.Unmarshal(cookies, jar); err != nil {
		return err
	}

	for k := range jar.Data {
		v := jar.Data[k]
		u, err := url.Parse(k)
		if err != nil {
			log.Printf("Invalid URL from cookies! Skipping: " + k)
			continue
		}
		bot.j.SetCookies(u, v)
		bot.History().Add(u.String())
	}
	return nil
}
