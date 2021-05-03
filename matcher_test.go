package routing

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_byHost_WithSpecificHosts(t *testing.T) {
	h := "test.com"
	h2 := "test2.com"

	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = h

	m := byHost(h)
	m2 := byHost(h2)

	assert.True(t, m(req))
	assert.False(t, m2(req))
}
