package bot

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"strings"
)

type Form struct {
	Method string
	ID string
	Name string
	Action string
	Fields url.Values
}

type Table struct {
	ID     string
	Class  string
	Header []string
	Data   [][]string
	RawCells [][]string
}

type Page struct {
	resp *http.Response
}

// Raw returns the raw bytes from the response as an io.Reader
func (page *Page) Raw() *http.Response {
	return page.resp
}

func (page *Page) Tables() ([]Table) {
	if page == nil {
		return nil
	}
	if page.resp == nil {
		return nil
	}
	doc, err := goquery.NewDocumentFromReader(page.resp.Body)
	if err != nil {
		panic(err)
	}
	debugf("Loaded document from response.")
	tables := make([]Table, 0)
	doc.Find("table").Each(func(i int, t *goquery.Selection) {
		table := Table{
			ID:    t.AttrOr("id", ""),
			Class: t.AttrOr("class", ""),
		}

		t.Find("tr").Each(func(j int, tr *goquery.Selection){
			rawRow := make([]string, 0)
			tr.Find("th").Each(func (k int, th *goquery.Selection) {
				// Append to the header
				table.Header = append(table.Header, th.Text())
				if raw, err := th.Html(); err == nil {
					rawRow = append(rawRow, raw)
				} else {
					debugf("Error parsing node HTML: %v", err)
				}
			})
			row := make([]string, 0, 10)
			tr.Find("td").Each(func (k int, td *goquery.Selection) {
				row = append(row, td.Text())
				if raw, err := td.Html() ; err == nil {
					rawRow = append(rawRow, raw)
				} else {
					debugf("Error parsing node HTML: %v", err)
				}
			})
			if len(row) > 0 {
				table.Data = append(table.Data, row)
			}
			table.RawCells = append(table.RawCells, rawRow)
		})

		tables = append(tables, table)
	})

	return tables
}

// Forms parses the page extracting all forms found as url.Values.
func (page *Page) Forms() ([]Form) {
	// create a new document from response body
	if page == nil {
		return nil
	}
	if page.resp == nil {
		return nil
	}
	doc, err := goquery.NewDocumentFromReader(page.resp.Body)
	if err != nil {
		panic(err)
	}
	debugf("Loaded document from response: %#v", doc)

	forms := make([]Form, 0)
	// Parse the forms in the document
	doc.Find("form").Each(func(i int, f *goquery.Selection) {
		formid := f.AttrOr("id", "")
		action := f.AttrOr("action", "")
		method := f.AttrOr("method", "GET")
		name := f.AttrOr("name", "")
		debugf("Found new form[id=%s, action=%s, method=%s]", formid, action, method)
		fields := make(url.Values)

		// Parse all input fields
		f.Find("input").Each(func(j int, input *goquery.Selection) {
			_type := input.AttrOr("type", "")
			name := input.AttrOr("name", "")
			value := input.AttrOr("value", "")
			debugf("> Parsing input[type=%s, name=%s, value=%s]", _type, name, value)
			switch strings.ToLower(_type) {
			case "text", "hidden":
				fields[name] = append(fields[name], value)
			case "radio":
				// We should only store the selected radio value
				if chk := input.AttrOr("checked", "unchecked"); chk != "unchecked" {
					fields[name] = []string{value}
				}
			case "submit":
				// If this is a submit button, we only store it if the name is not empty
				if name != "" {
					fields[name] = []string{value}
				}
			}
		})

		// Parse all select fields
		f.Find("select").Each(func(j int, _select *goquery.Selection) {
			name := _select.AttrOr("name", "")
			debugf("Processando select[name=%s]", name)
			if name != "" {
				_select.Find("option").Each(func(k int, option *goquery.Selection) {
					if option.AttrOr("selected", "notselected") != "notselected" {
						if value, has := option.Attr("value"); has {
							fields[name] = append(fields[name], value)
						} else {
							fields[name] = append(fields[name], strings.TrimSpace(option.Text()))
						}
					}
				})
			}
		})

		forms = append(forms, Form{
			ID:     formid,
			Method: method,
			Action: action,
			Name:   name,
			Fields: fields,
		})
	})
	return forms
}