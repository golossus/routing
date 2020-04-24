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

func ExampleRegisterHandler_UsingHandlerAlias() {
	h, _ := routing.GetHandler("myhandler")

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h(w, r)

	fmt.Println(w.Body.String())
	// Output: Hello MyHandler
}

func ExampleRegisterHandler_UsingHandlerCanonicalName() {
	h, _ := routing.GetHandler("github.com/golossus/routing_test.MyHandler")

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h(w, r)

	fmt.Println(w.Body.String())
	// Output: Hello MyHandler
}
