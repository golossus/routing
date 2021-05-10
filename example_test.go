package routing_test

import (
	"fmt"
	"github.com/golossus/routing"
	"net/http"
	"net/http/httptest"
)

func init() {
	routing.AddHandler(MyHandler, "myhandler")
}

func MyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello MyHandler")
}

func ExampleAddHandler_usingHandlerAlias() {
	h, _ := routing.GetHandler("myhandler")

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h(w, r)

	fmt.Println(w.Body.String())
	// Output: Hello MyHandler
}

func ExampleAddHandler_usingHandlerCanonicalName() {
	h, _ := routing.GetHandler("github.com/golossus/routing_test.MyHandler")

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h(w, r)

	fmt.Println(w.Body.String())
	// Output: Hello MyHandler
}

func ExampleRouter_NewRoute_routeRegistration() {
	h, _ := routing.GetHandler("github.com/golossus/routing_test.MyHandler")

	router := routing.NewRouter()
	_ = router.NewRoute().Path("/").Method("GET").Handler(h).Register()

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	fmt.Println(w.Body.String())
	// Output: Hello MyHandler
}
