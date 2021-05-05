package routing

import (
	"net/http"
	"strings"
)

func getAutoMethodOptionsHandler(router *Router) http.HandlerFunc {

	return func(writer http.ResponseWriter, request *http.Request) {
		availVerbs := getAvailableMethods(router, request)

		writer.WriteHeader(http.StatusNoContent)
		writer.Header().Set("Allow", strings.Join(availVerbs, ", "))
	}
}

func getRedirectHandler(url string, code ...int) http.HandlerFunc {
	defaultCode := http.StatusFound
	if len(code) > 0 && code[0] >= http.StatusMultipleChoices && code[0] <= http.StatusPermanentRedirect {
		defaultCode = code[0]
	}
	return func(writer http.ResponseWriter, request *http.Request, ) {
		http.Redirect(writer, request, url, defaultCode)
	}
}

func getAvailableMethods(router *Router, request *http.Request) []string {
	availVerbs := make([]string, 0, 9)
	for verb, tree := range router.trees {
		n := tree.find(request)
		if n != nil {
			availVerbs = append(availVerbs, verb)
		}
	}
	return availVerbs
}
