package http_router

import (
	"bufio"
	"bytes"
	"strings"
)

const (
	tSlash = iota
	tStatic
	tOpenVar
	tVar
	tCloseVar
	tEnd
	tExpReg
)

const (
	tModeStatic = iota
	tModeIdentifier
)

type lexer struct {
	mode int
	buf  *bufio.Reader
}

func newLexer(path string) *lexer {
	reader := strings.NewReader(path)
	return &lexer{buf: bufio.NewReader(reader), mode: tModeStatic}
}

func (l *lexer) scan() token {
	ch, _, err := l.buf.ReadRune()

	if nil != err {
		return createEndToken()
	}

	if l.mode == tModeIdentifier {
		l.mode = tModeStatic

		if isIdentifierRune(ch) {
			_ = l.buf.UnreadRune()
			return l.scanIdentifier()
		}
	}

	if isColon(ch) {
		return l.scanExpRegular()
	}

	if isCloseBrace(ch) {
		return createCloseVarToken()
	}

	// defaults to Static mode

	if isSlash(ch) {
		return createSlashToken()
	}

	if isOpenBrace(ch) {
		l.mode = tModeIdentifier
		return createOpenVarToken()
	}
	_ = l.buf.UnreadRune()
	return l.scanStatic()
}

func (l *lexer) scanStatic() token {
	var out bytes.Buffer

	for {
		ch, _, err := l.buf.ReadRune()

		if nil != err {
			break
		}

		if !isStatic(ch) {
			_ = l.buf.UnreadRune()
			break
		}

		out.WriteRune(ch)
	}

	return createStaticToken(out.String())
}

func (l *lexer) scanIdentifier() token {
	var out bytes.Buffer

	for {
		ch, _, err := l.buf.ReadRune()

		if nil != err {
			break
		}

		if !isIdentifierRune(ch) {
			_ = l.buf.UnreadRune()
			break
		}
		out.WriteRune(ch)
	}

	return createVarToken(out.String())
}

func (l *lexer) scanExpRegular() token {
	var out bytes.Buffer

	braces := 0
	for {
		ch, _, err := l.buf.ReadRune()

		if nil != err {
			break
		}

		if isOpenBrace(ch) {
			braces++
		}

		if isCloseBrace(ch) {
			if braces == 0 {
				_ = l.buf.UnreadRune()
				break
			}
			braces--
		}
		out.WriteRune(ch)
	}

	return createExpRegularToken(out.String())
}

func (l *lexer) scanAll() []token {

	tokens := make([]token, 0, 10)
	for {
		token := l.scan()
		tokens = append(tokens, token)
		if tEnd == token.t {
			break
		}
	}
	return tokens
}

func createSlashToken() token {
	return token{t: tSlash, v: "/"}
}

func createStaticToken(value string) token {
	return token{t: tStatic, v: value}
}

func createOpenVarToken() token {
	return token{t: tOpenVar, v: "{"}
}

func createVarToken(value string) token {
	return token{t: tVar, v: value}
}

func createExpRegularToken(value string) token {
	return token{t: tExpReg, v: value}
}

func createCloseVarToken() token {
	return token{t: tCloseVar, v: "}"}
}

func createEndToken() token {
	return token{t: tEnd, v: ""}
}

func isSlash(ch rune) bool {
	return '/' == ch
}

func isColon(ch rune) bool {
	return ':' == ch
}

func isOpenBrace(ch rune) bool {
	return '{' == ch
}

func isCloseBrace(ch rune) bool {
	return '}' == ch
}

func isAlpha(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9')
}
func isIdentifierRune(ch rune) bool {
	return isAlpha(ch) || (ch == '_')
}
func isStatic(ch rune) bool {
	return isAlpha(ch) ||
		(ch == '.') ||
		(ch == '-') ||
		(ch == '_')
}
