package bot_test

import (
	"fmt"

	"github.com/ronoaldo/bot"
)

func Example() {
	b := bot.New()
	page, err := b.GET("https://en.wikipedia.org/wiki/Go_(programming_language)")
	if err != nil {
		fmt.Printf("Error loading page: %v", err)
		return
	}

	// Some helper functions can make it easier to find forms and tables:
	forms, err := page.Forms()
	if err != nil {
		fmt.Printf("Error parsing forms: %v", err)
	}
	for _, form := range forms {
		if form.ID == "searchform" {
			fmt.Println("Search action:", form.Action)
			fmt.Println("Search fields:", form.Fields.Encode())
		}
	}

	// For better parsing, we can just use the page.Query method and use
	// jQuery like sintax (provided by GoQuery)
	html, err := page.Query()
	if err != nil {
		fmt.Printf("Error parsing document: %v", err)
	}

	h1 := html.Find("h1#firstHeading").First().Text()
	fmt.Println("Article Title:", h1)

	// Output:
	// Search action: /w/index.php
	// Search fields: search=&title=Special%3ASearch
	// Article Title: Go (programming language)
}
