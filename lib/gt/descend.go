package gt

import (
	"fmt"
	"log"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph/path"
	"github.com/cayleygraph/quad"
)

type Descendant struct {
	Word string
	From string
}

func GetDescendants(graph *cayley.Handle, lang string, word string) []Descendant {
	w := quad.String(fmt.Sprintf("%s/%s", lang, word))

	s := cayley.StartPath(graph, w)
	ps := findParents(s).Tag("parent")

	// Find children of parent or second-degree parent
	p := findChildren(findParents(ps).Tag("parent")).
		Or(findChildren(ps))

	rs, ts, err := QueryGraph(graph, p)
	if err != nil {
		log.Fatalf("Unable to execute query: %s", err)
	}

	ds := []Descendant{}

	for i, r := range rs {
		ds = append(ds, Descendant{Word: r, From: ts[i]})
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
