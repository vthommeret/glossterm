package gt

import (
	"fmt"
	"log"
	"strings"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph/path"
	"github.com/cayleygraph/quad"
)

type Cognate struct {
	Word string `json:"word" firestore:"word"`
	From string `json:"from" firestore:"from"`
}

func GetCognates(graph *cayley.Handle, lang string, word string) []*Cognate {
	prefix := fmt.Sprintf("%s/", lang)
	w := quad.String(prefix + word)

	s := cayley.StartPath(graph, w)
	ps := findParents(s).Tag("parent")

	// Find children of parent or second-degree parent
	p := findChildren(findParents(ps).Tag("parent")).
		Or(findChildren(ps))

	rs, ts, err := QueryGraph(graph, p)
	if err != nil {
		log.Fatalf("Unable to execute query: %s", err)
	}

	ds := []*Cognate{}

	for i, r := range rs {
		if strings.HasPrefix(r, prefix) {
			continue
		}
		ds = append(ds, &Cognate{Word: r, From: ts[i]})
	}

	return ds
}

func findParents(p *path.Path) *path.Path {
	return p.
		Out("borrowing-from").
		Or(p.Out("derived-from")).
		Or(p.Out("inherited-from")).
		Or(p.Out("mentions")).
		Or(p.Out("suffix"))
}

func findChildren(p *path.Path) *path.Path {
	return p.
		Out("descendant").
		Or(p.Out("cognate"))
}

/*

// Gremlin query

var input = g.V("fr/pelouse")

function findParents(g) {
  return g.out("borrowing-from")
    .or(g.out("derived-from"))
    .or(g.out("inherited-from"))
    .or(g.out("mentions"))
    .or(g.out("suffix"))
}

var parents = findParents(input)

findParents(parents).out("descendant")
  .or(parents.out("descendant"))
  .all()

*/
