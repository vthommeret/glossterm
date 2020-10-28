package gt

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/quad"
	"github.com/cayleygraph/quad/nquads"

	_ "github.com/cayleygraph/cayley/graph/kv/bolt"
	"github.com/cayleygraph/cayley/graph/path"
)

// GetGraphNquads returns graph from nquads file
func GetGraphNquads(path string) (*cayley.Handle, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Unable to open %s input: %s", path, err)
	}

	r := bufio.NewReader(f)
	nr := nquads.NewReader(r, false)

	store, err := cayley.NewMemoryGraph()
	if err != nil {
		log.Fatalf("Unable to create memory store: %s", err)
	}

	for {
		q, err := nr.ReadQuad()
		if err != nil {
			if err != io.EOF {
				log.Fatalf("Unable to read quad: %s", err)
			}
			break
		}
		if !q.IsValid() {
			fmt.Printf("Got invalid quad: %+v\n", q)
		}
		store.AddQuad(q)
	}

	return store, nil
}

// GetGraph returns graph either from path or compressed path.
func GetGraph(path string) (*cayley.Handle, error) {
	if !exists(path) {
		compressed := fmt.Sprintf("%s.gz", path)
		log.Printf("Uncompressing %q.", compressed)
		cf, err := os.Open(compressed)
		defer cf.Close()
		if err != nil {
			return nil, err
		}
		gr, err := gzip.NewReader(cf)
		defer gr.Close()
		if err != nil {
			return nil, err
		}
		tmp, err := ioutil.TempFile("", "words-graph")
		tmpName := tmp.Name()
		defer os.Remove(tmpName)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(tmp, gr)
		if err != nil {
			log.Fatalf("Unable to decompress db: %s", err)
		}
		path = tmpName
	}
	graph.InitQuadStore("bolt", path, nil)
	return cayley.NewGraph("bolt", path, nil)
}

// QueryGraph queries graph.
func QueryGraph(g *cayley.Handle, p *path.Path) (rs []string, ts []string, err error) {
	err = p.Iterate(nil).EachValue(nil, func(v quad.Value) {
		rs = append(rs, quad.NativeOf(v).(string))
	})
	if err != nil {
		return nil, nil, err
	}

	err = p.Iterate(nil).TagValues(nil, func(m map[string]quad.Value) {
		if v, ok := m["parent"]; ok {
			ts = append(ts, quad.NativeOf(v).(string))
		}
	})
	if err != nil {
		return nil, nil, err
	}

	return rs, ts, nil
}
