// Copyright 2015 Ronoaldo JLP <ronoaldo@gmail.com

/*
Package bot implements a simple HTTP automation utility.


Motivation

The web today is widelly used by several companies,
but not all web sites are implemented equal.
Besides the variety of technologies that can be used,
some apps simply does not offer a clean, REST or SOAP like
interface of interaction.

Bot is a simple tool that allows developers to interact with
legacy websites and webservices, that are not yet in the new age.
It is a statefull HTTP client, meaning that it uses a Cookie Jar
implementation, allowing you to automate tasks that would usually
be done manually.


Getting Started

Bot is designed to allow for you to embedd it in your own bots.
One primary use case is to allow you to submit an authentication form,
then issue other requests with the same session cookies
that were set by the authentication request.

Bot is highly integrated with GoQuery,
allowing you to quicly process the HTML responses.
Here is a sample snippet that logins into a service,
and updates a form.

	b := bot.New().BaseURL("http://website")
	page, err := b.POST("/login", url.Values{ "username": {"myuser"}, "password": {"not-a-secret"}})
	if err != nil {
		// If the request is not in the 2xx range, it is an error
		log.Fatal(err)
	}
	for _, f := range page.Forms() {
		if f.ID == "preferences" {
			f.Fields["mailings"] = []string{""}
			if _, err := b.POST("/account/", f.Fields) {
				log.Fatal(err)
			}
		}
	}

Bot does not provide any JavaScript inmplementation, nor a Rendering engine.
It is just a headless, statefull HTTP client.
*/
package bot // import "ronoaldo.gopkg.net/bot"
