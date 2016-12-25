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
	_ "github.com/blevesearch/bleve/analysis/analyzers/simple_analyzer"
)

const defaultIndexPath = "words.bleve"
const defaultPort = 8080

var indexPath string
var port int

var index bleve.Index

func init() {
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

	bleveIndex, err := bleve.Open(indexPath)
	if err != nil {
		log.Fatalf("Unable to open %q index: %s", indexPath, err)
	}
	index = bleveIndex

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
		"success": true,
		"results": searchResults,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
