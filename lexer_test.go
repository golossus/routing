package routing

import "testing"

func TestSimplePath(t *testing.T) {
	lexer := newLexer("/path1")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: tSlash},
		{v: "path1", t: tStatic},
		{v: "", t: tEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestDoubleStatic(t *testing.T) {
	lexer := newLexer("/path1/path2")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: tSlash},
		{v: "path1", t: tStatic},
		{v: "/", t: tSlash},
		{v: "path2", t: tStatic},
		{v: "", t: tEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestHome(t *testing.T) {
	lexer := newLexer("/")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: tSlash},
		{v: "", t: tEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestEndSlash(t *testing.T) {
	lexer := newLexer("/path1/")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: tSlash},
		{v: "path1", t: tStatic},
		{v: "/", t: tSlash},
		{v: "", t: tEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestWithParameter(t *testing.T) {
	lexer := newLexer("/path1/{id}/path2")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: tSlash},
		{v: "path1", t: tStatic},
		{v: "/", t: tSlash},
		{v: "{", t: tOpenVar},
		{v: "id", t: tVar},
		{v: "}", t: tCloseVar},
		{v: "/", t: tSlash},
		{v: "path2", t: tStatic},
		{v: "", t: tEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestInvalidRoute(t *testing.T) {
	lexer := newLexer("/path1/{id/path2")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: tSlash},
		{v: "path1", t: tStatic},
		{v: "/", t: tSlash},
		{v: "{", t: tOpenVar},
		{v: "id", t: tVar},
		{v: "/", t: tSlash},
		{v: "path2", t: tStatic},
		{v: "", t: tEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestInvalidRouteMissingOpen(t *testing.T) {
	lexer := newLexer("/path1/id}/path2")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: tSlash},
		{v: "path1", t: tStatic},
		{v: "/", t: tSlash},
		{v: "id", t: tStatic},
		{v: "}", t: tCloseVar},
		{v: "/", t: tSlash},
		{v: "path2", t: tStatic},
		{v: "", t: tEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestWithEmptyVar(t *testing.T) {
	lexer := newLexer("/path1/{}/path2")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: tSlash},
		{v: "path1", t: tStatic},
		{v: "/", t: tSlash},
		{v: "{", t: tOpenVar},
		{v: "}", t: tCloseVar},
		{v: "/", t: tSlash},
		{v: "path2", t: tStatic},
		{v: "", t: tEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestWithDoubleOpenVar(t *testing.T) {
	lexer := newLexer("/path1/{{id}/path2")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: tSlash},
		{v: "path1", t: tStatic},
		{v: "/", t: tSlash},
		{v: "{", t: tOpenVar},
		{v: "{", t: tOpenVar},
		{v: "id", t: tVar},
		{v: "}", t: tCloseVar},
		{v: "/", t: tSlash},
		{v: "path2", t: tStatic},
		{v: "", t: tEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestWithFinalDoubleOpenVar(t *testing.T) {
	lexer := newLexer("/path1/{id}/path2{")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: tSlash},
		{v: "path1", t: tStatic},
		{v: "/", t: tSlash},
		{v: "{", t: tOpenVar},
		{v: "id", t: tVar},
		{v: "}", t: tCloseVar},
		{v: "/", t: tSlash},
		{v: "path2", t: tStatic},
		{v: "{", t: tOpenVar},
		{v: "", t: tEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestWithDoubleCloseVar(t *testing.T) {
	lexer := newLexer("/path1/{id}}/path2")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: tSlash},
		{v: "path1", t: tStatic},
		{v: "/", t: tSlash},
		{v: "{", t: tOpenVar},
		{v: "id", t: tVar},
		{v: "}", t: tCloseVar},
		{v: "}", t: tCloseVar},
		{v: "/", t: tSlash},
		{v: "path2", t: tStatic},
		{v: "", t: tEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestWithStaticValueAfterCloseVar(t *testing.T) {
	lexer := newLexer("/path1/{id}_name/path2")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: tSlash},
		{v: "path1", t: tStatic},
		{v: "/", t: tSlash},
		{v: "{", t: tOpenVar},
		{v: "id", t: tVar},
		{v: "}", t: tCloseVar},
		{v: "_name", t: tStatic},
		{v: "/", t: tSlash},
		{v: "path2", t: tStatic},
		{v: "", t: tEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestWithExpRegular(t *testing.T) {
	lexer := newLexer("/path1/{id:[0-9]{4}-[0-9]{4}}")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: tSlash},
		{v: "path1", t: tStatic},
		{v: "/", t: tSlash},
		{v: "{", t: tOpenVar},
		{v: "id", t: tVar},
		{v: "[0-9]{4}-[0-9]{4}", t: tExpReg},
		{v: "}", t: tCloseVar},
		{v: "", t: tEnd},
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
