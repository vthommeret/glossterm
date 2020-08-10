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
	"strings"

	"github.com/vthommeret/glossterm/lib/gt"
	"github.com/vthommeret/glossterm/lib/radix"

	"github.com/cayleygraph/cayley"
)

const sourceLang = "fr"

const defaultWordsPath = "data/words.gob"
const defaultIndexPath = "data/index.gob"
const defaultGraphPath = "data/graph.db"
const defaultPort = 8080

const max = 10

var wordsPath string
var indexPath string
var graphPath string
var port int

var words map[string]*gt.Word
var index *radix.Tree
var graph *cayley.Handle

func init() {
	flag.StringVar(&wordsPath, "w", defaultWordsPath, "Words path (gob format)")
	flag.StringVar(&indexPath, "i", defaultIndexPath, "Index path (gob format)")
	flag.StringVar(&graphPath, "g", defaultGraphPath, "Graph path (boltdb format)")
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
	log.Printf("Loading words from %q.", wordsPath)
	ws, err := gt.GetWords(wordsPath)
	if err != nil {
		log.Fatalf("Unable to get words: %s", err)
	}
	words = ws

	// Get index
	log.Printf("Loading index from %q.", indexPath)
	t, err := gt.GetIndex(indexPath)
	if err != nil {
		log.Fatalf("Unable to get radix tree: %s", err)
	}
	index = t

	// Get graph
	log.Printf("Loading graph from %q.", graphPath)
	g, err := gt.GetGraph(graphPath)
	if err != nil {
		log.Fatalf("Unable to get graph: %s", err)
	}
	graph = g

	// Setup handlers
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/search", searchHandler)

	// Listen
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
	if index == nil {
		http.Error(w, "Index not yet indexed.", http.StatusInternalServerError)
		return
	}

	q := strings.TrimSpace(r.URL.Query().Get("query"))

	if word, ok := words[q]; ok && maybeWriteDescendants(w, word) {
		return
	}

	rs := index.FindWordsWithPrefix(strings.ToLower(q), max)
	if len(rs) > max {
		rs = rs[:max]
	} else if len(rs) == 1 {
		eq := gt.Normalize(string(rs[0])) == gt.Normalize(q)
		if eq && maybeWriteDescendants(w, words[string(rs[0])]) {
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(map[string]interface{}{
		"type":    "results",
		"results": rs,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func maybeWriteDescendants(w http.ResponseWriter, word *gt.Word) bool {
	descendants := latinDescendants(word)
	if descendants != nil {
		w.Header().Set("Content-Type", "application/json")
		b, err := json.Marshal(map[string]interface{}{
			"type":        "descendants",
			"descendants": descendants,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return true
		}
		w.Write(b)
		return true
	}
	return false
}

type Descendant struct {
	Lang string
	Word string
}

// Latin descendants of French words.
func latinDescendants(w *gt.Word) []Descendant {
	ds := gt.GetDescendants(graph, sourceLang, w.Name)
	return toDescendants(ds)
}

func toDescendants(is []interface{}) (ds []Descendant) {
	for _, i := range is {
		lang, word := idParts(i)
		ds = append(ds, Descendant{
			Lang: lang,
			Word: word,
		})
	}
	return ds
}

func idParts(id interface{}) (string, string) {
	parts := strings.Split(id.(string), "/")
	return parts[0], parts[1]
}
