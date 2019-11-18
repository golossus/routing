package hw14_go

import (
	"net/http"
	"testing"
)

func TestIsNilOnCreation(t *testing.T) {
	router := Router{}
	if nil != router.handlers {
		t.Errorf("Router is not empty on creation")
	}
}

func TestAddOneRouteHandler(t *testing.T) {
	router := Router{}
	flag := false
	f := func(response http.ResponseWriter, request *http.Request) {
		flag = true
	}
	router.AddHandler("/path1", f)

	if nil == router.handlers {
		t.Errorf("Router is empty")
	}

	if "/path1" != router.handlers[0].path {
		t.Errorf("Path not valid")
	}

	if nil == router.handlers[0].handler {
		t.Errorf("Handler not valid")
	}

	router.handlers[0].handler(nil, nil)
	if !flag {
		t.Errorf("Handler not match ")
	}
}

func TestAddSeveralRoutesHandler(t *testing.T) {
	router := Router{}
	flag := false
	f := func(response http.ResponseWriter, request *http.Request) {
		flag = true
	}
	router.AddHandler("/path1", nil)
	router.AddHandler("/path2", f)

	if 2 != len(router.handlers) {
		t.Errorf("Invalid size")
	}

	if "/path1" != router.handlers[0].path {
		t.Errorf("Path not valid")
	}

	if nil != router.handlers[0].handler {
		t.Errorf("Handler not valid")
	}

	if "/path2" != router.handlers[1].path {
		t.Errorf("Path not valid")
	}

	if nil == router.handlers[1].handler {
		t.Errorf("Handler not valid")
	}

	router.handlers[1].handler(nil, nil)
	if !flag {
		t.Errorf("Handler not match ")
	}
}

func TestRouteMatch(t *testing.T) {
	router := Router{}
	flag := false
	f := func(response http.ResponseWriter, request *http.Request) {
		flag = true
	}
	router.AddHandler("/path1", nil)
	router.AddHandler("/path2", f)

	request, _ := http.NewRequest("GET", "/path2", nil)
	router.ServeHTTP(nil, request)

	if !flag {
		t.Errorf("Handler not match ")
	}
}
