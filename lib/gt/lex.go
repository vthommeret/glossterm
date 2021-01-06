package gt

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Adapted from https://github.com/golang/go/blob/master/src/text/template/parse/lex.go
// Also see https://www.youtube.com/watch?v=HxaD_trXwRE

// TODO: Can HasPrefix use be minimized vs. using peek/accept from original lex.go?

// Pos represents a byte position in the original input text from which
// this template was parsed.
type Pos int

// --------

// item represents a token or text string returned from the scanner.
type item struct {
	typ      itemType // The type of this item.
	pos      Pos      // The starting position, in bytes, of this item in the input string.
	val      string   // The value of this item.
	depth    int      // Used by itemHeaderStart and itemHeaderEnd to indicate header depth.
	balanced bool     // Used for itemLeftTemplate to indicate whether it has a matching itemRightTemplate.
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
	itemAction      // Template action
	itemParamDelim  // Template parameter delimiter
	itemParamName   // Template parameter name
	itemParamText   // Template parameter text
	itemHeaderStart // Header start
	itemHeaderEnd   // Header end

	itemUnorderedListItemStart   // Unordered list item start
	itemUnorderedDefinitionStart // Unordered definition start
	itemUnorderedOrderedStart    // Unordered ordered start

	itemOrderedListItemStart   // Ordered list item start
	itemOrderedDefinitionStart // Ordered definition start
	itemOrderedUnorderedStart  // Ordered unordered start

	itemListItemPrefix // List item prefix
	itemListItemEnd    // List item end
	itemText           // Plain text
	itemLeftTemplate   // Left template delimiter
	itemRightTemplate  // Right template delimiter
	itemLeftLink       // Left link delimiter
	itemRightLink      // Right link delimiter
	itemLink           // Link text
	itemLinkDelim      // Link delimiter
	itemLinkName       // Link name
)

// Make the types prettyprint.
var itemName = map[itemType]string{
	itemError:       "error",
	itemEOF:         "EOF",
	itemAction:      "action",
	itemParamDelim:  "param delim",
	itemParamName:   "param name",
	itemParamText:   "param text",
	itemHeaderStart: "header start",
	itemHeaderEnd:   "header end",

	itemUnorderedListItemStart:   "unordered list item start",
	itemUnorderedDefinitionStart: "unordered definition start",
	itemUnorderedOrderedStart:    "unordered ordered start",

	itemOrderedListItemStart:   "ordered list item start",
	itemOrderedDefinitionStart: "ordered definition start",
	itemOrderedUnorderedStart:  "ordered unordered start",

	itemListItemPrefix: "list item prefix",
	itemListItemEnd:    "list item end",
	itemText:           "text",
	itemLeftTemplate:   "left template",
	itemRightTemplate:  "right template",
	itemLeftLink:       "left link",
	itemRightLink:      "right link",
	itemLink:           "link",
	itemLinkDelim:      "link delim",
	itemLinkName:       "link name",
}

// From http://w3c.github.io/html/syntax.html#void-elements
var voidTags = []string{
	"area",
	"base",
	"br",
	"col",
	"embed",
	"hr",
	"img",
	"input",
	"link",
	"menuitem",
	"meta",
	"param",
	"source",
	"track",
	"wbr",
}
var voidTagMap = map[string]bool{}

type buffer struct {
	buffering bool
	items     []*item // items buffered
	tpls      []*item // pointers to template opens
	openTpls  stack   // stack of open template indices
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
	items      chan item // channel of scanned items
	buffered   buffer    // buffer of items before they're emitted
	tplDepth   int       // template depth
	inListItem bool
	debug      bool // output debug messages automatically
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

	if l.debug {
		l.printDebug("next", "")
	}

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
	if l.debug {
		l.printDebug("backup", "")
	}
}

// advanced steps forward one rune. Can only be called once per call of next.
func (l *lexer) advance() {
	l.pos += l.width
	if l.debug {
		l.printDebug("advance", "")
	}
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	if l.debug {
		var val string
		if l.pos > l.start {
			val = fmt.Sprintf(" → %s", l.input[l.start:l.pos])
		}
		l.printDebug("emit", fmt.Sprintf("%s%s", t.String(), val))
	}
	if l.buffered.buffering {
		l.buffer(t)
	} else {
		l.items <- item{t, l.start, l.input[l.start:l.pos], 0, false}
		l.start = l.pos
	}
}

// buffer an item for emit, in case we need to backtrack.
func (l *lexer) buffer(t itemType) {
	i := item{t, l.start, l.input[l.start:l.pos], 0, false}
	l.buffered.items = append(l.buffered.items, &i)
	switch t {
	case itemLeftTemplate:
		l.buffered.tpls = append(l.buffered.tpls, &i)
		l.buffered.openTpls.push(len(l.buffered.tpls) - 1)
	case itemRightTemplate:
		lastOpen := l.buffered.openTpls.pop()
		l.buffered.tpls[lastOpen].balanced = true
	}
	l.start = l.pos
}

// emitDepth passes a item with depth back to the client.
func (l *lexer) emitDepth(t itemType, depth int) {
	if l.debug {
		l.printDebug("emit", fmt.Sprintf("%s → %s", t.String(), l.input[l.start:l.pos]))
	}
	l.items <- item{t, l.start, l.input[l.start:l.pos], depth, false}
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
	if l.debug {
		l.printDebug("ignore", "")
	}
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
	l.items <- item{itemError, l.start, fmt.Sprintf(format, args...), 0, false}
	return nil
}

func (l *lexer) lexListItem(startDelim string, startItem itemType, endItem itemType, definitionItem itemType, alternateDelim string, alternateItem itemType) (bool, stateFn) {
	for i := listItems; i > 0; i-- {
		delim := strings.Repeat(startDelim, i)

		// Definition list item (#:)
		defnPrefix := fmt.Sprintf("%s%s", delim, definitionStart)
		if strings.HasPrefix(l.input[l.pos:], defnPrefix) {
			if l.inListItem {
				l.emit(endItem)
			}
			l.inListItem = true
			return true, lexDelim(definitionItem, defnPrefix, i)
		}

		// For ordered list within unordered list, or visa versa...
		// From "nonhospital":
		// # An [[institution]] that is not a hospital.
		// #* '''1988''', LaVonne Straub, ‎Norman Walzer, ''Financing Rural Health Care'' (page 97)
		alternatePrefix := fmt.Sprintf("%s%s", delim, alternateDelim)
		if strings.HasPrefix(l.input[l.pos:], alternatePrefix) {
			if l.inListItem {
				l.emit(endItem)
			}
			l.inListItem = true
			return true, lexDelim(alternateItem, alternatePrefix, i)
		}

		// Regular list item
		prefix := fmt.Sprintf("%s", delim)
		if strings.HasPrefix(l.input[l.pos:], prefix) {
			if l.inListItem {
				l.emit(endItem)
			} else {
				l.inListItem = true
			}
			if l.pos > l.start {
				l.emit(itemText)
			}
			l.ignore()
			return true, lexDelim(startItem, prefix, i)
		}
	}
	return false, nil
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
		input: strings.TrimSpace(input),
		items: make(chan item),
	}
	go l.run()
	return l
}

func (l *lexer) printDebug(fn, val string) {
	var leftCursor, inner, rightCursor string
	width := l.pos - l.start
	prefix := strings.Repeat(" ", int(l.start))
	if width == 0 {
		leftCursor = "↑"
	} else {
		leftCursor = "↑"
		rightCursor = "↑"
		inner = strings.Repeat(" ", int(width-1))
	}
	action := fmt.Sprintf("# %s", fn)
	if val != "" {
		action = fmt.Sprintf("%s %s", action, val)
	}
	fmt.Printf("%s  %s\n%s%s%s%s\n\n", l.input, action, prefix, leftCursor, inner, rightCursor)
}

// Print prints all items
func (l *lexer) Print() {
	fmt.Print(l.String(l.Items()))
}

func (l *lexer) String(items []item) string {
	var b strings.Builder
	for _, i := range items {
		if i.typ == itemEOF || i.typ == itemError {
			if i.typ == itemError {
				fmt.Fprintf(&b, "Error: %s\n", i.val)
			}
			break
		}
		if i.typ == itemHeaderStart || i.typ == itemHeaderEnd || i.typ == itemUnorderedListItemStart || i.typ == itemOrderedListItemStart || i.typ == itemOrderedDefinitionStart || i.typ == itemUnorderedDefinitionStart {
			fmt.Fprintf(&b, "%s, %d: %q\n", i.typ, i.depth, i.val)
		} else {
			fmt.Fprintf(&b, "%s: %q\n", i.typ, i.val)
		}
	}
	return b.String()
}

func (l *lexer) Items() []item {
	var items []item
	for {
		i := l.NextItem()
		if i.typ == itemEOF {
			break
		}
		items = append(items, i)
	}
	return items
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for l.state = lexText; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.items)
}

// drainTplBuffer drains the tpl buffer
func (l *lexer) drainTplBuffer() {
	s := item{typ: itemText}
	openTpls := stack{}
	for n, i := range l.buffered.items {
		if i.typ == itemLeftTemplate {
			openTpls.push(n)
		}
		tpl := openTpls.peek()
		if tpl != -1 && l.buffered.items[tpl].balanced {
			if s.val != "" {
				l.items <- s
				s = item{typ: itemText}
			}
			l.items <- *i
		} else {
			if s.val == "" {
				s.pos = i.pos
			}
			s.val += i.val
		}
		if i.typ == itemRightTemplate {
			openTpls.pop()
		}
	}
	if s.val != "" {
		l.items <- s
	}
	l.buffered = buffer{}
}

// drainLinkBuffer drains the link buffer
func (l *lexer) drainLinkBuffer(asText bool) {
	s := item{typ: itemText}
	for _, i := range l.buffered.items {
		if asText {
			if s.val == "" {
				s.pos = i.pos
			}
			s.val += i.val
		} else {
			l.items <- *i
		}
	}
	if s.val != "" {
		l.items <- s
	}
	l.buffered = buffer{}
}

// state functions

const (
	headers                = 6
	headerDelim            = "="
	listItems              = 6
	orderedListItemStart   = "#"
	unorderedListItemStart = "*"
	definitionStart        = ":"
	listItemEndDelim       = ": "
	leftTemplate           = "{{"
	rightTemplate          = "}}"
	paramEqual             = "="
	paramDelim             = "|"
	leftLink               = "[["
	rightLink              = "]]"
	linkDelim              = "|"
	spaceChars             = " \t\r\n"
)

// lexText scans until a header or template delimiter.
func lexText(l *lexer) stateFn {
Loop:
	for {
		startOfLine := true
		if l.pos > 0 {
			startOfLine = (l.input[l.pos-1:l.pos] == "\n")
		}

		if startOfLine {
			// Headers
			for i := headers; i > 0; i-- {
				delim := strings.Repeat(headerDelim, i)
				if strings.HasPrefix(l.input[l.pos:], delim) {
					if l.inListItem {
						l.inListItem = false
						l.emit(itemListItemEnd)
					}
					if l.pos > l.start {
						l.emit(itemText)
					}
					l.ignore()
					return lexDelim(itemHeaderStart, delim, i)
				}
			}
			// Unordered list items
			if lexed, nextFn := l.lexListItem(unorderedListItemStart, itemUnorderedListItemStart, itemListItemEnd, itemUnorderedDefinitionStart, orderedListItemStart, itemOrderedUnorderedStart); lexed {
				return nextFn
			}
			// Ordered list items
			if lexed, nextFn := l.lexListItem(orderedListItemStart, itemOrderedListItemStart, itemListItemEnd, itemOrderedDefinitionStart, unorderedListItemStart, itemUnorderedOrderedStart); lexed {
				return nextFn
			}

			// If newline without header or another list item, end existing list item
			if l.inListItem {
				l.inListItem = false
				l.emit(itemListItemEnd)
			}
		}

		for i := headers; i > 0; i-- {
			delim := strings.Repeat(headerDelim, i)
			if strings.HasPrefix(l.input[l.pos:], delim+"\n") || l.input[l.pos:] == delim {
				if l.pos > l.start {
					l.emit(itemText)
				}
				l.ignore()
				return lexDelim(itemHeaderEnd, delim, i)
			}
		}
		/*
			if strings.HasPrefix(l.input[l.pos:], listItemEndDelim) {
				if l.pos > l.start {
					l.emit(itemListItemPrefix)
				}
				l.pos += Pos(len(listItemEndDelim))
				l.emit(itemListItemEnd)
			}
		*/
		if strings.HasPrefix(l.input[l.pos:], leftTemplate) {
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexLeftTemplate
		} else if strings.HasPrefix(l.input[l.pos:], rightTemplate) {
			// TODO: Can this condition be removed?
			// Need to properly support unclosed "}}".
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexRightTemplate
		} else if strings.HasPrefix(l.input[l.pos:], leftLink) {
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexLeftLink
		}

		switch r := l.next(); {
		case r == eof:
			l.drainTplBuffer()
			if l.inListItem {
				if l.pos > l.start {
					l.emit(itemText)
				}
				l.emit(itemListItemEnd)
				l.inListItem = false
			}
			break Loop
		case isEndOfLine(r):
			if l.pos > l.start {
				l.emit(itemText)
			}
			continue Loop
		default:
			// absorb
		}
	}
	// Correctly reached EOF.
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.emit(itemEOF)
	return nil
}

func lexDelim(it itemType, delim string, depth int) stateFn {
	return func(l *lexer) stateFn {
		l.pos += Pos(len(delim))
		l.emitDepth(it, depth)
		return lexText
	}
}

// lexLeftTemplate scans the left template delimiter.
func lexLeftTemplate(l *lexer) stateFn {
	l.buffered.buffering = true
	l.pos += Pos(len(leftTemplate))
	l.emit(itemLeftTemplate)
	return lexAction
}

// lexRightTemplate scans the right template delimiter.
func lexRightTemplate(l *lexer) stateFn {
	l.pos += Pos(len(rightTemplate))
	l.emit(itemRightTemplate)
	if len(l.buffered.openTpls) == 0 {
		l.drainTplBuffer()
	}
	return lexText
}

// lexAction scans a template action.
func lexAction(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == eof:
			// Unclosed template -- eof.
			l.drainTplBuffer()
			return lexText
		case isEndOfLine(r):
			if strings.HasPrefix(l.input[l.pos:], headerDelim) {
				// Unclosed template -- header start.
				l.backup()
				if l.pos > l.start {
					l.emit(itemAction)
				}
				l.drainTplBuffer()
				return lexText
			} else {
				l.backup()
				if l.pos > l.start {
					l.emit(itemAction)
				}
				l.ignore()
				return lexParam
			}
		case r == '}':
			if strings.HasPrefix(l.input[l.pos:], "}") {
				l.backup()
				if l.pos > l.start {
					l.emit(itemAction)
				}
				return lexRightTemplate
			}
		case r == '|':
			l.backup()
			if l.pos > l.start {
				l.emit(itemAction)
			}
			return lexParam
		default:
			// absorb.
		}
	}
	return lexText
}

// lexParam scans a template parameter.
func lexParam(l *lexer) stateFn {
	var inParam bool
	var inNamedParam bool
	var emittedEndOfLineParam bool

	var openStartTag bool
	var inTag bool
	var openCloseTag bool
	var openStartTagPos Pos

	var inLink bool

	for {
		switch r := l.next(); {
		case r == eof:
			// Unclosed template -- eof.
			l.backup()
			if l.pos > l.start {
				l.emit(itemParamText)
			}
			l.drainTplBuffer()
			return lexText
		case isEndOfLine(r):
			if emittedEndOfLineParam {
				emittedEndOfLineParam = false
			} else {
				if strings.HasPrefix(l.input[l.pos:], "|") || strings.HasPrefix(l.input[l.pos:], rightTemplate) {
					l.backup()
					if l.pos > l.start {
						l.emit(itemParamText)
						emittedEndOfLineParam = true
					}
					l.pos += Pos(1)
					l.ignore()
					inNamedParam = false
				} else if strings.HasPrefix(l.input[l.pos:], "=") {
					// Unclosed template -- header start.
					l.backup()
					if l.pos > l.start {
						l.emit(itemParamText)
					}
					l.drainTplBuffer()
					return lexText
				}
			}
		case r == '<':
			if openStartTag {
				openStartTag = false
			} else if inTag {
				inTag = false
				openCloseTag = true
			} else {
				openStartTag = true
				openStartTagPos = l.pos
			}
		case r == '>':
			if openStartTag {
				tag := l.input[openStartTagPos : l.pos-1]
				if _, ok := voidTagMap[tag]; !ok {
					inTag = true
				}
				openStartTag = false
			} else if inTag {
				inTag = false
			} else if openCloseTag {
				openCloseTag = false
			}
		case r == '[':
			if !inTag {
				if n := l.peek(); n == '[' {
					inLink = true
				}
			}
			if openStartTag {
				openStartTag = false
			}
			if openCloseTag {
				openCloseTag = false
			}
		case r == ']':
			if !inTag {
				if n := l.peek(); n == ']' {
					inLink = false
				}
			}
			if openStartTag {
				openStartTag = false
			}
			if openCloseTag {
				openCloseTag = false
			}
		case r == '=':
			if !inTag && !inNamedParam {
				l.backup()
				l.emit(itemParamName)
				l.pos += Pos(len(paramEqual))
				l.ignore()
				inNamedParam = true
			}
			if openStartTag {
				openStartTag = false
			}
			if openCloseTag {
				openCloseTag = false
			}
		case r == '|':
			if !inTag && !inLink {
				if inParam {
					l.backup()
					if emittedEndOfLineParam {
						emittedEndOfLineParam = false
					} else {
						l.emit(itemParamText)
					}
					l.pos += Pos(len(paramDelim))
				}
				l.emit(itemParamDelim)
				inParam = true
				inNamedParam = false
			}
			if openStartTag {
				openStartTag = false
			}
			if openCloseTag {
				openCloseTag = false
			}
		case r == '{':
			if !inTag {
				if strings.HasPrefix(l.input[l.pos:], "{") {
					l.backup()
					if l.pos > l.start {
						l.emit(itemParamText)
					}
					return lexLeftTemplate
				}
			}
			if openStartTag {
				openStartTag = false
			}
			if openCloseTag {
				openCloseTag = false
			}
		case r == '}':
			if strings.HasPrefix(l.input[l.pos:], "}") {
				l.backup()
				if !emittedEndOfLineParam && l.pos > l.start {
					l.emit(itemParamText)
				}
				return lexRightTemplate
			}
			if openStartTag {
				openStartTag = false
			}
			if openCloseTag {
				openCloseTag = false
			}
		default:
			if !isAlphaNumeric(r) {
				if openStartTag {
					openStartTag = false
				}
				if openCloseTag {
					openCloseTag = false
				}
			}
			if !inParam && !isWhitespace(r) {
				// Unclosed template (invalid character)
				l.backup()
				l.ignore() // Ignore previous whitespace.
				l.drainTplBuffer()
				return lexText
			}
			// absorb.
		}
	}
	return lexAction
}

// lexLeftLink scans the left link delimiter.
func lexLeftLink(l *lexer) stateFn {
	_ = "breakpoint"
	l.buffered.buffering = true
	l.pos += Pos(len(leftLink))
	l.emit(itemLeftLink)
	return lexLink
}

// lexRightLink scans the right link delimiter.
func lexRightLink(l *lexer) stateFn {
	_ = "breakpoint"
	l.pos += Pos(len(rightLink))
	l.emit(itemRightLink)
	l.drainLinkBuffer(false)
	return lexText
}

// lexLink scans a link.
func lexLink(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == eof:
			_ = "breakpoint"
			// Unclosed link -- eof.
			if l.pos > l.start {
				l.emit(itemLink)
			}
			l.drainLinkBuffer(true)
			return lexText
		case isEndOfLine(r):
			_ = "breakpoint"
			// Unclosed link -- end of line.
			l.emit(itemLink)
			l.drainLinkBuffer(true)
			return lexText
		case r == '[', r == '{', r == '}':
			_ = "breakpoint"
			// Invalid link -- template start/close.
			l.backup()
			if l.pos > l.start {
				l.emit(itemLink)
			}
			l.drainLinkBuffer(true)
			return lexText
		case r == ']':
			_ = "breakpoint"
			if strings.HasPrefix(l.input[l.pos:], "]") {
				l.backup()
				if l.pos > l.start {
					l.emit(itemLink)
				}
				return lexRightLink
			}
		case r == '|':
			l.backup()
			if l.pos > l.start {
				l.emit(itemLink)
			}
			l.pos += Pos(len(linkDelim))
			l.emit(itemLinkDelim)
			return lexLinkName
		default:
			_ = "breakpoint"
			// absorb.
		}
	}
	return lexText
}

// lexLinkName scans a link name.
func lexLinkName(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == eof:
			_ = "breakpoint"
			// Unclosed link name -- eof.
			if l.pos > l.start {
				l.emit(itemLinkName)
			}
			l.drainLinkBuffer(true)
			return lexText
		case isEndOfLine(r):
			_ = "breakpoint"
			// Unclosed link name -- end of line.
			l.emit(itemLinkName)
			l.drainLinkBuffer(true)
			return lexText
		case r == '[', r == '{', r == '}':
			_ = "breakpoint"
			// Invalid link name -- template start/close.
			l.backup()
			if l.pos > l.start {
				l.emit(itemLinkName)
			}
			l.drainLinkBuffer(true)
			return lexText
		case r == ']':
			_ = "breakpoint"
			if strings.HasPrefix(l.input[l.pos:], "]") {
				l.backup()
				if l.pos > l.start {
					l.emit(itemLinkName)
				}
				return lexRightLink
			}
		default:
			_ = "breakpoint"
			// absorb.
		}
	}
	return lexText
}

// isEndOfLine reports whether r is an end-of-line character.
func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

// isAlphaNumeric reports whether r is an ASCII letter or digit.
func isAlphaNumeric(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
}

// isWhitespace reports whether r is a whitespace character.
func isWhitespace(r rune) bool {
	return (r == ' ' || r == '\t' || r == '\r' || r == '\n')
}

func init() {
	for _, t := range voidTags {
		voidTagMap[t] = true
	}
}

// Simple stack implementation. Doesn't handle popping/peeking empty stacks
// since it should only be called when the buffer can be balanced.

type stack []int

func (s *stack) push(v int) {
	*s = append(*s, v)
}

func (s *stack) pop() int {
	res := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return res
}

func (s *stack) peek() int {
	return (*s)[len(*s)-1]
}
