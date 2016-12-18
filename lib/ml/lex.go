package ml

import (
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
	itemAction        // Template action
	itemParam         // Template parameter
	itemParamName     // Template parameter name
	itemParamValue    // Template parameter value
	itemPipe          // Pipe symbol
	itemHeaderStart   // Header start
	itemHeaderEnd     // Header end
	itemText          // Plain text
	itemLeftTemplate  // Left template delimiter
	itemRightTemplate // Right template delimiter
)

// Make the types prettyprint.
var itemName = map[itemType]string{
	itemError:         "error",
	itemEOF:           "EOF",
	itemAction:        "action",
	itemParam:         "param",
	itemParamName:     "param name",
	itemParamValue:    "param value",
	itemPipe:          "pipe",
	itemHeaderStart:   "header start",
	itemHeaderEnd:     "header end",
	itemText:          "text",
	itemLeftTemplate:  "left template",
	itemRightTemplate: "right template",
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
	input string    // the string being scanned
	state stateFn   // the next lexing function to enter
	pos   Pos       // current position in the input
	start Pos       // start position of this item
	width Pos       // width of last rune read from input
	items chan item // channel of scanned items
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
	return item
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

// Simplifies parsing.
func normalize(s string) string {
	return fmt.Sprintf("\n%s\n", strings.TrimSpace(s))
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
	headerDelim   = "="
	headers       = 6
	leftTemplate  = "{{"
	rightTemplate = "}}"
)

// lexText scans until a header or template delimiter.
func lexText(l *lexer) stateFn {
	for {
		for i := headers; i > 0; i-- {
			delim := strings.Repeat(headerDelim, i)

			leftHeader := "\n" + delim
			if strings.HasPrefix(l.input[l.pos:], leftHeader) {
				if l.pos > l.start {
					l.emit(itemText)
				}
				l.ignore()
				return lexHeaderDelim(itemHeaderStart, leftHeader, i)
			}

			rightHeader := delim + "\n"
			if strings.HasPrefix(l.input[l.pos:], rightHeader) {
				if l.pos > l.start {
					l.emit(itemText)
				}
				l.ignore()
				return lexHeaderDelim(itemHeaderEnd, rightHeader, i)
			}
		}
		if strings.HasPrefix(l.input[l.pos:], leftTemplate) {
			if l.pos > l.start {
				l.emit(itemText)
			}
			l.ignore()
			return lexLeftTemplate
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
		return lexText
	}
}

// lexLeftTemplate scans the left template delimiter.
func lexLeftTemplate(l *lexer) stateFn {
	l.pos += Pos(len(leftTemplate))
	l.emit(itemLeftTemplate)
	l.ignore()
	return lexInsideAction
}

// lexRightTemplate scans the right template delimiter.
func lexRightTemplate(l *lexer) stateFn {
	l.pos += Pos(len(rightTemplate))
	l.emit(itemRightTemplate)
	return lexText
}

// lexInsideAction scans the elements inside action delimiters.
func lexInsideAction(l *lexer) stateFn {
	// Either number, quoted string, or identifier.
	// Spaces separate arguments; runs of spaces turn into itemSpace.
	// Pipe symbols separate and are emitted.
	if strings.HasPrefix(l.input[l.pos:], rightTemplate) {
		return lexRightTemplate
	}
	switch r := l.next(); {
	case r == '|':
		l.emit(itemPipe)
	default:
		l.backup()
		return lexAction
	}
	return lexInsideAction
}

// lexAction scans a template action.
func lexAction(l *lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case r == eof || isEndOfLine(r):
			return l.errorf("unclosed action")
		case r == '}':
			l.backup()
			l.emit(itemAction)
			break Loop
		default:
			// absorb.
			if l.peek() == '|' {
				l.emit(itemAction)
				return lexParam
			}
		}
	}
	return lexInsideAction
}

// lexParam scans a template parameter.
func lexParam(l *lexer) stateFn {
	var kv bool
Loop:
	for {
		switch r := l.next(); {
		case r == eof || isEndOfLine(r):
			return l.errorf("unclosed action")
		case r == '=':
			l.ignore()
		case r == '|':
			l.ignore()
		case r == '}':
			l.backup()
			l.emitParam(kv)
			break Loop
		default:
			// absorb.
			switch l.peek() {
			case '=':
				l.emit(itemParamName)
				kv = true
			case '|':
				l.emitParam(kv)
				kv = false
			}
		}
	}
	return lexInsideAction
}

func (l *lexer) emitParam(kv bool) {
	if kv {
		l.emit(itemParamValue)
	} else {
		l.emit(itemParam)
	}
}

// isEndOfLine reports whether r is an end-of-line character.
func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}
