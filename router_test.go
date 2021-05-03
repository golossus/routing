package routing

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testHandlerFunc = func(response http.ResponseWriter, request *http.Request) {
	_, _ = fmt.Fprint(response, request.URL.Path)
}

var testDummyHandlerFunc = func(response http.ResponseWriter, request *http.Request) {
	_, _ = fmt.Fprint(response, "dummy")
}

func assertPathFound(t *testing.T, router Router, method, path string) {
	r, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 200 || w.Body.String() != path {
		t.Errorf("%s %s not found", method, path)
	}
}

func assertPathWithHostFound(t *testing.T, router Router, method, path, host string) {
	r, _ := http.NewRequest(method, path, nil)
	r.Host = host

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
			t.Errorf("size of parameter bag doesn't match %d != %d", len(bag.params), len(par.params))
		}

		for i := 0; i < len(bag.params); i++ {
			bagValue, _ := bag.GetByIndex(uint(i))
			parValue, _ := par.GetByIndex(uint(i))
			if bagValue != parValue {
				t.Errorf("parameter at index %d don't match", i)
			}
		}

		testHandlerFunc(w, r)
	}
}

func assertRouteNameHasHandler(t *testing.T, mainRouter Router, method, path, routeName string) {
	leaf, ok := mainRouter.routes[routeName]
	if !ok {
		t.Errorf("route name %s not found", routeName)
	}

	r, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	leaf.handler(w, r)

	if w.Result().StatusCode != 200 || w.Body.String() != path {
		t.Errorf("%s %s not found as %s route name", method, path, routeName)
	}
}

func TestRouter_ServeHTTP_FindsPaths(t *testing.T) {
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
	_ = router.Register(http.MethodGet, "/path1/{id}/path2", testHandlerFunc)
	_ = router.Register(http.MethodGet, "/path1/{id}-path2", testHandlerFunc)
	_ = router.Register(http.MethodGet, "/{date}/", testHandlerFunc)
	_ = router.Register(http.MethodGet, "/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}/", testHandlerFunc)
	_ = router.Register(http.MethodGet, "/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}", testHandlerFunc)
	_ = router.Register(http.MethodGet, "/path3/{slug:[0-9]+}", testHandlerFunc)
	_ = router.Register(http.MethodGet, "/path3/{slug:.*}", testHandlerFunc)
	_ = router.Register(http.MethodGet, "/path4/{id:[0-9]+}", testHandlerFunc)
	_ = router.Register(http.MethodGet, "/path4/{id:[0-9]+}/{slug:[a-z]+}", testHandlerFunc)

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
	assertPathFound(t, router, "GET", "/path1/100/path2")
	assertPathFound(t, router, "GET", "/path1/100-path2")
	assertPathFound(t, router, "GET", "/october/")
	assertPathFound(t, router, "GET", "/2019-02-20")
	assertPathFound(t, router, "GET", "/2019-02-20/")
	assertPathFound(t, router, "GET", "/path3/00545")
	assertPathFound(t, router, "GET", "/path3/00545/5456/file/file.jpg")
	assertPathFound(t, router, "GET", "/path4/00545")
	assertPathFound(t, router, "GET", "/path4/00545/abc")

	assertPathNotFound(t, router, "GET", "/path1/100/123")
}

func TestGetURLParameters(t *testing.T) {
	mainRouter := Router{}
	postsRouter := Router{}

	bag := newURLParameterBag(2)
	bag.add("id", "100")
	bag.add("name", "dummy")
	f := assertRequestHasParameterHandler(t, bag)
	_ = mainRouter.Register(http.MethodGet, "/path1/{id}/{name:[a-z]{1,5}}", f)

	bag2 := newURLParameterBag(2)
	bag2.add("name", "dummy/file/src/image.jpg")
	f2 := assertRequestHasParameterHandler(t, bag2)
	_ = mainRouter.Register(http.MethodGet, "/path1/{name:.*}", f2)

	bag3 := newURLParameterBag(2)
	bag3.add("name", "2020-05-05")
	f3 := assertRequestHasParameterHandler(t, bag3)
	_ = mainRouter.Register(http.MethodGet, "/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}", f3)

	bag4 := newURLParameterBag(2)
	bag4.add("id", "123")
	bag4.add("name", "2020-05-05")
	f4 := assertRequestHasParameterHandler(t, bag4)
	_ = postsRouter.Register(http.MethodGet, "/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}", f4)

	_ = mainRouter.Prefix("/posts/{id}", &postsRouter)

	assertPathFound(t, mainRouter, "GET", "/path1/100/dummy")
	assertPathFound(t, mainRouter, "GET", "/path1/dummy/file/src/image.jpg")
	assertPathFound(t, mainRouter, "GET", "/2020-05-05")
	assertPathFound(t, mainRouter, "GET", "/posts/123/2020-05-05")
}

func TestGetURLParameters_ContainsHostParameters(t *testing.T) {
	mainRouter := Router{}

	bag := newURLParameterBag(2)
	bag.add("id", "100")
	bag.add("subdomain", "dummy")
	bag.add("domain", "test")

	f := assertRequestHasParameterHandler(t, bag)
	options := NewMatchingOptions()
	options.Host = "{subdomain:[a-z]+}.{domain:[a-z]+}.com"

	_ = mainRouter.Register(http.MethodGet, "/path1/{id}", f, options)

	assertPathWithHostFound(t, mainRouter, "GET", "/path1/100", "dummy.test.com")
}

func TestRouter_AllVerbs(t *testing.T) {
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

func TestGetURLParameters_ReturnsEmptyBagIfNoContextValueExists(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/dummy", nil)

	bag := GetURLParameters(r)

	assertBagSetting(t, bag, 0)
}

func TestRouter_ServeHTTP_FindsPathsWhenPrefixingRouters(t *testing.T) {
	mainRouter := Router{}
	_ = mainRouter.Register(http.MethodGet, "/path1", testHandlerFunc)
	_ = mainRouter.Register(http.MethodGet, "/path2", testHandlerFunc)
	_ = mainRouter.Register(http.MethodGet, "/1/classes/{className}/{objectId}", testHandlerFunc)
	_ = mainRouter.Register(http.MethodGet, "/1/classes/{className}", testHandlerFunc)
	_ = mainRouter.Register(http.MethodGet, "/1/classes/{className}", testHandlerFunc)
	_ = mainRouter.Register(http.MethodPost, "/1/classes/{className}", testHandlerFunc)
	_ = mainRouter.Register(http.MethodGet, "/activities/{activityId}/people/{collection}", testHandlerFunc)
	_ = mainRouter.Register(http.MethodGet, "/activities/{activityId}/comments", testHandlerFunc)
	_ = mainRouter.Register(http.MethodGet, "/users/{user}/starred", testHandlerFunc)
	_ = mainRouter.Register(http.MethodGet, "/user/starred", testHandlerFunc)
	_ = mainRouter.Register(http.MethodGet, "/users", testHandlerFunc)
	_ = mainRouter.Register(http.MethodPost, "/path1", testHandlerFunc)
	_ = mainRouter.Register(http.MethodGet, "/{date}/", testHandlerFunc)
	_ = mainRouter.Register(http.MethodGet, "/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}/", testHandlerFunc)
	_ = mainRouter.Register(http.MethodGet, "/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}", testHandlerFunc)
	_ = mainRouter.Register(http.MethodGet, "/path3/{slug:[0-9]+}", testHandlerFunc)
	_ = mainRouter.Register(http.MethodGet, "/path3/{slug:.*}", testHandlerFunc)
	_ = mainRouter.Register(http.MethodGet, "/path4/{id:[0-9]+}", testHandlerFunc)
	_ = mainRouter.Register(http.MethodGet, "/path4/{id:[0-9]+}/{slug:[a-z]+}", testHandlerFunc)

	path1Router := Router{}
	_ = path1Router.Register(http.MethodGet, "/{id}", testHandlerFunc)
	_ = path1Router.Register(http.MethodPut, "/{id}", testHandlerFunc)
	_ = path1Router.Register(http.MethodDelete, "/{id}", testHandlerFunc)
	_ = path1Router.Register(http.MethodGet, "/{id:[0-9]+}/{name:[a-z]+}", testHandlerFunc)
	_ = path1Router.Register(http.MethodGet, "/{id}/path2", testHandlerFunc)
	_ = path1Router.Register(http.MethodGet, "/{id}-path2", testHandlerFunc)

	userRouter := Router{}
	_ = userRouter.Register(http.MethodGet, "/profile", testHandlerFunc)
	_ = userRouter.Register(http.MethodGet, "/{date}/posts", testHandlerFunc)

	_ = mainRouter.Prefix("/path1", &path1Router)
	_ = mainRouter.Prefix("/user/{id}", &userRouter)

	assertPathFound(t, mainRouter, "GET", "/path1")
	assertPathFound(t, mainRouter, "GET", "/path2")
	assertPathFound(t, mainRouter, "GET", "/1/classes/{className}/{objectId}")
	assertPathFound(t, mainRouter, "GET", "/1/classes/{className}")
	assertPathFound(t, mainRouter, "GET", "/1/classes/{className}")
	assertPathFound(t, mainRouter, "POST", "/1/classes/{className}")
	assertPathFound(t, mainRouter, "GET", "/activities/{activityId}/people/{collection}")
	assertPathFound(t, mainRouter, "GET", "/activities/{activityId}/comments")
	assertPathFound(t, mainRouter, "GET", "/users/{user}/starred")
	assertPathFound(t, mainRouter, "GET", "/user/starred")
	assertPathFound(t, mainRouter, "GET", "/users")
	assertPathFound(t, mainRouter, "GET", "/path1/{id}")
	assertPathFound(t, mainRouter, "POST", "/path1")
	assertPathFound(t, mainRouter, "PUT", "/path1/{id}")
	assertPathFound(t, mainRouter, "DELETE", "/path1/{id}")
	assertPathFound(t, mainRouter, "GET", "/path1/100/abc")
	assertPathFound(t, mainRouter, "GET", "/path1/100/path2")
	assertPathFound(t, mainRouter, "GET", "/path1/100-path2")
	assertPathFound(t, mainRouter, "GET", "/october/")
	assertPathFound(t, mainRouter, "GET", "/2019-02-20")
	assertPathFound(t, mainRouter, "GET", "/2019-02-20/")
	assertPathFound(t, mainRouter, "GET", "/path3/00545")
	assertPathFound(t, mainRouter, "GET", "/path3/00545/5456/file/file.jpg")
	assertPathFound(t, mainRouter, "GET", "/path4/00545")
	assertPathFound(t, mainRouter, "GET", "/path4/00545/abc")
	assertPathFound(t, mainRouter, "GET", "/user/5/profile")
	assertPathFound(t, mainRouter, "GET", "/user/5/2020-03-01/posts")

	assertPathNotFound(t, mainRouter, "GET", "/path1/100/123")
}

func TestRouter_As_AssignsRouteNames(t *testing.T) {
	mainRouter := Router{}

	_ = mainRouter.As("users.get").Get("/users", testHandlerFunc)
	_ = mainRouter.As("users.create").Post("/users", testHandlerFunc)
	_ = mainRouter.As("users.create").Post("/users/create", testHandlerFunc)
	_ = mainRouter.As("users.update").Put("/users/{id}", testHandlerFunc)
	_ = mainRouter.As("users.delete").Delete("/users/{id}", testDummyHandlerFunc)
	_ = mainRouter.As("users.softDelete").Delete("/users/{id}", testHandlerFunc)
	_ = mainRouter.Get("/users/profile", testDummyHandlerFunc)

	apiRouter := Router{}
	_ = apiRouter.As("users.account").Get("/users/account", testHandlerFunc)
	_ = apiRouter.As("users.profile").Get("/users/profile", testHandlerFunc)

	_ = mainRouter.As("api.").Prefix("/api", &apiRouter)

	assertRouteNameHasHandler(t, mainRouter, http.MethodGet, "/users", "users.get")
	assertRouteNameHasHandler(t, mainRouter, http.MethodPost, "/users/create", "users.create")
	assertRouteNameHasHandler(t, mainRouter, http.MethodPut, "/users/100", "users.update")
	assertRouteNameHasHandler(t, mainRouter, http.MethodDelete, "/users/100", "users.delete")
	assertRouteNameHasHandler(t, mainRouter, http.MethodDelete, "/users/100", "users.softDelete")

	assertRouteNameHasHandler(t, mainRouter, http.MethodGet, "/api/users/account", "api.users.account")
	assertRouteNameHasHandler(t, mainRouter, http.MethodGet, "/api/users/profile", "api.users.profile")
}

func TestRouter_MatchingOptions_AssignsRouteNames(t *testing.T) {
	mainRouter := Router{}

	_ = mainRouter.Get("/users", testHandlerFunc, MatchingOptions{"users.get", "", []string{}})
	_ = mainRouter.Post("/users", testHandlerFunc, MatchingOptions{"users.create", "", []string{}})
	_ = mainRouter.Post("/users/create", testHandlerFunc, MatchingOptions{"users.create", "", []string{}})
	_ = mainRouter.Put("/users/{id}", testHandlerFunc, MatchingOptions{"users.update", "", []string{}})
	_ = mainRouter.Delete("/users/{id}", testDummyHandlerFunc, MatchingOptions{"users.delete", "", []string{}})
	_ = mainRouter.Delete("/users/{id}", testHandlerFunc, MatchingOptions{"users.softDelete", "", []string{}})
	_ = mainRouter.Get("/users/profile", testDummyHandlerFunc)

	apiRouter := Router{}
	_ = apiRouter.Get("/users/account", testHandlerFunc, MatchingOptions{"users.account", "", []string{}})
	_ = apiRouter.Get("/users/profile", testHandlerFunc, MatchingOptions{"users.profile", "", []string{}})

	_ = mainRouter.Prefix("/api", &apiRouter, "api.")

	assertRouteNameHasHandler(t, mainRouter, http.MethodGet, "/users", "users.get")
	assertRouteNameHasHandler(t, mainRouter, http.MethodPost, "/users/create", "users.create")
	assertRouteNameHasHandler(t, mainRouter, http.MethodPut, "/users/100", "users.update")
	assertRouteNameHasHandler(t, mainRouter, http.MethodDelete, "/users/100", "users.delete")
	assertRouteNameHasHandler(t, mainRouter, http.MethodDelete, "/users/100", "users.softDelete")

	assertRouteNameHasHandler(t, mainRouter, http.MethodGet, "/api/users/account", "api.users.account")
	assertRouteNameHasHandler(t, mainRouter, http.MethodGet, "/api/users/profile", "api.users.profile")
}

func TestRouter_MatchingOptions_AssignsRouteNamesOverAsMethod(t *testing.T) {
	mainRouter := Router{}

	_ = mainRouter.As("users.getAs").Get("/users", testHandlerFunc, MatchingOptions{"users.get", "", []string{}})
	_ = mainRouter.As("users.createAs").Post("/users", testHandlerFunc, MatchingOptions{"users.create", "", []string{}})
	_ = mainRouter.As("users.createAs").Post("/users/create", testHandlerFunc, MatchingOptions{"users.create", "", []string{}})
	_ = mainRouter.As("users.updateAs").Put("/users/{id}", testHandlerFunc, MatchingOptions{"users.update", "", []string{}})
	_ = mainRouter.As("users.deleteAs").Delete("/users/{id}", testDummyHandlerFunc, MatchingOptions{"users.delete", "", []string{}})
	_ = mainRouter.As("users.softDeleteAs").Delete("/users/{id}", testHandlerFunc, MatchingOptions{"users.softDelete", "", []string{}})
	_ = mainRouter.Get("/users/profile", testDummyHandlerFunc)

	apiRouter := Router{}
	_ = apiRouter.Get("/users/account", testHandlerFunc, MatchingOptions{"users.account", "", []string{}})
	_ = apiRouter.Get("/users/profile", testHandlerFunc, MatchingOptions{"users.profile", "", []string{}})

	_ = mainRouter.Prefix("/api", &apiRouter, "api.")

	assertRouteNameHasHandler(t, mainRouter, http.MethodGet, "/users", "users.get")
	assertRouteNameHasHandler(t, mainRouter, http.MethodPost, "/users/create", "users.create")
	assertRouteNameHasHandler(t, mainRouter, http.MethodPut, "/users/100", "users.update")
	assertRouteNameHasHandler(t, mainRouter, http.MethodDelete, "/users/100", "users.delete")
	assertRouteNameHasHandler(t, mainRouter, http.MethodDelete, "/users/100", "users.softDelete")

	assertRouteNameHasHandler(t, mainRouter, http.MethodGet, "/api/users/account", "api.users.account")
	assertRouteNameHasHandler(t, mainRouter, http.MethodGet, "/api/users/profile", "api.users.profile")
}

func TestRouter_MatchingOptions_MatchesByHost(t *testing.T) {
	mainRouter := Router{}

	_ = mainRouter.Get("/users", testHandlerFunc, NewMatchingOptions())
	_ = mainRouter.Get("/users/{id}", testHandlerFunc, NewMatchingOptions())
	_ = mainRouter.Get("/users/{id}/create", testHandlerFunc, MatchingOptions{"", "test.com", []string{}})

	apiRouter := Router{}
	_ = apiRouter.Get("/users/account", testHandlerFunc, MatchingOptions{"", "api.test.com", []string{}})
	_ = mainRouter.Prefix("/api", &apiRouter, "api.")

	req, _ := http.NewRequest("GET", "/users/1/create", nil)
	req.Host = "test.com"
	res := httptest.NewRecorder()
	mainRouter.ServeHTTP(res, req)

	assertEqual(t, 200, res.Code)

	req, _ = http.NewRequest("GET", "/users/1/create", nil)
	res = httptest.NewRecorder()
	mainRouter.ServeHTTP(res, req)

	assertEqual(t, 404, res.Code)

	req, _ = http.NewRequest("GET", "/api/users/account", nil)
	res = httptest.NewRecorder()
	mainRouter.ServeHTTP(res, req)

	assertEqual(t, 404, res.Code)

	req, _ = http.NewRequest("GET", "/api/users/account", nil)
	req.Host = "api.test.com"
	res = httptest.NewRecorder()
	mainRouter.ServeHTTP(res, req)

	assertEqual(t, 200, res.Code)
}

func TestRouter_MatchingOptions_MatchesBySchemas(t *testing.T) {
	mainRouter := Router{}

	_ = mainRouter.Get("/users", testHandlerFunc, NewMatchingOptions())
	_ = mainRouter.Get("/users/{id}", testHandlerFunc, MatchingOptions{"", "", []string{"Http", "ftp"}})
	_ = mainRouter.Get("/users/{id}/create", testHandlerFunc, MatchingOptions{"", "", []string{"https"}})

	req, _ := http.NewRequest("GET", "/users/1/create", nil)
	req.URL.Scheme = "https"
	res := httptest.NewRecorder()
	mainRouter.ServeHTTP(res, req)

	assertEqual(t, 200, res.Code)

	req, _ = http.NewRequest("GET", "/users/1/create", nil)
	req.URL.Scheme = "http"
	res = httptest.NewRecorder()
	mainRouter.ServeHTTP(res, req)

	assertEqual(t, 404, res.Code)

	req, _ = http.NewRequest("GET", "/users/1", nil)
	req.URL.Scheme = "https"
	res = httptest.NewRecorder()
	mainRouter.ServeHTTP(res, req)

	assertEqual(t, 404, res.Code)

	req, _ = http.NewRequest("GET", "/users/1", nil)
	req.URL.Scheme = "ftp"
	res = httptest.NewRecorder()
	mainRouter.ServeHTTP(res, req)

	assertEqual(t, 200, res.Code)

	req, _ = http.NewRequest("GET", "/users/1", nil)
	req.URL.Scheme = "http"
	res = httptest.NewRecorder()
	mainRouter.ServeHTTP(res, req)

	assertEqual(t, 200, res.Code)
}

func TestRouter_MatchingOptions_MatchesByHostReturnsErrorWhenMalformedHost(t *testing.T) {
	mainRouter := Router{}

	err := mainRouter.Get("/users", testHandlerFunc, MatchingOptions{"", "app.{subdomain:[a-z]+}{m}.test2.com", []string{}})
	assertNotNil(t, err)
}

func TestRouter_GenerateURL_GenerateValidRoutes(t *testing.T) {
	mainRouter := Router{}

	_ = mainRouter.As("path1").Register(http.MethodGet, "/path1", testHandlerFunc)
	_ = mainRouter.As("path1.id.name").Register(http.MethodGet, "/path1/{id}/{name:[a-z]{1,5}}", testHandlerFunc)
	_ = mainRouter.As("path1.file").Register(http.MethodGet, "/path1/{file:.*}", testHandlerFunc)
	_ = mainRouter.As("date").Register(http.MethodGet, "/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}", testHandlerFunc)

	postsRouter := Router{}
	_ = postsRouter.As("date").Register(http.MethodGet, "/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}", testHandlerFunc)

	_ = mainRouter.As("posts.id.").Prefix("/posts/{id}", &postsRouter)

	assertRouteIsGenerated(t, mainRouter, "path1", "/path1", map[string]string{})
	assertRouteIsGenerated(t, mainRouter, "path1.id.name", "/path1/100/abc", map[string]string{"id": "100", "name": "abc"})
	assertRouteIsGenerated(t, mainRouter, "path1.file", "/path1/100/2098939/image.jpg", map[string]string{"file": "100/2098939/image.jpg"})
	assertRouteIsGenerated(t, mainRouter, "date", "/2020-05-05", map[string]string{"date": "2020-05-05"})
	assertRouteIsGenerated(t, mainRouter, "posts.id.date", "/posts/10/2020-05-05", map[string]string{"id": "10", "date": "2020-05-05"})
}

func assertRouteIsGenerated(t *testing.T, mainRouter Router, name, url string, params map[string]string) {
	bag := URLParameterBag{}
	for key, value := range params {
		bag.add(key, value)
	}
	route2, err2 := mainRouter.GenerateURL(name, bag)
	if err2 != nil {
		t.Errorf("route %s not generated", name)
	}
	if route2 != url {
		t.Errorf("route %s not valid", name)
	}
}

type sliceLoader []RouteDef

func (l *sliceLoader) Load() []RouteDef {
	var x []RouteDef = *l
	return x
}

func TestRouter_Load_RegisterRoutes(t *testing.T) {
	AddHandler(testHandlerFunc, "users.Handler")

	router := NewRouter()
	loader := sliceLoader{
		RouteDef{Name: "get.users", Method: "GET", Schema: "/users", Handler: "users.Handler"},
	}
	err := router.Load(&loader)
	if err != nil {
		t.Error(err)
	}
	assertRouteIsGenerated(t, router, "get.users", "/users", nil)
	assertPathFound(t, router, "GET", "/users")
}

func TestRouter_Load_FailsWhenHandlerNoExists(t *testing.T) {
	AddHandler(testHandlerFunc, "users.Handler")

	router := NewRouter()
	loader := sliceLoader{
		RouteDef{Name: "get.users", Method: "GET", Schema: "/users", Handler: "users.Handler.no.exists"},
	}
	err := router.Load(&loader)
	if err == nil {
		t.Errorf("no registered Handler %s has been loaded", "users.Handler.no.exists")
	}
}

func TestRouter_Load_FailsWhenSchemaIsInvalid(t *testing.T) {
	AddHandler(testHandlerFunc, "users.Handler")

	router := NewRouter()
	loader := sliceLoader{
		RouteDef{Name: "get.users", Method: "GET", Schema: "users", Handler: "users.Handler.no.exists"},
	}
	err := router.Load(&loader)
	if err == nil {
		t.Errorf("invalid Schema %s has been loaded", "users")
	}
}

func assertEqual(t *testing.T, expected, value int) {
	if expected != value {
		t.Errorf("%v is not equal to %v", expected, value)
	}
}
