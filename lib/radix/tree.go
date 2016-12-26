package radix

import (
	"fmt"
	"strings"
)

// Copied from https://github.com/gojp/nihongo/blob/master/lib/dictionary/radix_tree.go

type EntryID string

type Tree struct {
	Root *Node
}

func (r Tree) String() string {
	return fmt.Sprintf("%v", r.Root.Edges)
}

func NewTree() *Tree {
	root := Node{
		Edges: []Edge{},
		Ids:   nil,
	}
	return &Tree{Root: &root}
}

type Edge struct {
	Target *Node
	Label  string
}

func (r Edge) String() string {
	return r.Label
}

type Node struct {
	Edges []Edge
	Ids   []EntryID
}

func (n Node) IsLeaf() bool {
	return len(n.Ids) > 0
}

func (n Node) Value() []EntryID {
	return n.Ids
}

func (n Node) FindPrefixedEntries(max int) (entries []EntryID) {
	entries = []EntryID{}

	stack := []*Node{&n}
	added := map[EntryID]bool{}

	var node *Node
	for len(stack) > 0 {
		node, stack = stack[len(stack)-1], stack[:len(stack)-1]
		if node.IsLeaf() {
			for _, v := range node.Value() {
				if _, ok := added[v]; !ok {
					entries = append(entries, v)
				}
				added[v] = true
			}
		}
		if len(entries) >= max {
			return
		}
		for i := range node.Edges {
			stack = append(stack, node.Edges[i].Target)
		}
	}
	return
}

func (r *Tree) findLastMatchingNode(key string) (n *Node, elementsFound int) {
	n = r.Root
	elementsFound = 0

	for n != nil && elementsFound < len(key) {
		// Get the next edge to explore based on the elements not yet found in key
		// select edge from n.Edges where edge.Label is a prefix of key.suffix(elementsFound)
		suffix := key[elementsFound:]
		var nextEdge *Edge
		for i := range n.Edges {
			if strings.HasPrefix(suffix, n.Edges[i].Label) {
				nextEdge = &n.Edges[i]
				break
			}
		}

		if nextEdge == nil {
			// terminate loop
			break
		}

		// Was an edge found?
		// Set the next node to explore
		n = nextEdge.Target

		// Increment elements found based on the label stored at the edge
		elementsFound += len(nextEdge.Label)
	}

	return n, elementsFound
}

func (r *Tree) Insert(key string, id EntryID) {
	n, elementsFound := r.findLastMatchingNode(key)
	if n == nil || elementsFound > len(key) {
		return
	}
	if elementsFound == len(key) {
		// key already exists, so add id to this node
		if n.Ids == nil {
			n.Ids = []EntryID{}
		}
		n.Ids = append(n.Ids, id)
	} else {
		// check if an outgoing edge shares a prefix with us
		suffix := key[elementsFound:]
		prefix := ""
		sharedEdge := -1

		for i := range n.Edges {
		inner:
			for u := 0; u < len(n.Edges[i].Label) && u < len(suffix); u++ {
				if n.Edges[i].Label[u] == suffix[u] {
					prefix += suffix[u : u+1]
				} else {
					break inner
				}
			}
			// there can be at most one outgoing edge that shares a prefix
			if len(prefix) > 0 {
				sharedEdge = i
				break
			}
		}

		if sharedEdge == -1 {
			// create a new edge and node
			n.Edges = append(n.Edges, Edge{Target: &Node{Edges: []Edge{}, Ids: []EntryID{id}}, Label: suffix})
		} else {
			oldEdge := n.Edges[sharedEdge]

			child := Node{Ids: []EntryID{}}
			n.Edges[sharedEdge] = Edge{Target: &child, Label: prefix}

			node := Node{Edges: []Edge{}, Ids: []EntryID{id}}
			var left Edge
			if prefix == suffix {
				child.Ids = append(child.Ids, id)
			} else {
				left = Edge{Target: &node, Label: suffix[len(prefix):]}
				child.Edges = []Edge{left}
			}
			right := Edge{Target: oldEdge.Target, Label: oldEdge.Label[len(prefix):]}
			child.Edges = append(child.Edges, right)
		}
	}
}

func (r *Tree) Get(key string) []EntryID {
	n, elementsFound := r.findLastMatchingNode(key)

	// A match is found if we arrive at a leaf node and have used up exactly len(key) elements
	if n != nil && elementsFound == len(key) {
		return n.Value()
	}

	return nil
}

func (r *Tree) FindWordsWithPrefix(key string, max int) []EntryID {
	words := []EntryID{}

	n, elementsFound := r.findLastMatchingNode(key)
	if n != nil {
		if elementsFound == len(key) {
			children := n.FindPrefixedEntries(max - len(words))
			words = append(words, children...)
		} else {
			// check if an outgoing edge shares a prefix with us
			suffix := key[elementsFound:]
			prefix := ""
			sharedEdge := -1

			for i := range n.Edges {
			inner:
				for u := 0; u < len(n.Edges[i].Label) && u < len(suffix); u++ {
					if n.Edges[i].Label[u] == suffix[u] {
						prefix += suffix[u : u+1]
					} else {
						break inner
					}
				}
				// there can be at most one outgoing edge that shares a prefix
				if len(prefix) > 0 {
					sharedEdge = i
					break
				}
			}

			if sharedEdge >= 0 {
				children := n.Edges[sharedEdge].Target.FindPrefixedEntries(max - len(words))
				words = append(words, children...)
			}
		}
	}

	return words
}
