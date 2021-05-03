package routing

import "net/http"

type matcher func(r *http.Request) bool

func byHost(host string) matcher {
	return func(r *http.Request) bool {
		return r.Host == host
	}
}
