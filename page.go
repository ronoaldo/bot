// Copyright 2015 Ronoaldo JLP <ronoaldo@gmail.com>
// Licensed under the Apache License, Version 2.0

package bot

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	errNilPage     = fmt.Errorf("bot: nil page")
	errNilResp     = fmt.Errorf("bot: nil response")
	errNilRespBody = fmt.Errorf("bot: nil response body")
)

// Form is a representation of an HTML form structure.
// This struct is used by Page to parse HTML forms embedded into the document.
// Fields is populated with the parsed form input and select fields.
type Form struct {
	Method string
	ID     string
	Name   string
	Action string

	// Fields contains the parsed form input and select elements.
	Fields url.Values
}

// Print pretty prints the form into a human-readable, line delimited string.
func (f *Form) Print() string {
	buff := new(bytes.Buffer)
	fmt.Fprintf(buff, "form#%s\n", f.ID)
	maxw := 0
	keys := make([]string, 0, len(f.Fields))
	for k := range f.Fields {
		if len(k) > maxw {
			maxw = len(k)
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	format := fmt.Sprintf("%%%ds:%%s\n", maxw)
	for _, k := range keys {
		if len(f.Fields[k]) > 0 {
			fmt.Fprintf(buff, format, k, f.Fields[k])
		}
	}
	return buff.String()
}

// Table represents an HTML data table.
// It is used by Page to store the parsed table data.
type Table struct {
	ID     string
	Class  string
	Header []string
	Data   [][]string

	// RawCells, unlike Data, contains the raw HTML elements inside
	// table cells.
	RawCells [][]string
}

// Page is a wrapper to an http.Response, with some usefull methods.
type Page struct {
	resp *http.Response
	body []byte
}

// Raw returns the raw Response, after reading all the data from the response Body.
// This is usefull to inspect the response size, headers, and other attributes.
// The body is already closed after you call this method, and if you need the raw
// response bytes, call Body() instead.
func (page *Page) Raw() (*http.Response, error) {
	// Load the body in memory before returning.
	if _, err := page.Body(); err != nil {
		return nil, err
	}
	return page.resp, nil
}

// Body returns a copy of the response body as a new Reader.
// Use this if you need to integrate with any third party that expectes the reader.
// Thi is, basically, a shortcut for bytes.NewReader(page.Body()).
func (page *Page) Body() (io.Reader, error) {
	b, err := page.Bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

// Bytes returns the page body bytes, so you can json.Unmarshal it.
func (page *Page) Bytes() ([]byte, error) {
	if err := page.sanityCheck(); err != nil {
		return nil, err
	}
	if err := page.ensureBodyReady(); err != nil {
		return nil, err
	}
	return page.body, nil
}

// Tables parses the response body, and extract all <table>s from it.
// The result is nil, if there is an error reading the response,
// or if there is an error building the document reader.
//
// The returned table is filled with the text content from the table.
// That usually means that Table.Data is the text in the table cells.
// <th> elements are stored in the Table.Header slice.
// You can find the raw HTML data from each table cell
// in the Table.RawCells. This is usefull if you need to parse links inside
// tables.
func (page *Page) Tables() ([]Table, error) {
	var (
		body io.Reader
		doc  *goquery.Document
		err  error
	)
	if body, err = page.Body(); err != nil {
		return nil, err
	}
	if doc, err = goquery.NewDocumentFromReader(body); err != nil {
		return nil, err
	}
	debugf("Loaded document from response.")
	var tables []Table
	doc.Find("table").Each(func(i int, t *goquery.Selection) {
		table := Table{
			ID:    t.AttrOr("id", ""),
			Class: t.AttrOr("class", ""),
		}

		t.Find("tr").Each(func(j int, tr *goquery.Selection) {
			var rawRow []string
			tr.Find("th").Each(func(k int, th *goquery.Selection) {
				// Append to the header
				table.Header = append(table.Header, th.Text())
				if raw, err := th.Html(); err == nil {
					rawRow = append(rawRow, raw)
				} else {
					debugf("Error parsing node HTML: %v", err)
				}
			})
			row := make([]string, 0, 10)
			tr.Find("td").Each(func(k int, td *goquery.Selection) {
				row = append(row, td.Text())
				if raw, err := td.Html(); err == nil {
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

	return tables, nil
}

// Forms parses the response extracting all <form> elements
// and returns them in the resulting Form slice.
// The return is nil if there is an error reading from the response,
// or if there is an error building the document reader.
//
// From the form tag, this method returns the properties:
// 	* id
// 	* action
// 	* method
// 	* name
//
// The parser scans all <input> and <select> elements inside the form,
// and decodes their values into the Form.Fields map.
// For selects, the returned value is the option marked with the "selected"
// attribute, or empty otherwise.
func (page *Page) Forms() ([]Form, error) {
	var (
		body io.Reader
		doc  *goquery.Document
		err  error
	)
	if body, err = page.Body(); err != nil {
		return nil, err
	}
	if doc, err = goquery.NewDocumentFromReader(body); err != nil {
		return nil, err
	}
	debugf("Loaded document from response: %#v", doc)

	var forms []Form
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
			case "text", "hidden", "":
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
	return forms, nil
}

// sanityCheck makes sure that the page is valid, and is wrapping a valid response.
func (page *Page) sanityCheck() error {
	if page == nil {
		return errNilPage
	}
	if page.resp == nil {
		return errNilResp
	}
	if page.resp.Body == nil {
		return errNilRespBody
	}
	return nil
}

// ensureBodyReady makes sure that the body is read once from the response.
func (page *Page) ensureBodyReady() error {
	if page.body == nil {
		var err error
		page.body, err = ioutil.ReadAll(page.resp.Body)
		if err != nil {
			return err
		}
		defer page.resp.Body.Close()
		// Let's do some magic here, to convert ISO-8859-1 (Latin1) pages to Unicode
		ct := strings.ToLower(page.resp.Header.Get("Content-Type"))
		if strings.Contains(ct, "charset=iso-8859-1") {
			b := make([]rune, len(page.body))
			for i, c := range page.body {
				b[i] = rune(c)
			}
			page.body = []byte(string(b))
		}
	}
	return nil
}
