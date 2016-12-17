package ml

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"
)

// Adapted from https://github.com/golang/go/blob/master/src/text/template/parse/lex.go
// Also see https://www.youtube.com/watch?v=HxaD_trXwRE

// From parse.go

// Pos represents a byte position in the original input text from which
// this template was parsed.
type Pos int

// --------

// item represents a token or text string returned from the scanner.
type item struct {
	typ   itemType // The type of this item.
	pos   Pos      // The starting position, in bytes, of this item in the input string.
	val   string   // The value of this item.
	depth int      // Used by itemHeaderStart and itemHeaderEnd to indicate header depth.
}

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	case len(i.val) > 10:
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

// itemType identifies the type of lex items.
type itemType int

const (
	itemError itemType = iota // error occurred; value is text of error
	itemEOF
	itemHeaderStart // Header start
	itemHeaderEnd   // Header end
	itemText        // plain text
)

// Make the types prettyprint.
var itemName = map[itemType]string{
	itemError:       "error",
	itemEOF:         "EOF",
	itemHeaderStart: "header start",
	itemHeaderEnd:   "header end",
	itemText:        "text",
}

func (i itemType) String() string {
	s := itemName[i]
	if s == "" {
		return fmt.Sprintf("item%d", int(i))
	}
	return s
}

const eof = -1

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	input      string    // the string being scanned
	state      stateFn   // the next lexing function to enter
	pos        Pos       // current position in the input
	start      Pos       // start position of this item
	width      Pos       // width of last rune read from input
	lastPos    Pos       // position of most recent item returned by nextItem
	items      chan item // channel of scanned items
	parenDepth int       // nesting depth of ( ) exprs
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.start, l.input[l.start:l.pos], 0}
	l.start = l.pos
}

// emitHeader passes a header item back to the client.
func (l *lexer) emitHeader(t itemType, depth int) {
	l.items <- item{t, l.start, l.input[l.start:l.pos], depth}
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

// lineNumber reports which line we're on, based on the position of
// the previous item returned by nextItem. Doing it this way
// means we don't have to worry about peek double counting.
func (l *lexer) lineNumber() int {
	return 1 + strings.Count(l.input[:l.lastPos], "\n")
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, l.start, fmt.Sprintf(format, args...), 0}
	return nil
}

// NextItem returns the next item from the input.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) NextItem() item {
	item := <-l.items
	l.lastPos = item.pos
	return item
}

// Drain drains the output so the lexing goroutine will exit.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) Drain() {
	for range l.items {
	}
}

// NewLexer creates a new scanner for the input string.
func NewLexer(input string) *lexer {
	l := &lexer{
		input: normalize(input),
		items: make(chan item),
	}
	go l.run()
	return l
}

// Print prints all items
func (l *lexer) Print() {
	for {
		i := l.NextItem()
		if i.typ == itemEOF || i.typ == itemError {
			if i.typ == itemError {
				fmt.Printf("Error: %s\n", i.val)
			}
			break
		}
		if i.typ == itemHeaderStart || i.typ == itemHeaderEnd {
			fmt.Printf("%s, %d: %q\n", i.typ, i.depth, i.val)
		} else {
			fmt.Printf("%s: %q\n", i.typ, i.val)
		}
	}
}

// PrettyPrint pretty prints all items
func (l *lexer) PrettyPrint() {
	inHeader := false
	depth := 0
	for {
		i := l.NextItem()
		if i.typ == itemEOF || i.typ == itemError {
			if i.typ == itemError {
				fmt.Printf("Error: %s\n", i.val)
			}
			break
		}

		if i.typ == itemHeaderStart {
			inHeader = true
			depth = i.depth
		} else if i.typ == itemHeaderEnd {
			inHeader = false
		}

		if i.typ == itemText {
			prefix := ""
			if depth > 0 {
				prefix = strings.Repeat("  ", depth-1)
			}
			if inHeader {
				headerPrefix := strings.Repeat("#", depth)
				fmt.Printf("%s%s %s\n\n", prefix, headerPrefix, i.val)
			} else {
				fmt.Printf("%s", indent(i.val, prefix))
			}
		}
	}
}

// Simplifies parsing.
func normalize(s string) string {
	return fmt.Sprintf("\n%s\n", strings.TrimSpace(s))
}

func indent(s, prefix string) string {
	var buffer bytes.Buffer
	for _, line := range strings.Split(s, "\n") {
		buffer.WriteString(prefix + line + "\n")
	}
	return buffer.String()
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for l.state = lexText; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.items)
}

// state functions

const (
	headerDelim = "="
	headers     = 6
)

// lexText scans until an opening action delimiter, "{{".
func lexText(l *lexer) stateFn {
	for {
		for i := headers; i > 0; i-- {
			delim := strings.Repeat(headerDelim, i)

			leftDelim := "\n" + delim
			if strings.HasPrefix(l.input[l.pos:], leftDelim) {
				if l.pos > l.start {
					l.emit(itemText)
				}
				l.ignore()
				return lexHeaderDelim(itemHeaderStart, leftDelim, i)
			}

			rightDelim := delim + "\n"
			if strings.HasPrefix(l.input[l.pos:], rightDelim) {
				if l.pos > l.start {
					l.emit(itemText)
				}
				l.ignore()
				return lexHeaderDelim(itemHeaderEnd, rightDelim, i)
			}
		}
		if l.next() == eof {
			break
		}
	}
	// Correctly reached EOF.
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.emit(itemEOF)
	return nil
}

func lexHeaderDelim(it itemType, delim string, depth int) stateFn {
	return func(l *lexer) stateFn {
		l.pos += Pos(len(delim))
		l.emitHeader(it, depth)
		l.ignore()
		l.parenDepth = 0
		return lexText
	}
}
