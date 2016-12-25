package main

import (
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/blevesearch/bleve"
	_ "github.com/blevesearch/bleve/analysis/analyzers/simple_analyzer"
	"github.com/vthommeret/memory.limited/lib/ml"
	"github.com/vthommeret/memory.limited/lib/tpl"
)

const defaultWordsPath = "words.gob"
const defaultIndexPath = "words.bleve"
const defaultPort = 8080

var wordsPath string
var indexPath string
var port int

var words map[string]*ml.Word
var index bleve.Index

func init() {
	flag.StringVar(&wordsPath, "w", defaultWordsPath, "Words path (gob format)")
	flag.StringVar(&indexPath, "i", defaultIndexPath, "Index path (bleve format)")

	flag.IntVar(&port, "p", defaultPort, "Port (default 8080)")
	flag.Parse()
}

func main() {
	if envPort := os.Getenv("PORT"); envPort != "" {
		intPort, err := strconv.Atoi(envPort)
		if err != nil {
			log.Fatalf("Unable to convert %q port: %s", envPort, err)
		}
		port = intPort
	}

	f, err := os.Open(wordsPath)
	if err != nil {
		log.Fatalf("Unable to open %q words: %s", wordsPath, err)
	}

	dec := gob.NewDecoder(f)

	var decoded map[string]*ml.Word
	err = dec.Decode(&decoded)
	if err != nil {
		log.Fatalf("Unable to decode gob: %s", err)
	}
	f.Close()

	words = decoded

	bleveIndex, err := bleve.Open(indexPath)
	if err != nil {
		log.Fatalf("Unable to open %q index: %s", indexPath, err)
	}
	index = bleveIndex
	defer index.Close()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/search", searchHandler)

	log.Printf("Listening on port %d.", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("assets/tpl/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")

	if word, ok := words[q]; ok {
		from, descendants := latinDescendants(word)
		if from != nil && descendants != nil {
			w.Header().Set("Content-Type", "application/json")
			b, err := json.Marshal(map[string]interface{}{
				"type":        "descendants",
				"from":        from,
				"descendants": descendants,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(b)
			return
		}
	}

	query := bleve.NewPrefixQuery(q)

	search := bleve.NewSearchRequest(query)
	search.Highlight = bleve.NewHighlight()

	searchResults, err := index.Search(search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(map[string]interface{}{
		"type":    "results",
		"results": searchResults,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

// Latin descendants of Spanish words.
func latinDescendants(w *ml.Word) (from *tpl.Mention, descendants []tpl.Link) {
Loop:
	for _, l := range w.Languages {
		if l.Code == "es" && l.Etymology.Mentions != nil {
			for _, m := range l.Etymology.Mentions {
				if m.Lang == "la" {
					from = &m
					break Loop
				}
			}
		}
	}
	if from == nil {
		return nil, nil
	}
	if w, ok := words[from.Word]; ok {
		for _, l := range w.Languages {
			if l.Code == "la" && l.Descendants != nil {
				for _, d := range l.Descendants {
					if d.Lang != "es" {
						descendants = append(descendants, d)
					}
				}
			}
		}
	}
	return from, descendants
}
