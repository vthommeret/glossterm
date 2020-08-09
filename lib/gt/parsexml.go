package gt

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"vthommeret/glossterm/lib/lang"
)

type Redirect struct {
	Title string `xml:"title,attr"`
}

type Page struct {
	XMLName xml.Name
	Title   string   `xml:"title"`
	Redir   Redirect `xml:"redirect"`
	Text    string   `xml:"revision>text"`
}

type Error struct {
	Message string
	Fatal   bool
}

const count = 1
const etymTree = "Template:etymtree/"

// ParseXMLPage returns page for cmd/gtpage.
func ParseXMLPage(r io.ReadCloser, title string, page chan<- Page, errors chan<- Error, done chan<- io.ReadCloser) {
	d := xml.NewDecoder(r)
Parse:
	for {
		t, err := d.Token()
		if err != nil {
			if err != io.EOF {
				errors <- Error{fmt.Sprintf("unable to decode token: %s", err), true}
			}
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "page" {
				var p Page
				d.DecodeElement(&p, &se)
				if p.Title == title {
					page <- p
					break Parse
				}
			}
		}
	}
	done <- r
}

// ParseXMLPages returns pages for cmd/gtsplit.
func ParseXMLPages(r io.ReadCloser, pages chan<- Page, errors chan<- Error, done chan<- io.ReadCloser) {
	d := xml.NewDecoder(r)

Parse:
	for {
		t, err := d.Token()
		if err != nil {
			if err != io.EOF {
				errors <- Error{fmt.Sprintf("unable to decode token: %s", err), true}
			}
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "page" {
				var p Page
				d.DecodeElement(&p, &se)
				// Exclude namespaced pages.
				if strings.Contains(p.Title, ":") && !strings.HasPrefix(p.Title, etymTree) {
					continue Parse
				}
				pages <- p
			}
		}
	}

	done <- r
}

// ParseXMLWords returns words and descendants for cmd/gtstream.
func ParseXMLWords(r io.ReadCloser, words chan<- Word, descendants chan<- Descendants, errors chan<- Error, done chan<- io.ReadCloser) {
	d := xml.NewDecoder(r)

Parse:
	for {
		t, err := d.Token()
		if err != nil {
			if err != io.EOF {
				errors <- Error{fmt.Sprintf("unable to decode token: %s", err), true}
			}
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "page" {
				var p Page
				d.DecodeElement(&p, &se)
				if strings.HasPrefix(p.Title, etymTree) {
					ds, err := ParseEtymTree(p, lang.DefaultLangMap)
					if err != nil {
						errors <- Error{fmt.Sprintf("unable to parse %q word: %s", p.Title, err), false}
						continue Parse
					}
					descendants <- *ds
				} else {
					w, err := ParseWord(p, lang.DefaultLangMap)
					if err != nil {
						errors <- Error{fmt.Sprintf("unable to parse %q word: %s", p.Title, err), false}
						continue Parse
					} else if w.IsEmpty() {
						continue Parse
					}
					words <- w
				}
			}
		}
	}

	done <- r
}
