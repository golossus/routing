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
