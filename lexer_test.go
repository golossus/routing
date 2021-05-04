package routing

import "testing"

func validateTokens(expectedTokens, tokens []token, t *testing.T) {
	for index, token := range tokens {
		if token.t != expectedTokens[index].t || token.v != expectedTokens[index].v {
			t.Errorf(
				"expected token: %v,%v but got %v,%v",
				expectedTokens[index].t,
				expectedTokens[index].v,
				token.t,
				token.v,
			)
		}
	}
}

func TestLexer_ScanAll_SimplePath(t *testing.T) {
	lexer := newLexer("/path1")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: tSlash},
		{v: "path1", t: tStatic},
		{v: "", t: tEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestLexer_ScanAll_DoubleStatic(t *testing.T) {
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

func TestLexer_ScanAll_Home(t *testing.T) {
	lexer := newLexer("/")

	tokens := lexer.scanAll()
	expectedTokens := []token{
		{v: "/", t: tSlash},
		{v: "", t: tEnd},
	}
	validateTokens(expectedTokens, tokens, t)
}

func TestLexer_ScanAll_EndSlash(t *testing.T) {
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

func TestLexer_ScanAll_WithParameter(t *testing.T) {
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

func TestLexer_ScanAll_InvalidRoute(t *testing.T) {
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

func TestLexer_ScanAll_InvalidRouteMissingOpen(t *testing.T) {
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

func TestLexer_ScanAll_WithEmptyVar(t *testing.T) {
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

func TestLexer_ScanAll_WithDoubleOpenVar(t *testing.T) {
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

func TestLexer_ScanAll_WithFinalDoubleOpenVar(t *testing.T) {
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

func TestLexer_ScanAll_WithDoubleCloseVar(t *testing.T) {
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

func TestLexer_ScanAll_WithStaticValueAfterCloseVar(t *testing.T) {
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

func TestLexer_ScanAll_WithExpRegular(t *testing.T) {
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
