package loaders

import (
	. "github.com/golossus/routing"
	"reflect"
	"testing"
)

func TestJsonFileLoader_LoadFile(t *testing.T) {

	loader := JsonFileLoader{}
	err := loader.FromFile("../fixtures/routes.json", "../fixtures/routes2.json", "../fixtures/routes3.json")
	if err != nil {
		t.Error(err)
	}

	routes := loader.Load()
	if len(routes) != 3 {
		t.Errorf("routes length doesn't match")
	}

	expected := RouteDef{
		Method:  "GET",
		Handler: "get.users.handler",
		Path:    "/users",
		Options: RouteDefOptions{
			Name: "get.users",
		},
	}
	if !reflect.DeepEqual(routes[0], expected) {
		t.Errorf("route %v not equals to %v", routes[0], expected)
	}

	expected = RouteDef{
		Method:  "POST",
		Handler: "post.users.handler",
		Path:    "/users",
		Options: RouteDefOptions{
			Name: "post.users",
		},
	}
	if !reflect.DeepEqual(routes[1], expected) {
		t.Errorf("route %v not equals to %v", routes[1], expected)
	}

	expected = RouteDef{
		Method:  "PUT",
		Handler: "put.users.handler",
		Path:    "/users",
		Options: RouteDefOptions{
			Name:          "put.users",
			Host:          "my.domain.com",
			Schemas:       []string{"http", "https"},
			Headers:       map[string]string{"X-Dummy": "dummy"},
			QueryParams:   map[string]string{"offset": "2"},
			CustomMatcher: "",
		},
	}
	if !reflect.DeepEqual(routes[2], expected) {
		t.Errorf("route %v not equals to %v", routes[2], expected)
	}
}
