package ml

import (
	"encoding/xml"
	"fmt"
	"io"
)

type Redirect struct {
	Title string `xml:"title,attr"`
}

type Page struct {
	Title string   `xml:"title"`
	Redir Redirect `xml:"redirect"`
	Text  string   `xml:"revision>text"`
}

const count = 10

func ParseXML(r io.Reader, word string) (Page, error) {
	d := xml.NewDecoder(r)

	i := 0
	for {
		t, err := d.Token()
		if err != nil {
			return Page{}, fmt.Errorf("unable to decode token: %s", err)
		}
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "page" {
				var p Page
				d.DecodeElement(&p, &se)

				if p.Title == word {
					return p, nil
				}

				i++
				if i == count {
					//break Loop
				}
			}
		}
	}
	return Page{}, fmt.Errorf("unable to parse XML")
}
