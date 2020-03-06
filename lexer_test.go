package http_router

import "testing"

func TestSimplePath(t *testing.T) {
	lexer := newLexer("/path1")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: TSlash},
		{v: "path1", t: TStatic},
		{v: "", t: TEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestDoubleStatic(t *testing.T) {
	lexer := newLexer("/path1/path2")

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
	lexer := newLexer("/")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: TSlash},
		{v: "", t: TEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestEndSlash(t *testing.T) {
	lexer := newLexer("/path1/")

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
	lexer := newLexer("/path1/{id}/path2")

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
	lexer := newLexer("/path1/{id/path2")

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
	lexer := newLexer("/path1/id}/path2")

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
	lexer := newLexer("/path1/{}/path2")

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
	lexer := newLexer("/path1/{{id}/path2")

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
	lexer := newLexer("/path1/{id}/path2{")

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
	lexer := newLexer("/path1/{id}}/path2")

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
	lexer := newLexer("/path1/{id}_name/path2")

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

func TestWithExpRegular(t *testing.T) {
	lexer := newLexer("/path1/{id:[0-9]{4}-[0-9]{4}}")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: TSlash},
		{v: "path1", t: TStatic},
		{v: "/", t: TSlash},
		{v: "{", t: TOpenVar},
		{v: "id", t: TVar},
		{v: "[0-9]{4}-[0-9]{4}", t: TExpReg},
		{v: "}", t: TCloseVar},
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
