package routing

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRouteBuilder(t *testing.T) {
	r := NewRouter()
	h := func(writer http.ResponseWriter, request *http.Request) {}
	m := func(r *http.Request) bool { return true }

	builder := &routeBuilder{&r, route{}}

	builder.Method("GET")
	builder.Path("/some")
	builder.Handler(h)
	builder.Name("name")
	builder.Host("domain.com")
	builder.Schemas("http", "ftp")
	builder.Header("x-1", "1")
	builder.Header("x-2", "2")
	builder.QueryParam("p1", "v1")
	builder.QueryParam("p2", "v2")
	builder.Matcher(m)

	expected := route{
		method:  "GET",
		path:    "/some",
		handler: h,
		options: MatchingOptions{
			Name:    "name",
			Host:    "domain.com",
			Schemas: []string{"http", "ftp"},
			Headers: map[string]string{
				"x-1": "1",
				"x-2": "2",
			},
			QueryParams: map[string]string{
				"p1": "v1",
				"p2": "v2",
			},
			Custom: m,
		},
	}

	assertStringEqual(t, fmt.Sprintln(expected), fmt.Sprintln(builder.curr))
}
