package gt

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/quad"

	_ "github.com/cayleygraph/cayley/graph/kv/bolt"
	"github.com/cayleygraph/cayley/graph/path"
)

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
