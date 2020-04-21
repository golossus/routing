package routing

import (
	"fmt"
	"strconv"
	"testing"
)

func TestParserValidatesValidPaths(t *testing.T) {
	paths := []string{
		"/",
		"/path1",
		"/path1/path2",
		"/path1/path2/",
		"/path1/{id}",
		"/path1/{id}/",
		"/path1/{id:[0-9]{4}-[0-9]{2}-[0-9]{2}}/",
		"/path1/{id:/ab+c/}",
		"/path1/{id}/path2",
		"/path1/{id:/ab+c/}/path2",
		"/{id}",
		"/{id}/",
		"/{id}/path1",
		"/asdf-{id}",
		"/{id}-adfasf",
		"/{id}/{name}",
		"/{id}/{name}/",
		"/{id}-{name}/",
		"/{id:[0-9]+}-{name:/ab+c/}/",
	}

	for _, path := range paths {
		parser := newParser(path)
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
		"/{/}",
		"/{",
		"/}",
		"/{name{id}}",
		"/{id$}",
		"/{.id}",
		"/{id/name}",
		":[0-9]+",
		"/:[0-9]+",
		"/path1:[0-9]+/{id}",
		"/path1/:[0-9]+{id}",
		"/path1/{:[0-9]+id}",
		"/path1/{id}:[0-9]+",
	}

	for _, path := range paths {
		parser := newParser(path)
		valid, _ := parser.parse()
		if valid {
			t.Errorf("Validated invalid path %v", path)
		}
	}
}

func TestParserReturnsErrorWhenExpressionIsInvalid(t *testing.T) {
	paths := []string{
		"/path1/{id:[0-9+}/",
	}

	for _, path := range paths {
		parser := newParser(path)
		_, err := parser.parse()

		if  err == nil {
			t.Errorf("The code did not return an error")
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
		"/{id:[0-9]+}/",
		"/path1/{id:/ab+c/}",
		"/path1/{id:/ab+c/}/path2",
		"/{id:[0-9]+}-{name:/ab+c/}/",
	}

	expectedChunks := []string{
		"[{0 / }]",
		"[{0 /path1 }]",
		"[{0 /path1/path2 }]",
		"[{0 /path1/path2/ }]",
		"[{0 /path1/ } {1 id }]",
		"[{0 /path1/ } {1 id } {0 / }]",
		"[{0 /path1/ } {1 id } {0 /path2 }]",
		"[{0 / } {1 id }]",
		"[{0 / } {1 id } {0 / }]",
		"[{0 / } {1 id } {0 /path1 }]",
		"[{0 /asdf- } {1 id }]",
		"[{0 / } {1 id } {0 -adfasf }]",
		"[{0 / } {1 id } {0 / } {1 name }]",
		"[{0 / } {1 id } {0 / } {1 name } {0 / }]",
		"[{0 / } {1 id } {0 - } {1 name } {0 / }]",
		"[{0 / } {1 id ^[0-9]+$} {0 / }]",
		"[{0 /path1/ } {1 id ^/ab+c/$}]",
		"[{0 /path1/ } {1 id ^/ab+c/$} {0 /path2 }]",
		"[{0 / } {1 id ^[0-9]+$} {0 - } {1 name ^/ab+c/$} {0 / }]",
	}

	for index, path := range paths {
		parser := newParser(path)
		valid, err := parser.parse()
		if !valid {
			t.Errorf("%v in path %v", err, path)
		}
		var chunks []string
		for _, chunk := range parser.chunks {
			expString := ""
			if chunk.exp != nil {
				expString = chunk.exp.String()
			}
			chunks = append(chunks, "{"+strconv.Itoa(chunk.t)+" "+chunk.v+" "+expString+"}")
		}
		if expectedChunks[index] != fmt.Sprintf("%v", chunks) {
			t.Errorf("parser error, Invalid chunk %v for path %v", chunks, path)
		}
	}
}
