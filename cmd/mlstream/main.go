package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/vthommeret/memory.limited/lib/ml"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Must specify file and word.")
	}

	fp := os.Args[1]
	f, err := os.Open(fp)
	if err != nil {
		log.Fatalf("Unable to open fp: %s", err)
	}

	w := os.Args[2]

	p, err := ml.ParseXML(f, w)
	if err != nil {
		log.Fatalf("Unable to parse XML: %s", err)
	}

	b, err := json.Marshal(p)
	if err != nil {
		log.Fatalf("Unable to marshal json: %s", err)
	}
	fmt.Print(string(b))
}
