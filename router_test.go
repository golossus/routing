package routing

import (
	"net/http"
	"testing"
)

func TestTreeRouter(t *testing.T) {
	router := Router{}

	flag := false
	f := func(response http.ResponseWriter, request *http.Request) {
		flag = true
	}
	router.Register(http.MethodGet, "/path1", nil)
	router.Register(http.MethodGet, "/path2", f)

	request, _ := http.NewRequest("GET", "/path2", nil)
	router.ServeHTTP(nil, request)

	if !flag {
		t.Errorf("Handler not match ")
	}
}

func TestTreeRouterRegistrationOrder(t *testing.T) {
	router := Router{}

	flag := false
	flag2 := false
	f := func(response http.ResponseWriter, request *http.Request) {
		flag = true
	}
	f2 := func(response http.ResponseWriter, request *http.Request) {
		flag2 = true
	}
	router.Register(http.MethodGet, "/1/classes/{className}/{objectId}", f)
	router.Register(http.MethodGet, "/1/classes/{className}", f2)

	request, _ := http.NewRequest("GET", "/1/classes/:className/:objectId", nil)
	router.ServeHTTP(nil, request)

	if !flag {
		t.Errorf("Handler not match ")
	}
	
	request, _ = http.NewRequest("GET", "/1/classes/:className", nil)
	router.ServeHTTP(nil, request)

	if !flag2 {
		t.Errorf("Handler2 not match ")
	}
}

func TestTreeRouterMethod(t *testing.T) {
	router := Router{}

	flag := false
	flag2 := false
	f := func(response http.ResponseWriter, request *http.Request) {
		flag = true
	}
	f2 := func(response http.ResponseWriter, request *http.Request) {
		flag2 = true
	}
	router.Register(http.MethodGet, "/1/classes/{className}", f)
	router.Register(http.MethodPost, "/1/classes/{className}", f2)

	request, _ := http.NewRequest("GET", "/1/classes/:className", nil)
	router.ServeHTTP(nil, request)

	if !flag {
		t.Errorf("Handler not match ")
	}

	request, _ = http.NewRequest("POST", "/1/classes/:className", nil)
	router.ServeHTTP(nil, request)

	if !flag2 {
		t.Errorf("Handler2 not match ")
	}
}

func TestGetURLParamatersBagInHandler(t *testing.T) {
	router := Router{}

	f := func(response http.ResponseWriter, request *http.Request) {
		urlParameterBag := GetURLParameters(request)
		if 2 != len(urlParameterBag.params) {
			t.Errorf("")
		}
		id, _ := urlParameterBag.GetByName("id")
		if "100" != id {
			t.Errorf("")
		}

		name, _ := urlParameterBag.GetByName("name")
		if "dummy" != name {
			t.Errorf("")
		}
	}
	router.Register(http.MethodGet, "/path1/{id:[0-9]+}/{name:[a-z]{1,5}}", f)

	request, _ := http.NewRequest("GET", "/path1/100/dummy", nil)
	router.ServeHTTP(nil, request)

}

func TestGetURLParamatersFailsIfRegExpFails(t *testing.T) {

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	router := Router{}

	f := func(response http.ResponseWriter, request *http.Request) {}
	router.Register(http.MethodGet, "/path1/{id:[0-9]+}/{name:[a-z]+}", f)

	request, _ := http.NewRequest("GET", "/path1/100/123", nil)
	router.ServeHTTP(nil, request)

}

func TestVariosVerbsMatching(t *testing.T) {
	router := Router{}

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

	router.Register(http.MethodGet, "/path1", f1)
	router.Register(http.MethodGet, "/path1/{id}", f2)
	router.Register(http.MethodPost, "/path1", f3)
	router.Register(http.MethodPut, "/path1/{id}", f4)
	router.Register(http.MethodDelete, "/path1/{id}", f5)

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

func TestVerbsMethodsAreWorking(t *testing.T) {
	router := Router{}

	flag := 0
	f1 := func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			t.Errorf("")
		}
		flag++
	}
	f2 := func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodHead {
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
		if request.Method != http.MethodPatch {
			t.Errorf("")
		}
		flag++
	}
	f6 := func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodDelete {
			t.Errorf("")
		}
		flag++
	}
	f7 := func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodConnect {
			t.Errorf("")
		}
		flag++
	}
	f8 := func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodOptions {
			t.Errorf("")
		}
		flag++
	}
	f9 := func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodTrace {
			t.Errorf("")
		}
		flag++
	}
	router.Get("/path1", f1)
	router.Head("/path1", f2)
	router.Post("/path1", f3)
	router.Put("/path1", f4)
	router.Patch("/path1", f5)
	router.Delete("/path1", f6)
	router.Connect("/path1", f7)
	router.Options("/path1", f8)
	router.Trace("/path1", f9)

	request, _ := http.NewRequest(http.MethodGet, "/path1", nil)
	router.ServeHTTP(nil, request)

	request, _ = http.NewRequest(http.MethodHead, "/path1", nil)
	router.ServeHTTP(nil, request)

	request, _ = http.NewRequest(http.MethodPost, "/path1", nil)
	router.ServeHTTP(nil, request)

	request, _ = http.NewRequest(http.MethodPut, "/path1", nil)
	router.ServeHTTP(nil, request)

	request, _ = http.NewRequest(http.MethodPatch, "/path1", nil)
	router.ServeHTTP(nil, request)

	request, _ = http.NewRequest(http.MethodDelete, "/path1", nil)
	router.ServeHTTP(nil, request)

	request, _ = http.NewRequest(http.MethodConnect, "/path1", nil)
	router.ServeHTTP(nil, request)

	request, _ = http.NewRequest(http.MethodOptions, "/path1", nil)
	router.ServeHTTP(nil, request)

	request, _ = http.NewRequest(http.MethodTrace, "/path1", nil)
	router.ServeHTTP(nil, request)

	if flag != 9 {
		t.Errorf("")
	}
}
