package http_router

import (
	"bytes"
	"fmt"
	"regexp"
)

const (
	TChunkStatic = iota
	TChunkDynamic
)

type token struct {
	v string
	t int
}

type parser struct {
	lexer  *lexer
	last   token
	chunks []chunk
	buf    bytes.Buffer
}

type chunk struct {
	t   int
	v   string
	exp *regexp.Regexp
}

func newParser(path string) *parser {
	l := newLexer(path)

	return &parser{lexer: l, chunks: make([]chunk, 0, 3)}
}

func (p *parser) parse() (bool, error) {
	return p.parseStart()
}

func (p *parser) parseStart() (bool, error) {
	token := p.lexer.scan()
	if !isSlashToken(token) {
		return false, fmt.Errorf("parser error, expected %s but got %s", "/", token.v)
	}
	p.last = token
	p.buf.Write([]byte(token.v))
	return p.parseStatic()
}

func (p *parser) parseVar() (bool, error) {
	token := p.lexer.scan()

	if !isVarToken(token) {
		return false, fmt.Errorf("parser error, expected %s but got %s", "var identifier", token.v)
	}
	p.buf.Write([]byte(token.v))

	token = p.lexer.scan()

	var regExp *regexp.Regexp
	if isRegExpressionToken(token) {
		regExp = regexp.MustCompile(fmt.Sprintf("^%s$", token.v))

		token = p.lexer.scan()
	}

	if !isCloseVarToken(token) {
		return false, fmt.Errorf("parser error, expected %s but got %s", "}", token.v)
	}
	p.chunks = append(p.chunks, chunk{t: TChunkDynamic, v: p.buf.String(), exp: regExp})
	p.buf.Reset()

	token = p.lexer.scan()
	if isEndToken(token) {
		return true, nil
	}
	if isOpenVarToken(token) || isCloseVarToken(token) || isRegExpressionToken(token) {
		return false, fmt.Errorf("parser error, expected %s but got %s", "static", token.v)
	}

	p.buf.Write([]byte(token.v))

	p.last = token
	return p.parseStatic()

}

func (p *parser) parseStatic() (bool, error) {
	token := p.lexer.scan()

	if isEndToken(token) {
		p.chunks = append(p.chunks, chunk{t: TChunkStatic, v: p.buf.String()})
		p.buf.Reset()
		return true, nil
	}

	if (!isSlashToken(p.last) && isSlashToken(token)) || isStaticToken(token) {
		p.buf.Write([]byte(token.v))
		p.last = token

		return p.parseStatic()
	}

	if isOpenVarToken(token) {
		p.chunks = append(p.chunks, chunk{t: TChunkStatic, v: p.buf.String()})
		p.buf.Reset()

		return p.parseVar()
	}

	return false, fmt.Errorf("parser error, unexpected token %s", token.v)
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

func isRegExpressionToken(t token) bool {
	return t.t == TExpReg
}

func isEndToken(t token) bool {
	return t.t == TEnd
}

func isStaticToken(t token) bool {
	return t.t == TStatic
}
