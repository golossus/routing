package routing

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testHandlerFunc = func(response http.ResponseWriter, request *http.Request) {
	fmt.Fprint(response, request.URL.Path)
}

func assertPathFound(t *testing.T, router Router, method, path string) {
	r, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 200 || w.Body.String() != path {
		t.Errorf("%s %s not found", method, path)
	}
}

func assertPathNotFound(t *testing.T, router Router, method, path string) {
	r, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 404 {
		t.Errorf("%s %s not found", method, path)
	}
}

func assertRequestHasParameterHandler(t *testing.T, bag URLParameterBag) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		par := GetURLParameters(r)

		if len(bag.params) != len(par.params) {
			t.Errorf("Size of parameter bag doesn't match %d != %d", len(bag.params), len(par.params))
		}

		for i := 0; i < len(bag.params); i++ {
			bagValue, _ := bag.GetByIndex(uint(i))
			parValue, _ := par.GetByIndex(uint(i))
			if bagValue != parValue {
				t.Errorf("Parameter at index %d don't match", i)
			}
		}

		testHandlerFunc(w, r)
	}
}

func TestTreeRouterFindsPaths(t *testing.T) {
	router := Router{}

	_ = router.Register(http.MethodGet, "/path1", testHandlerFunc)
	_ = router.Register(http.MethodGet, "/path2", testHandlerFunc)
	_ = router.Register(http.MethodGet, "/1/classes/{className}/{objectId}", testHandlerFunc)
	_ = router.Register(http.MethodGet, "/1/classes/{className}", testHandlerFunc)
	_ = router.Register(http.MethodGet, "/1/classes/{className}", testHandlerFunc)
	_ = router.Register(http.MethodPost, "/1/classes/{className}", testHandlerFunc)
	_ = router.Register(http.MethodGet, "/activities/{activityId}/people/{collection}", testHandlerFunc)
	_ = router.Register(http.MethodGet, "/activities/{activityId}/comments", testHandlerFunc)
	_ = router.Register(http.MethodGet, "/users/{user}/starred", testHandlerFunc)
	_ = router.Register(http.MethodGet, "/user/starred", testHandlerFunc)
	_ = router.Register(http.MethodGet, "/users", testHandlerFunc)
	_ = router.Register(http.MethodGet, "/path1/{id}", testHandlerFunc)
	_ = router.Register(http.MethodPost, "/path1", testHandlerFunc)
	_ = router.Register(http.MethodPut, "/path1/{id}", testHandlerFunc)
	_ = router.Register(http.MethodDelete, "/path1/{id}", testHandlerFunc)
	_ = router.Register(http.MethodGet, "/path1/{id:[0-9]+}/{name:[a-z]+}", testHandlerFunc)

	assertPathFound(t, router, "GET", "/path1")
	assertPathFound(t, router, "GET", "/path2")
	assertPathFound(t, router, "GET", "/1/classes/{className}/{objectId}")
	assertPathFound(t, router, "GET", "/1/classes/{className}")
	assertPathFound(t, router, "GET", "/1/classes/{className}")
	assertPathFound(t, router, "POST", "/1/classes/{className}")
	assertPathFound(t, router, "GET", "/activities/{activityId}/people/{collection}")
	assertPathFound(t, router, "GET", "/activities/{activityId}/comments")
	assertPathFound(t, router, "GET", "/users/{user}/starred")
	assertPathFound(t, router, "GET", "/user/starred")
	assertPathFound(t, router, "GET", "/users")
	assertPathFound(t, router, "GET", "/path1/{id}")
	assertPathFound(t, router, "POST", "/path1")
	assertPathFound(t, router, "PUT", "/path1/{id}")
	assertPathFound(t, router, "DELETE", "/path1/{id}")
	assertPathFound(t, router, "GET", "/path1/100/abc")
	assertPathNotFound(t, router, "GET", "/path1/100/123")
}

func TestGetURLParamatersBagInHandler(t *testing.T) {
	router := Router{}

	bag := newURLParameterBag(2, false)
	bag.add("id", "100")
	bag.add("name", "dummy")

	f := assertRequestHasParameterHandler(t, bag)
	_ = router.Register(http.MethodGet, "/path1/{id:[0-9]+}/{name:[a-z]{1,5}}", f)

	assertPathFound(t, router, "GET", "/path1/100/dummy")
}

func TestVerbsMethodsAreWorking(t *testing.T) {
	path := "/path1"

	router := Router{}
	_ = router.Get(path, testHandlerFunc)
	_ = router.Head(path, testHandlerFunc)
	_ = router.Post(path, testHandlerFunc)
	_ = router.Put(path, testHandlerFunc)
	_ = router.Patch(path, testHandlerFunc)
	_ = router.Delete(path, testHandlerFunc)
	_ = router.Connect(path, testHandlerFunc)
	_ = router.Options(path, testHandlerFunc)
	_ = router.Trace(path, testHandlerFunc)

	assertPathFound(t, router, "GET", path)
	assertPathFound(t, router, "HEAD", path)
	assertPathFound(t, router, "POST", path)
	assertPathFound(t, router, "PUT", path)
	assertPathFound(t, router, "PATCH", path)
	assertPathFound(t, router, "DELETE", path)
	assertPathFound(t, router, "CONNECT", path)
	assertPathFound(t, router, "OPTIONS", path)
	assertPathFound(t, router, "TRACE", path)
}
