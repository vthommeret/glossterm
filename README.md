# glossterm

## Pipeline

In order to generate files for the web app, you need to grab an English
Wiktionary dump, put it in `data/` and run the following commands.

1. `gtdump`    downloads Wiktionary dump to en.xml.bz2.

1. `gtsplit`   splits Wiktionary dump into N files so it can be parsed
               in parallel. N is set to the current number of cores.

1. `gtparse`   parses split files into words.gob and descendants.gob.

1. `gtresolve` reads words.gob and looks up DescendantTrees references
               in descendants.gob, and inlines them.

1. `gtindex`   indexes terms for each words to power autocomplete.

1. `gtquads`   generates quads for each word to power graph lookups,
               e.g. find all descendants for the Latin roots of a
               given word.

Once you've run those commands, you can run `gtweb` which will launch
the web app.

## Debugging a single word

1. `gtpage <word>`                 extracts a single XML page for a given
                                   word.

1. `gtlex <word.xml>`              lexes a single XML page for a given word.

1. `gtparseword <word.xml>`        parses a single XML page for a given word.

1. `gtparsedescendants <word.xml>` parses the descendants for a single
                                   XML page for a given word.

1. `gtdescend <word>`              shows the descendants from any words
                                   mentioned for a given word.

## Additional commands

1. `gtread`           reads words.gob.

1. `gtsearch <query>` searches the index for a given word.

## Dependencies

### Updating dependencies

`go install`

### Adding dependencies

`glide get github.com/dustin/go-humanize`
