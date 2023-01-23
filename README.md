# glossterm 0.5

glossterm is a pipeline that extracts, lexes, and parses wiktionary data.

## Pipeline

In order to generate files for the web app, you need to grab an English
Wiktionary dump, put it in `data/` and run the following commands.

You can run the commands by doing e.g. `go run cmd/gtdump/main.go` or running `make`
to install globally available commands that can be run as e.g. `gtdump`.

1. `gtdump`
   downloads Wiktionary dump to en.xml.bz2.

1. `gtsplit`
   splits Wiktionary dump into N files so it can be parsed in parallel.
   N is set to the current number of cores.

1. `gtparse`
   parses split files into words.gob and descendants.gob.
   Use --no-backup after initial change to index to edit index in place and
   compare to previously committed index.

1. `gtresolve`
   reads words.gob and looks up DescendantTrees references in
   descendants.gob, and inlines them.

1. `gtquads`
   generates quads for each word to power graph lookups, e.g. find all
   descendants for the Latin roots of a given word.

1. `gtbeam`
   fetches cognates in parallel using Apache Beam local runner.

1. `gtcognates`
   inlines cognates from `gtbeam` into words.gob

1. `gtcompare`
   compares new index to old index. always use to manually verify parsing changes

1. `gtindex`
   incrementally indexes (additions, deletions, updates) words in Firestore

## Debugging a single word

1. `gtpage <word>`
   extracts a single XML page for a given word.
   Example: `gtpage helado`

1. `gtlex <word.xml>`
   lexes a single XML page for a given word.
   Example: `gtpage hombre | gtlex`

1. `gtparseword <word.xml>`
   parses a single XML word.
   Example: `gtpage horno | gtparseword`

1. `gtparseetymtree <word.xml>`
   parses a single etymtree XML page.
   Example: `gtpage Template:etymtree/la/germanus | gtparseetymtree`

1. `gtdescend <word>`
   shows the descendants from any words mentioned for a given word.

1. `gtread <word>`
   reads word from words.gob.
   Example: `gtpage pt/nariz`

1. `gtsearch <query>`
   searches the index for a given word.
