package http_router

import (
	"fmt"
	"testing"
)

func TestSimplePath(t *testing.T) {
	lexer := NewLexer("/path1")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: TSlash},
		{v: "path1", t: TStatic},
		{v: "", t: TEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestDoubleStatic(t *testing.T) {
	lexer := NewLexer("/path1/path2")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: TSlash},
		{v: "path1", t: TStatic},
		{v: "/", t: TSlash},
		{v: "path2", t: TStatic},
		{v: "", t: TEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestHome(t *testing.T) {
	lexer := NewLexer("/")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: TSlash},
		{v: "", t: TEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestEndSlash(t *testing.T) {
	lexer := NewLexer("/path1/")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: TSlash},
		{v: "path1", t: TStatic},
		{v: "/", t: TSlash},
		{v: "", t: TEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestWithParameter(t *testing.T) {
	lexer := NewLexer("/path1/{id}/path2")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: TSlash},
		{v: "path1", t: TStatic},
		{v: "/", t: TSlash},
		{v: "{", t: TOpenVar},
		{v: "id", t: TVar},
		{v: "}", t: TCloseVar},
		{v: "/", t: TSlash},
		{v: "path2", t: TStatic},
		{v: "", t: TEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestInvalidRoute(t *testing.T) {
	lexer := NewLexer("/path1/{id/path2")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: TSlash},
		{v: "path1", t: TStatic},
		{v: "/", t: TSlash},
		{v: "{", t: TOpenVar},
		{v: "id", t: TVar},
		{v: "/", t: TSlash},
		{v: "path2", t: TStatic},
		{v: "", t: TEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestInvalidRouteMissingOpen(t *testing.T) {
	lexer := NewLexer("/path1/id}/path2")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: TSlash},
		{v: "path1", t: TStatic},
		{v: "/", t: TSlash},
		{v: "id", t: TStatic},
		{v: "}", t: TCloseVar},
		{v: "/", t: TSlash},
		{v: "path2", t: TStatic},
		{v: "", t: TEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestWithEmptyVar(t *testing.T) {
	lexer := NewLexer("/path1/{}/path2")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: TSlash},
		{v: "path1", t: TStatic},
		{v: "/", t: TSlash},
		{v: "{", t: TOpenVar},
		{v: "}", t: TCloseVar},
		{v: "/", t: TSlash},
		{v: "path2", t: TStatic},
		{v: "", t: TEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestWithDoubleOpenVar(t *testing.T) {
	lexer := NewLexer("/path1/{{id}/path2")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: TSlash},
		{v: "path1", t: TStatic},
		{v: "/", t: TSlash},
		{v: "{", t: TOpenVar},
		{v: "{", t: TOpenVar},
		{v: "id", t: TVar},
		{v: "}", t: TCloseVar},
		{v: "/", t: TSlash},
		{v: "path2", t: TStatic},
		{v: "", t: TEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestWithFinalDoubleOpenVar(t *testing.T) {
	lexer := NewLexer("/path1/{id}/path2{")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: TSlash},
		{v: "path1", t: TStatic},
		{v: "/", t: TSlash},
		{v: "{", t: TOpenVar},
		{v: "id", t: TVar},
		{v: "}", t: TCloseVar},
		{v: "/", t: TSlash},
		{v: "path2", t: TStatic},
		{v: "{", t: TOpenVar},
		{v: "", t: TEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestWithDoubleCloseVar(t *testing.T) {
	lexer := NewLexer("/path1/{id}}/path2")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: TSlash},
		{v: "path1", t: TStatic},
		{v: "/", t: TSlash},
		{v: "{", t: TOpenVar},
		{v: "id", t: TVar},
		{v: "}", t: TCloseVar},
		{v: "}", t: TCloseVar},
		{v: "/", t: TSlash},
		{v: "path2", t: TStatic},
		{v: "", t: TEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestWithStaticValueAfterCloseVar(t *testing.T) {
	lexer := NewLexer("/path1/{id}_name/path2")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: TSlash},
		{v: "path1", t: TStatic},
		{v: "/", t: TSlash},
		{v: "{", t: TOpenVar},
		{v: "id", t: TVar},
		{v: "}", t: TCloseVar},
		{v: "_name", t: TStatic},
		{v: "/", t: TSlash},
		{v: "path2", t: TStatic},
		{v: "", t: TEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func validateTokens(expectedTokens, tokens []token, t *testing.T) {
	for index, token := range tokens {
		if token.t != expectedTokens[index].t || token.v != expectedTokens[index].v {
			t.Errorf(
				"Expected token: %v,%v but got %v,%v",
				expectedTokens[index].t,
				expectedTokens[index].v,
				token.t,
				token.v,
			)
		}
	}
}

func TestParserValidatesValidPaths(t *testing.T) {
	paths := []string{
		"/",
		"/path1",
		"/path1/path2",
		"/path1/path2/",
		"/path1/{id}",
		"/path1/{id}/",
		"/path1/{id}/path2",
		"/{id}",
		"/{id}/",
		"/{id}/path1",
		"/asdf-{id}",
		"/{id}-adfasf",
		"/{id}/{name}",
		"/{id}/{name}/",
		"/{id}-{name}/",
	}

	for _, path := range paths {
		parser := NewParser(path)
		valid, err := parser.parse()
		if !valid {
			t.Errorf("%v in path %v", err, path)
		}
	}
}

func TestParserDoesNotValidateInvalidPaths(t *testing.T) {
	paths := []string{
		"",
		"//",
		"/path1//path2",
		"/path1//",
		"path1/path2/",
		"/path1/{id}}",
		"/path1/{{id}",
		"/path1/id}",
		"/path1/{id",
		"/path1/{id{",
		"/path1/{id/",
		"/path1/{}",
		"/path1/{id}{name}",
		"/{}",
		"/{",
		"/}",
		"/{name{id}}",
		"/{id$}",
		"/{.id}",
		"/{id/name}",
	}

	for _, path := range paths {
		parser := NewParser(path)
		valid, _ := parser.parse()
		if valid {
			t.Errorf("Validated invalid path %v", path)
		}
	}
}

func TestParserChechChunks(t *testing.T) {
	paths := []string{
		"/",
		"/path1",
		"/path1/path2",
		"/path1/path2/",
		"/path1/{id}",
		"/path1/{id}/",
		"/path1/{id}/path2",
		"/{id}",
		"/{id}/",
		"/{id}/path1",
		"/asdf-{id}",
		"/{id}-adfasf",
		"/{id}/{name}",
		"/{id}/{name}/",
		"/{id}-{name}/",
	}

	expectedChunks := []string{
		"[{0 /}]",
		"[{0 /path1}]",
		"[{0 /path1/path2}]",
		"[{0 /path1/path2/}]",
		"[{0 /path1/} {1 id}]",
		"[{0 /path1/} {1 id} {0 /}]",
		"[{0 /path1/} {1 id} {0 /path2}]",
		"[{0 /} {1 id}]",
		"[{0 /} {1 id} {0 /}]",
		"[{0 /} {1 id} {0 /path1}]",
		"[{0 /asdf-} {1 id}]",
		"[{0 /} {1 id} {0 -adfasf}]",
		"[{0 /} {1 id} {0 /} {1 name}]",
		"[{0 /} {1 id} {0 /} {1 name} {0 /}]",
		"[{0 /} {1 id} {0 -} {1 name} {0 /}]",
	}

	for index, path := range paths {
		parser := NewParser(path)
		valid, err := parser.parse()
		if !valid {
			t.Errorf("%v in path %v", err, path)
		}
		if expectedChunks[index] != fmt.Sprintf("%v", parser.chunks) {
			t.Errorf("Parser error, Invalid chunk %v for path %v", parser.chunks, path)
		}
	}
}
