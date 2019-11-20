package hw14_go

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

const (
	TSlash = iota
	TStatic
	TOpenVar
	TVar
	TCloseVar
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

		if isAlpha(ch) {
			l.buf.UnreadRune()
			return l.scanIdentifier()
		}
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
	l.buf.UnreadRune()
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
		out.WriteRune(ch)
	}

	return createVarToken(out.String())
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

type Parser struct {
	lexer *Lexer
	last  token
}

func NewParser(path string) *Parser {
	l := NewLexer(path)

	return &Parser{lexer: l}
}

func isSlashToken(t token) bool {
	return t.t == TSlash
}

func isOpenVarToken(t token) bool {
	return t.t == TOpenVar
}

func isCloseVarToken(t token) bool {
	return t.t == TCloseVar
}

func isVarToken(t token) bool {
	return t.t == TVar
}

func isEndToken(t token) bool {
	return t.t == TEnd
}

func isStaticToken(t token) bool {
	return t.t == TStatic
}

//By now we only validate the path with the parser. In the
//future the goal is to return a list of things to use with
//the router.
func (p *Parser) parse() (bool, error) {
	return p.parseStart()
}

func (p *Parser) parseStart() (bool, error) {
	token := p.lexer.scan()
	if !isSlashToken(token) {
		return false, fmt.Errorf("Parser error, expected %s but got %s", "/", token.v)
	}

	token = p.lexer.scan()
	p.last = token
	switch {
	case isOpenVarToken(token):
		return p.parseVar()
	case isEndToken(token):
		return true, nil
	case isSlashToken(token) || isCloseVarToken(token):
		return false, fmt.Errorf("Parser error, expected %s but got %s", "alphanum", token.v)
	}

	return p.parseStatic()
}

func (p *Parser) parseVar() (bool, error) {
	token := p.lexer.scan()

	if !isVarToken(token) {
		return false, fmt.Errorf("Parser error, expected %s but got %s", "var identifier", token.v)
	}

	token = p.lexer.scan()
	if !isCloseVarToken(token) {
		return false, fmt.Errorf("Parser error, expected %s but got %s", "}", token.v)
	}

	token = p.lexer.scan()
	if isEndToken(token) {
		return true, nil
	}
	if isOpenVarToken(token) || isCloseVarToken(token) {
		return false, fmt.Errorf("Parser error, expected %s but got %s", "static", token.v)
	}

	p.last = token

	return p.parseStatic()

}

func (p *Parser) parseStatic() (bool, error) {
	token := p.lexer.scan()

	if isEndToken(token) {
		return true, nil
	}

	if isSlashToken(token) && isSlashToken(p.last) {
		return false, fmt.Errorf("Parser error, unexpected %v", token.v)
	}

	if isSlashToken(token) || isStaticToken(token) {
		p.last = token
		return p.parseStatic()
	}

	if isCloseVarToken(token) {
		return false, fmt.Errorf("Parser error, expected %s but got %s", "}", token.v)
	}

	if isOpenVarToken(token) {
		p.last = token
		return p.parseVar()
	}

	return false, fmt.Errorf("Parser error, unexpected token %s", token.v)
}
