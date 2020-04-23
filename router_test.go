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

func TestGetURLParamatersBagInHandler(t *testing.T) {
	mainRouter := Router{}
	postsRouter := Router{}

	bag := newURLParameterBag(2, false)
	bag.add("id", "100")
	bag.add("name", "dummy")
	f := assertRequestHasParameterHandler(t, bag)
	_ = mainRouter.Register(http.MethodGet, "/path1/{id}/{name:[a-z]{1,5}}", f)

	bag2 := newURLParameterBag(2, false)
	bag2.add("name", "dummy/file/src/image.jpg")
	f2 := assertRequestHasParameterHandler(t, bag2)
	_ = mainRouter.Register(http.MethodGet, "/path1/{name:.*}", f2)

	bag3 := newURLParameterBag(2, false)
	bag3.add("name", "2020-05-05")
	f3 := assertRequestHasParameterHandler(t, bag3)
	_ = mainRouter.Register(http.MethodGet, "/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}", f3)

	bag4 := newURLParameterBag(2, false)
	bag4.add("id", "123")
	bag4.add("name", "2020-05-05")
	f4 := assertRequestHasParameterHandler(t, bag4)
	_ = postsRouter.Register(http.MethodGet, "/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}", f4)

	mainRouter.Prefix("/posts/{id}", &postsRouter)

	assertPathFound(t, mainRouter, "GET", "/path1/100/dummy")
	assertPathFound(t, mainRouter, "GET", "/path1/dummy/file/src/image.jpg")
	assertPathFound(t, mainRouter, "GET", "/2020-05-05")
	assertPathFound(t, mainRouter, "GET", "/posts/123/2020-05-05")
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

func TestGetURLParametersReturnsEmptyBagIfNoContextValueExists(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/dummy", nil)

	bag := GetURLParameters(r)

	assertBagSetting(t, bag, 0, true)
}

func TestTreeRouterFindsPathsWhenPrefixingRouters(t *testing.T) {
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