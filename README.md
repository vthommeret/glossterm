# glossterm

## Commands

In order to generate files for the web app, you need to grab an English
Wiktionary dump, put it in `data/` and run the following commands.

1. `gtsplit`   splits Wiktionary dump into N files so it can be parsed
               in parallel. N is set to the current number of cores.

2. `gtstream`  parses split files into words.gob and descendants.gob.

3. `gtresolve` reads words.gob and looks up DescendantTrees references
               in descendants.gob, and inlines them.

4. `gtindex`   indexes terms for each words to power autocomplete.

5. `gtquads`   generates quads for each word to power graph lookups,
               e.g. find all descendants for the Latin roots of a
               given word.

Once you've run those commands, you can run `gtweb` which will launch
the web app.

## Mobile

In order to generate a Mobile framework for iOS/Android, you can run `gomobile bind -target=ios github.com/vthommeret/glossterm/lib/mobile`.

More information is available here: https://github.com/golang/go/wiki/Mobile
