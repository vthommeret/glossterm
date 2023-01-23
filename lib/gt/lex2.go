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
//
// General flow is characters within template parameter values, HTML tags and markup are
// all considered text and are processed by "lexText". This allows arbitrary of nesting
// of templates, HTML tags, markup within text. HTML attributes, HTML comments, template
// action names have their own parsers and don't support nesting (e.g. evaluating a
// template within an HTML attribute value)
func (l *lexer) run2() {
	if l.debug {
		l.printDebug("start", "")
	}
	for l.state = lexText2; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.items)
}

func (l *lexer) close() {
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.emit(itemEOF)
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
				return lexHeaderDelim2(itemHeaderEnd, delim, i)
			}

			// Start delim
			if startOfLine && strings.HasPrefix(remaining, delim) {
				l.emitAnyText()
				return lexHeaderDelim2(itemHeaderStart, delim, i)
			}
		}

		// Lex templates
		if strings.HasPrefix(remaining, leftTemplate) {
			l.emitAnyText()
			return lexLeftTemplate2
		}
		if strings.HasPrefix(remaining, rightTemplate) {
			l.emitAnyText()
			return lexRightTemplate2
		}

		// Lex links
		if strings.HasPrefix(remaining, leftLink) {
			l.emitAnyText()
			return lexLeftLink2
		}
		if strings.HasPrefix(remaining, rightLink) {
			l.emitAnyText()
			return lexRightLink2
		}

		// Lex markup
		if strings.HasPrefix(remaining, strongEmphasized) {
			l.emitAnyText()
			return lexStrongEmphasized2
		}
		if strings.HasPrefix(remaining, strong) {
			l.emitAnyText()
			return lexStrong2
		}
		if strings.HasPrefix(remaining, emphasized) {
			l.emitAnyText()
			return lexEmphasized2
		}

		// Lex HTML
		if strings.HasPrefix(remaining, tagCommentLeft) {
			l.emitAnyText()
			return lexTagCommentLeft
		}
		if strings.HasPrefix(remaining, closeTagLeft) {
			l.emitAnyText()
			return lexCloseTagLeft
		}
		if strings.HasPrefix(remaining, openTagLeft) {
			l.emitAnyText()
			return lexOpenTagLeft
		}
		if strings.HasPrefix(remaining, tagRight) {
			l.emitAnyText()
			return lexTagRight
		}

		// Lex remaining characters
		switch r := l.next(); {
		case r == eof:
			break Loop

			// Template parameter delimiter
		case r == '|':
			if l.tplDepth > 0 {
				l.backup()
				l.emitAnyText()
				l.advance()
				l.emit(itemParamDelim)
			}
			// Template parameter name delimiter
		case r == '=':
			if l.tplDepth > 0 {
				l.backup()
				l.emit(itemParamName)
				l.advance()
				l.ignore()
			}
		}

	}

	// Correctly reached EOF
	l.close()

	return nil
}

// lexHeaderDelim2 scans header delimiters
func lexHeaderDelim2(it itemType, delim string, depth int) stateFn {
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
		case r == eof:
			l.close()
			return nil
		}
	}
}

// lexStrongEmphasized2 scans the strong emphasized text
func lexStrongEmphasized2(l *lexer) stateFn {
	l.pos += Pos(len(strongEmphasized))
	l.emit(itemStrongEmphasized)
	return lexText2
}

// lexStrong2 scans the strong text
func lexStrong2(l *lexer) stateFn {
	l.pos += Pos(len(strong))
	l.emit(itemStrong)
	return lexText2
}

// lexEmphasized2 scans the strong text
func lexEmphasized2(l *lexer) stateFn {
	l.pos += Pos(len(emphasized))
	l.emit(itemEmphasized)
	return lexText2
}

// lexTagCommentLeft scans HTML comment left delimiters (<!--)
func lexTagCommentLeft(l *lexer) stateFn {
	l.pos += Pos(len(tagCommentLeft))
	l.emit(itemTagCommentLeft)
	return lexTagComment
}

// lexTagComment scans HTML comments
func lexTagComment(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == '-':
			if strings.HasPrefix(l.remaining(), "->") {
				l.backup()
				if l.pos > l.start {
					l.emit(itemTagComment)
				}
				return lexTagCommentRight
			}
		}
	}
}

// lexTagCommentRight scans HTML comment right delimiters (-->)
func lexTagCommentRight(l *lexer) stateFn {
	l.pos += Pos(len(tagCommentRight))
	l.emit(itemTagCommentRight)
	return lexText2
}

// lexOpenTagLeft scans open HTML tag left delimiters (<)
func lexOpenTagLeft(l *lexer) stateFn {
	l.pos += Pos(len(openTagLeft))
	l.emit(itemOpenTagLeft)
	return lexOpenTagName
}

// lexCloseTagLeft scans close HTML tag left delimiters (</)
func lexCloseTagLeft(l *lexer) stateFn {
	l.pos += Pos(len(closeTagLeft))
	l.emit(itemCloseTagLeft)
	return lexCloseTagName
}

// lexCloseTagRight scans close HTML tag right delimiters (/>)
func lexCloseTagRight(l *lexer) stateFn {
	l.pos += Pos(len(closeTagRight))
	l.emit(itemCloseTagRight)
	return lexText2
}

// lexTagRight scans HTML tag right delimiters (>)
func lexTagRight(l *lexer) stateFn {
	l.pos += Pos(len(tagRight))
	l.emitTrim(itemTagRight)
	return lexText2
}

// lexOpenTagName scans HTML opening tag names (e.g. span)
func lexOpenTagName(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == ' ':
			l.backup()
			if l.pos > l.start {
				l.emit(itemTagName)
			}
			return lexTagAttrName
		case r == '/':
			if strings.HasPrefix(l.remaining(), ">") {
				l.backup()
				if l.pos > l.start {
					l.emit(itemTagName)
				}
				return lexCloseTagRight
			}
		case r == '>':
			l.backup()
			if l.pos > l.start {
				l.emit(itemTagName)
			}
			return lexTagRight
		}
	}
}

// lexCloseTagName scans HTML closing tag names (e.g. span)
func lexCloseTagName(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == '>':
			l.backup()
			if l.pos > l.start {
				l.emit(itemTagName)
			}
			return lexTagRight
		}
	}
}

// lexTagAttrName scans HTML tag attribute names (e.g. style)
func lexTagAttrName(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == '=':
			l.backup()
			if l.pos > l.start {
				l.emitTrim(itemTagAttrName)
			}
			l.advance()
			l.ignore()
			return lexTagAttrValueLeft
		case r == '/':
			if strings.HasPrefix(l.remaining(), ">") {
				l.backup()
				if l.pos > l.start {
					l.emitTrim(itemTagAttrName)
				}
				return lexCloseTagRight
			}
		case r == '>':
			l.backup()
			if l.pos > l.start {
				l.emitTrim(itemTagAttrName)
			}
			return lexTagRight
		}
	}
}

// lexTagAttrValueLeft scans HTML tag attribute left delimiter (")
func lexTagAttrValueLeft(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == '\'', r == '"':
			l.ignore()
			return lexTagAttrValue
		case isWhitespace(r):
			l.ignore()
		default:
			return lexTagAttrValue
		}
	}
}

// lexTagAttrValue scans HTML tag attribute value (e.g. "color: red")
func lexTagAttrValue(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == '\'', r == '"':
			l.backup()
			if l.pos > l.start {
				l.emit(itemTagAttrValue)
			}
			l.advance()
			l.ignore()
			return lexTagAttrName
		case r == '/':
			if strings.HasPrefix(l.remaining(), ">") {
				l.backup()
				if l.pos > l.start {
					l.emit(itemTagAttrValue)
				}
				return lexCloseTagRight
			}
		case r == '>':
			l.backup()
			if l.pos > l.start {
				l.emit(itemTagAttrValue)
			}
			return lexTagRight
		}
	}
}
