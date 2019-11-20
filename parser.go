package hw14_go

import (
	"bufio"
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

type Lexer struct {
	buf *bufio.Reader
}

func NewLexer(path string) *Lexer {
	reader := strings.NewReader(path)
	return Lexer{buf: bufio.Reader{rd: reader}}
}
