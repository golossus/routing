package hw14_go

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
	TFinishVar
	TEnd
)

const (
	TModeStatic = iota
	TModeIdentifier
)


type token struct {
	v string
	t int
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

func createCloseVarToken() token {
	return token{t: TCloseVar, v: "}"}
}

func createFinishVarToken(value string) token {
	return token{t: TFinishVar, v: value}
}

func createEndToken() token {
	return token{t: TEnd, v: ""}
}


func isSlash(ch rune) bool {
	return '/' == ch
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
func isStatic(ch rune) bool {
	return isAlpha(ch) ||
		(ch == '.') ||
		(ch == '-') ||
		(ch == '_')
}

type Lexer struct {
	mode int
	buf *bufio.Reader
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

		if (isCloseBrace(ch)) {
			l.mode = TModeStatic
			return createCloseVarToken()
		}

		return l.scanIdentifier()

	}

	// defaults to Static mode

	if isSlash(ch) {
		return createSlashToken()
	}

	if isOpenBrace(ch) {
		l.mode = TModeIdentifier
		return createOpenVarToken()
	}

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
			l.buf.UnreadRune()
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

		if !isAlpha(ch) {
			l.buf.UnreadRune()
			break
		}

		// is valid

		out.WriteRune(ch)
	}

	return createVarToken(out.String())
}

