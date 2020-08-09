package mobile

import (
	"fmt"
	"log"
	"strings"
	"vthommeret/glossterm/lib/gt"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/quad"
)

type Graph struct {
	graph *cayley.Handle
}

func NewGraph(path string) (*Graph, error) {
	g, err := gt.GetGraph(path)
	if err != nil {
		return &Graph{}, err
	}
	return &Graph{
		graph: g,
	}, nil
}

func (g *Graph) Query(q string) string {
	w := quad.String(fmt.Sprintf("es/%s", q))

	v := cayley.StartPath(g.graph, w)
	bs := v.Out(quad.String("borrowing-from"))
	ds := v.Out(quad.String("derived-from"))
	is := v.Out(quad.String("inherited-from"))
	ms := v.Out(quad.String("mentions"))
	p := bs.Or(ds).Or(is).Or(ms).Out(quad.String("descendant"))

	rs, err := gt.QueryGraph(g.graph, p)
	if err != nil {
		log.Fatalf("Unable to execute query: %s", err)
	}

	var words []string
	for _, r := range rs {
		words = append(words, r.(string))
	}

	return strings.Join(words, ", ")
}
