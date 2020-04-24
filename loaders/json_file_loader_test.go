package loaders

import (
	. "github.com/golossus/routing"
	"reflect"
	"testing"
)

func TestJsonFileLoader_LoadFile(t *testing.T) {

	loader := JsonFileLoader{}
	err := loader.FromFile("../fixtures/routes.json", "../fixtures/routes2.json")
	if err != nil {
		t.Error(err)
	}

	routes := loader.Load()
	if len(routes) != 2 {
		t.Errorf("routes length doesn't match")
	}

	expected := RouteDef{
		Method:  "GET",
		Handler: "get.users.handler",
		Schema:  "/users",
		Name:    "get.users",
	}
	if !reflect.DeepEqual(routes[0], expected) {
		t.Errorf("route %v not equals to %v", routes[0], expected)
	}

	expected = RouteDef{
		Method:  "POST",
		Handler: "post.users.handler",
		Schema:  "/users",
		Name:    "post.users",
	}
	if !reflect.DeepEqual(routes[1], expected) {
		t.Errorf("route %v not equals to %v", routes[1], expected)
	}
}
