package gt

import (
	"fmt"
	"log"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph/path"
	"github.com/cayleygraph/quad"
)

func GetDescendants(graph *cayley.Handle, lang string, word string) []interface{} {
	w := quad.String(fmt.Sprintf("%s/%s", lang, word))

	s := cayley.StartPath(graph, w)
	ps := findParents(s)
	p := findParents(ps).Out(quad.String("descendant")).
		Or(ps.Out(quad.String("descendant")))

	rs, err := QueryGraph(graph, p)
	if err != nil {
		log.Fatalf("Unable to execute query: %s", err)
	}

	return rs
}

func findParents(p *path.Path) *path.Path {
	return p.Out(quad.String("borrowing-from")).
		Or(p.Out("derived-from")).
		Or(p.Out("inherited-from")).
		Or(p.Out("mentions")).
		Or(p.Out("suffix"))
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
