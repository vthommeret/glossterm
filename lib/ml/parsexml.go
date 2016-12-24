package ml

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

type Redirect struct {
	Title string `xml:"title,attr"`
}

type Page struct {
	Title string   `xml:"title"`
	Redir Redirect `xml:"redirect"`
	Text  string   `xml:"revision>text"`
}

type Error struct {
	Message string
}

const count = 1

func ParseXMLWord(r io.Reader, w string) (*Page, error) {
	d := xml.NewDecoder(r)
	for {
		t, err := d.Token()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "page" {
				var p Page
				d.DecodeElement(&p, &se)
				if p.Title == w {
					return &p, nil
				}
			}
		}
	}
	return nil, nil
}

func ParseXML(r io.Reader, pages chan<- Page, errors chan<- Error, done chan<- bool) {
	d := xml.NewDecoder(r)

Parse:
	for {
		t, err := d.Token()
		if err != nil {
			if err != io.EOF {
				errors <- Error{fmt.Sprintf("unable to decode token: %s", err)}
			}
			break
		}
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "page" {
				var p Page
				d.DecodeElement(&p, &se)

				// Exclude namespaced pages.
				if strings.Contains(p.Title, ":") {
					continue Parse
				}

				pages <- p
			}
		}
	}

	done <- true
}
