package hw14_go

import (
	"net/http"
	"testing"
)

func TestSliceIsNilOnCreation(t *testing.T) {
	router := SliceRouter{}
	if nil != router.handlers {
		t.Errorf("SliceRouter is not empty on creation")
	}
}

func TestSliceAddOneRouteHandler(t *testing.T) {
	router := SliceRouter{}
	flag := false
	f := func(response http.ResponseWriter, request *http.Request) {
		flag = true
	}
	router.AddHandler("/path1", f)

	if nil == router.handlers {
		t.Errorf("SliceRouter is empty")
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

func TestSliceAddSeveralRoutesHandler(t *testing.T) {
	router := SliceRouter{}
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

func TestSliceRouteMatch(t *testing.T) {
	router := SliceRouter{}
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

func TestMapIsNilOnCreation(t *testing.T) {
	router := MapRouter{}
	if nil != router.handlers {
		t.Errorf("MapRouter is not empty on creation")
	}
}

func TestMapAddOneRouteHandler(t *testing.T) {
	router := MapRouter{}
	flag := false
	f := func(response http.ResponseWriter, request *http.Request) {
		flag = true
	}
	router.AddHandler("/path1", f)

	if nil == router.handlers {
		t.Errorf("MapRouter is empty")
	}

	handler, found := router.handlers["/path1"]
	if !found {
		t.Errorf("Path not valid")
	}

	if nil == handler {
		t.Errorf("Handler not valid")
	}

	handler(nil, nil)
	if !flag {
		t.Errorf("Handler not match ")
	}
}

func TestMapAddSeveralRoutesHandler(t *testing.T) {
	router := MapRouter{}
	flag := false
	f := func(response http.ResponseWriter, request *http.Request) {
		flag = true
	}
	router.AddHandler("/path1", nil)
	router.AddHandler("/path2", f)

	if 2 != len(router.handlers) {
		t.Errorf("Invalid size")
	}

	handler, found := router.handlers["/path1"]
	if !found {
		t.Errorf("Path not valid")
	}

	if nil != handler {
		t.Errorf("Handler not valid")
	}

	handler, found = router.handlers["/path2"]
	if !found {
		t.Errorf("Path not valid")
	}

	if nil == handler {
		t.Errorf("Handler not valid")
	}

	handler(nil, nil)
	if !flag {
		t.Errorf("Handler not match ")
	}

}

func TestMapRouteMatch(t *testing.T) {
	router := MapRouter{}
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

func TestPrefixTreeRouter(t *testing.T) {
	router := PrefixTreeRouter{}

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

func TestUrlParameterBagEmptyOnCreation(t *testing.T) {
	bag := NewUrlParameterBag()

	if bag.params != nil {
		t.Errorf("")
	}
}

func TestUrlParameterBagAddsParameter(t *testing.T) {
	bag := NewUrlParameterBag()

	bag.addParameter(urlParameter{name: "param1", value: "v1"})

	if bag.params == nil {
		t.Errorf("")
	}

	if len(bag.params) == 0 {
		t.Errorf("")
	}

	if bag.params[0].name != "param1" {
		t.Errorf("")
	}
}

func TestUrlParameterBagGetByName(t *testing.T) {
	bag := NewUrlParameterBag()

	bag.addParameter(urlParameter{name: "param1", value: "v1"})
	bag.addParameter(urlParameter{name: "param2", value: "v2"})
	bag.addParameter(urlParameter{name: "param3", value: "v3"})

	if bag.params == nil {
		t.Errorf("")
	}

	if len(bag.params) != 3 {
		t.Errorf("")
	}

	if bag.GetByName("param1", "") != "v1" {
		t.Errorf("")
	}

	if bag.GetByName("param2", "") != "v2" {
		t.Errorf("")
	}

	if bag.GetByName("param3", "") != "v3" {
		t.Errorf("")
	}

	if bag.GetByName("param4", "v4") != "v4" {
		t.Errorf("")
	}
}

func TestUrlParameterBagGetByIndex(t *testing.T) {
	bag := NewUrlParameterBag()

	bag.addParameter(urlParameter{name: "param1", value: "v1"})
	bag.addParameter(urlParameter{name: "param2", value: "v2"})
	bag.addParameter(urlParameter{name: "param3", value: "v3"})

	if bag.params == nil {
		t.Errorf("")
	}

	if len(bag.params) != 3 {
		t.Errorf("")
	}

	if bag.GetByIndex(0, "") != "v1" {
		t.Errorf("")
	}

	if bag.GetByIndex(1, "") != "v2" {
		t.Errorf("")
	}

	if bag.GetByIndex(2, "") != "v3" {
		t.Errorf("")
	}

	if bag.GetByIndex(3, "v4") != "v4" {
		t.Errorf("")
	}
}
