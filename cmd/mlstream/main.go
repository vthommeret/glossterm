package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

const count = 10

type Redirect struct {
	Title string `xml:"title,attr"`
}

type Page struct {
	Title string   `xml:"title"`
	Redir Redirect `xml:"redirect"`
	Text  string   `xml:"revision>text"`
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Must specify file and word.")
	}

	fp := os.Args[1]
	f, err := os.Open(fp)
	if err != nil {
		log.Fatalf("Unable to open fp: %s", err)
	}
	d := xml.NewDecoder(f)

	w := os.Args[2]

	i := 0
Loop:
	for {
		t, err := d.Token()
		if err != nil {
			log.Fatalf("Unable to decode token: %s", err)
		}
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "page" {
				var p Page
				d.DecodeElement(&p, &se)

				if p.Title == w {
					fmt.Printf("%+s\n", p.Text)
					break Loop
				}

				i++
				if i == count {
					//break Loop
				}
			}
		}
	}
}
