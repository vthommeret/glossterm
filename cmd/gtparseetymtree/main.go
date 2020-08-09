package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"vthommeret/glossterm/lib/gt"
	"vthommeret/glossterm/lib/lang"
)

func main() {
	// Return file info for stdin.
	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatalf("Unable to stat stdin.")
	}

	var f io.Reader

	// Set f to stdin or specified file.
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		f = os.Stdin
	} else {
		if len(os.Args) < 2 {
			log.Fatalf("Must specify file.")
		}
		fp := os.Args[1]
		f, err = os.Open(fp)
		if err != nil {
			log.Fatalf("Unable to open fp: %s", err)
		}
	}

	d := xml.NewDecoder(f)

	var p gt.Page
	err = d.Decode(&p)
	if err != nil {
		log.Fatalf("Unable to unmarshal JSON: %s", err)
	}

	w, err := gt.ParseEtymTree(p, lang.DefaultLangMap)
	if err != nil {
		log.Fatalf("Unable to parse word: %s", err)
	}

	b, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		log.Fatalf("Unable to marshal JSON: %s", err)
	}

	fmt.Println(string(b))
}
