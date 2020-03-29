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
	router.addHandler(http.MethodGet, "/path1", nil)
	router.addHandler(http.MethodGet, "/path2", f)

	request, _ := http.NewRequest("GET", "/path2", nil)
	router.ServeHTTP(nil, request)

	if !flag {
		t.Errorf("Handler not match ")
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
	router.addHandler(http.MethodGet, "/path1/{id:[0-9]+}/{name:[a-z]{1,5}}", f)

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
	router.addHandler(http.MethodGet, "/path1/{id:[0-9]+}/{name:[a-z]+}", f)

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

	router.addHandler(http.MethodGet, "/path1", f1)
	router.addHandler(http.MethodGet, "/path1/{id}", f2)
	router.addHandler(http.MethodPost, "/path1", f3)
	router.addHandler(http.MethodPut, "/path1/{id}", f4)
	router.addHandler(http.MethodDelete, "/path1/{id}", f5)

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
