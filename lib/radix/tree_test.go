package radix

import "testing"

// Copied from https://github.com/gojp/nihongo/blob/master/lib/dictionary/radix_tree_test.go

var getTests = []string{
	"apple",
	"shoe",
	"tree",
	"banana",
	"band",
}

var notIncluded = []string{
	"orange",
	"car",
	"beach",
}

func TestGet(t *testing.T) {
	r := NewTree()
	for i, entry := range getTests {
		r.Insert(entry, EntryID(i))
		got := r.Get(entry)
		if len(got) != 1 {
			t.Fatalf("%q len(got) = %d, want %d", entry, len(got), 1)
		}
		if got[0] != EntryID(i) {
			t.Fatalf("got[0] = %q, want %q", got[0], entry)
		}
	}

	for _, entry := range notIncluded {
		got := r.Get(entry)
		if len(got) != 0 {
			t.Fatalf("%q len(got) = %d, want %d", entry, len(got), 0)
		}
	}
}

func TestFindWordsWithPrefix(t *testing.T) {
	r := NewTree()
	for i, entry := range getTests {
		r.Insert(entry, EntryID(i))
	}

	got := len(r.FindWordsWithPrefix("sho", 10))
	if got != 1 {
		t.Fatalf("%q len(got) = %d, want %d", "sho", got, 1)
	}

	got = len(r.FindWordsWithPrefix("ban", 10))
	if got != 2 {
		t.Fatalf("%q len(got) = %d, want %d", "ban", got, 2)
	}
}
