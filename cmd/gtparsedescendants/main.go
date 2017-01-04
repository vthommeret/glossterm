package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"

	"github.com/vthommeret/glossterm/lib/gt"
	"github.com/vthommeret/glossterm/lib/lang"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Must specify file.")
	}
	fp := os.Args[1]
	f, err := os.Open(fp)
	if err != nil {
		log.Fatalf("Unable to open fp: %s", err)
	}

	d := xml.NewDecoder(f)

	var p gt.Page
	err = d.Decode(&p)
	if err != nil {
		log.Fatalf("Unable to unmarshal JSON: %s", err)
	}

	w, err := gt.ParseDescendants(p, lang.DefaultLangMap)
	if err != nil {
		log.Fatalf("Unable to parse word: %s", err)
	}

	b, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		log.Fatalf("Unable to marshal JSON: %s", err)
	}

	fmt.Println(string(b))
}
