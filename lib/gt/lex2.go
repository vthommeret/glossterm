package gt

import (
	"strings"
)

const newline = '\n'

// Adapted from https://github.com/golang/go/blob/master/src/text/template/parse/lex.go
// Also see https://www.youtube.com/watch?v=HxaD_trXwRE

// NewLexer creates a new scanner for the input string.
func NewLexer2(input string, debug bool) *lexer {
	l := &lexer{
		input: input,
		items: make(chan item),
		debug: debug,
	}
	go l.run2()
	return l
}

// run runs the state machine for the lexer.
// TODO: Ignore whitespace
func (l *lexer) run2() {
	if l.debug {
		l.printDebug("start", "")
	}
	for l.state = lexText2; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.items)
}

func (l *lexer) remaining() string {
	return l.input[l.pos:]
}

func (l *lexer) emitAnyText() {
	if l.pos > l.start {
		l.emit(itemText)
	}
}

func lexText2(l *lexer) stateFn {
Loop:
	for {
		startOfLine := (l.pos == 0 || l.input[l.pos-1:l.pos] == string(newline))
		remaining := l.remaining()

		// Lex headers
		for i := headers; i > 0; i-- {
			delim := strings.Repeat(headerDelim, i)
			endDelim := delim + string(newline)

			// End delim
			if remaining == delim || strings.HasPrefix(remaining, endDelim) {
				l.emitAnyText()
				return lexDelim2(itemHeaderEnd, delim, i)
			}

			// Start delim
			if startOfLine && strings.HasPrefix(remaining, delim) {
				l.emitAnyText()
				return lexDelim2(itemHeaderStart, delim, i)
			}
		}

		// Lex templates
		if strings.HasPrefix(remaining, leftTemplate) {
			l.emitAnyText()
			return lexLeftTemplate2
		}
		// TODO: Share with template parsing code
		if strings.HasPrefix(remaining, rightTemplate) {
			l.emitAnyText()
			return lexRightTemplate2
		}

		// Lex links
		if strings.HasPrefix(remaining, leftLink) {
			l.emitAnyText()
			return lexLeftLink2
		}
		// TODO: Share with link parsing code
		if strings.HasPrefix(remaining, rightLink) {
			l.emitAnyText()
			return lexRightLink2
		}

		switch r := l.next(); {
		case r == eof:
			break Loop

			// TODO: Share cases with template parsing code
		case r == '|':
			// TODO: Update to check that active context is a template
			if l.tplDepth > 0 {
				l.backup()
				l.emitAnyText()
				l.advance()
				l.emit(itemParamDelim)
			}
		case r == '=':
			// TODO: Update to check that active context is a template
			if l.tplDepth > 0 {
				l.backup()
				l.emit(itemParamName)
				l.advance()
				l.ignore()
			}
		}

	}

	// Correctly reached EOF
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.emit(itemEOF)

	return nil
}

// Rename lexHeader?
func lexDelim2(it itemType, delim string, depth int) stateFn {
	return func(l *lexer) stateFn {
		l.pos += Pos(len(delim))
		l.emitDepth(it, depth)
		return lexText2
	}
}

// lexLeftTemplate scans the left template delimiter.
func lexLeftTemplate2(l *lexer) stateFn {
	l.pos += Pos(len(leftTemplate))
	l.tplDepth++
	l.emit(itemLeftTemplate)
	return lexAction2
}

// lexRightTemplate scans the right template delimiter.
func lexRightTemplate2(l *lexer) stateFn {
	l.pos += Pos(len(rightTemplate))
	l.tplDepth--
	l.emit(itemRightTemplate)
	return lexText2
}

// lexAction scans template action
func lexAction2(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == '|':
			l.backup()
			if l.pos > l.start {
				l.emit(itemAction)
			}
			l.advance()
			l.emit(itemParamDelim)
			return lexText2
		case r == '}':
			if strings.HasPrefix(l.remaining(), "}") {
				l.backup()
				if l.pos > l.start {
					l.emit(itemAction)
				}
				return lexRightTemplate2
			}
		}
	}
}

// lexLeftLink2 scans the left link delimiter.
func lexLeftLink2(l *lexer) stateFn {
	l.pos += Pos(len(leftLink))
	l.emit(itemLeftLink)
	return lexLink2
}

// lexRightLink2 scans the right link delimiter.
func lexRightLink2(l *lexer) stateFn {
	l.pos += Pos(len(rightLink))
	l.emit(itemRightLink)
	return lexText2
}

// lexLink scans a link
func lexLink2(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == '|':
			l.backup()
			if l.pos > l.start {
				l.emit(itemLink)
			}
			l.advance()
			l.emit(itemLinkDelim)
			return lexText2
		case r == ']':
			if strings.HasPrefix(l.remaining(), "]") {
				l.backup()
				if l.pos > l.start {
					l.emit(itemLink)
				}
				return lexRightLink2
			}
		}
	}
}
