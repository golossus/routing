package routing

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMiddlewarePipe_ReturnsEmptyPipe(t *testing.T) {

	pipe := NewMiddlewarePipe()
	assertEqual(t, 0, len(pipe.middlewares))
}

func TestMiddlewarePipe_Next(t *testing.T) {
	pipe := NewMiddlewarePipe()
	pipe.Next(
		func(handlerFunc http.HandlerFunc) http.HandlerFunc {
			return handlerFunc
		},
		func(handlerFunc http.HandlerFunc) http.HandlerFunc {
			return handlerFunc
		},
	)

	assertEqual(t, 2, len(pipe.middlewares))
}

func TestMiddlewarePipe_Pipe(t *testing.T) {
	pipe := NewMiddlewarePipe()
	pipe.Next(
		func(handlerFunc http.HandlerFunc) http.HandlerFunc {
			return handlerFunc
		},
		func(handlerFunc http.HandlerFunc) http.HandlerFunc {
			return handlerFunc
		},
	)

	pipe2 := NewMiddlewarePipe()
	pipe2.Next(
		func(handlerFunc http.HandlerFunc) http.HandlerFunc {
			return handlerFunc
		},
	)

	pipe.Pipe(pipe2)

	assertEqual(t, 3, len(pipe.middlewares))
}

func TestMiddlewarePipe_Then(t *testing.T) {
	pipe := NewMiddlewarePipe()
	handler := pipe.Next(func(handlerFunc http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Access-Control-Allow-Origin", "https://test.com")
			handlerFunc(writer, request)
		}
	}).Next(func(handlerFunc http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			handlerFunc(writer, request)
		}
	}).Then(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusAccepted)
	})

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	handler(response, request)

	assertEqual(t, http.StatusAccepted, response.Code)
	assertStringEqual(t, "application/json", response.Header().Get("Content-Type"))
	assertStringEqual(t, "https://test.com", response.Header().Get("Access-Control-Allow-Origin"))

}
