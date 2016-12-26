package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/blevesearch/bleve"
	"github.com/vthommeret/memory.limited/lib/ml"
	"github.com/vthommeret/memory.limited/lib/tpl"
)

const defaultWordsPath = "data/words.gob"
const defaultIndexPath = "data/words.bleve"
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

	// Get words
	ws, err := ml.GetWords(wordsPath)
	if err != nil {
		log.Fatalf("Unable to get words: %s", err)
	}
	words = ws

	// Get index
	go setupIndex(indexPath)

	// Setup handlers
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/search", searchHandler)

	// Listen
	log.Printf("Listening on port %d.", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func setupIndex(indexPath string) {
	i, err := ml.GetIndex(indexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		i, err = ml.CreateIndex(indexPath)
		if err != nil {
			log.Fatalf("Unable to create index: %s", err)
		}
		err := ml.Index(i, words)
		if err != nil {
			log.Fatalf("Unable to index words: %s", err)
		}
		index = i
	} else if err != nil {
		log.Fatalf("Unable to get index: %s", err)
	} else {
		index = i
	}
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
	if index == nil {
		http.Error(w, "Index not yet indexed.", http.StatusInternalServerError)
		return
	}

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

type From struct {
	Lang string
	Word string
}

// Latin descendants of Spanish words.
func latinDescendants(w *ml.Word) (from *From, descendants []tpl.Link) {
	var mention *tpl.Mention
	var derived *tpl.Derived

Loop:
	for _, l := range w.Languages {
		if l.Code == "es" {
			for _, m := range l.Etymology.Mentions {
				if m.Lang == "la" {
					mention = &m
					break Loop
				}
			}
			for _, d := range l.Etymology.Derived {
				if d.FromLang == "la" || d.FromLang == "LL" { // Latin or Late Latin.
					derived = &d
					break Loop
				}
			}
		}
	}
	var wordName string
	if mention != nil {
		wordName = mention.Word
		from = &From{mention.Lang, mention.Word}
	} else if derived != nil {
		wordName = derived.FromWord
		from = &From{derived.FromLang, derived.FromWord}
	} else {
		return nil, nil
	}
	if w, ok := words[wordName]; ok {
		for _, l := range w.Languages {
			if l.Code == "la" {
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
