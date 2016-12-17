package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/vthommeret/memory.limited/lib/ml"
)

func main() {
	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatalf("Unable to stat stdin.")
	}

	var f io.Reader

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

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("Unable to read file: %s", err)
	}

	var p ml.Page
	err = json.Unmarshal(b, &p)
	if err != nil {
		log.Fatalf("Unable to unmarshal JSON: %s", err)
	}

	w := ml.Parse(p)

	fmt.Printf("%s\n\n", w.Value)
	for _, s := range w.Sections {
		fmt.Printf("%s%s\n", strings.Repeat("  ", s.Depth-1), s.Name)
	}
}
