package http_router

import (
	"bufio"
	"bytes"
	"strings"
)

const (
	TSlash = iota
	TStatic
	TOpenVar
	TVar
	TCloseVar
	TEnd
	TExpReg
)

const (
	TModeStatic = iota
	TModeIdentifier
)

type Lexer struct {
	mode int
	buf  *bufio.Reader
}

func NewLexer(path string) *Lexer {
	reader := strings.NewReader(path)
	return &Lexer{buf: bufio.NewReader(reader), mode: TModeStatic}
}

func (l *Lexer) scan() token {
	ch, _, err := l.buf.ReadRune()

	if nil != err {
		return createEndToken()
	}

	if l.mode == TModeIdentifier {
		l.mode = TModeStatic

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
		l.mode = TModeIdentifier
		return createOpenVarToken()
	}
	_ = l.buf.UnreadRune()
	return l.scanStatic()
}

func (l *Lexer) scanStatic() token {
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

func (l *Lexer) scanIdentifier() token {
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

func (l *Lexer) scanExpRegular() token {
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

func (l *Lexer) scanAll() []token {

	tokens := make([]token, 0, 10)
	for {
		token := l.scan()
		tokens = append(tokens, token)
		if TEnd == token.t {
			break
		}
	}
	return tokens
}

func createSlashToken() token {
	return token{t: TSlash, v: "/"}
}

func createStaticToken(value string) token {
	return token{t: TStatic, v: value}
}

func createOpenVarToken() token {
	return token{t: TOpenVar, v: "{"}
}

func createVarToken(value string) token {
	return token{t: TVar, v: value}
}

func createExpRegularToken(value string) token {
	return token{t: TExpReg, v: value}
}

func createCloseVarToken() token {
	return token{t: TCloseVar, v: "}"}
}

func createEndToken() token {
	return token{t: TEnd, v: ""}
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
