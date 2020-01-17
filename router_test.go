package http_router

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
	router.AddHandler(http.MethodGet, "/path1", f)

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
	router.AddHandler(http.MethodGet, "/path1", nil)
	router.AddHandler(http.MethodGet, "/path2", f)

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
	router.AddHandler(http.MethodGet, "/path1", nil)
	router.AddHandler(http.MethodGet, "/path2", f)

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
	router.AddHandler(http.MethodGet, "/path1", f)

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
	router.AddHandler(http.MethodGet, "/path1", nil)
	router.AddHandler(http.MethodGet, "/path2", f)

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
	router.AddHandler(http.MethodGet, "/path1", nil)
	router.AddHandler(http.MethodGet, "/path2", f)

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
	router.AddHandler(http.MethodGet, "/path1", nil)
	router.AddHandler(http.MethodGet, "/path2", f)

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

func TestGetURLParamatersBagInHandler(t *testing.T) {
	router := PrefixTreeRouter{}

	f := func(response http.ResponseWriter, request *http.Request) {
		urlParameterBag := request.Context().Value(ParamsBagKey).(UrlParameterBag)
		if 2 != len(urlParameterBag.params) {
			t.Errorf("")
		}
		id := urlParameterBag.GetByName("id", "0")
		if "100" != id {
			t.Errorf("")
		}

		name := urlParameterBag.GetByName("name", "")
		if "dummy" != name {
			t.Errorf("")
		}
	}
	router.AddHandler(http.MethodGet, "/path1/{id}/{name}", f)

	request, _ := http.NewRequest("GET", "/path1/100/dummy", nil)
	router.ServeHTTP(nil, request)

}

func TestVariosVerbsMatching(t *testing.T) {
	router := PrefixTreeRouter{}

	flag := 0
	f1 := func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			t.Errorf("")
		}
		flag++
	}

	f2 := func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			t.Errorf("")
		}
		flag++
	}
	f3 := func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			t.Errorf("")
		}
		flag++
	}
	f4 := func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPut {
			t.Errorf("")
		}
		flag++
	}
	f5 := func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodDelete {
			t.Errorf("")
		}
		flag++
	}

	router.AddHandler(http.MethodGet, "/path1", f1)
	router.AddHandler(http.MethodGet, "/path1/{id}", f2)
	router.AddHandler(http.MethodPost, "/path1", f3)
	router.AddHandler(http.MethodPut, "/path1/{id}", f4)
	router.AddHandler(http.MethodDelete, "/path1/{id}", f5)

	request, _ := http.NewRequest(http.MethodGet, "/path1", nil)
	router.ServeHTTP(nil, request)

	request, _ = http.NewRequest(http.MethodGet, "/path1/100", nil)
	router.ServeHTTP(nil, request)

	request, _ = http.NewRequest(http.MethodPost, "/path1", nil)
	router.ServeHTTP(nil, request)

	request, _ = http.NewRequest(http.MethodPut, "/path1/100", nil)
	router.ServeHTTP(nil, request)

	request, _ = http.NewRequest(http.MethodDelete, "/path1/100", nil)
	router.ServeHTTP(nil, request)

	if flag != 5 {
		t.Errorf("")
	}

}
