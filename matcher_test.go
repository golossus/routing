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

	m, _ := byHost(h)
	m2, _ := byHost(h2)

	matches, leaf := m(req)
	assert.True(t, matches)
	assert.False(t, leaf.hasParameters())

	matches2, leaf2 := m2(req)
	assert.False(t, matches2)
	assert.False(t, leaf2.hasParameters())
}

func Test_byHost_WithDynamicHosts(t *testing.T) {
	h := "app.{subdomain:[a-z]+}.test2.com"

	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = "app.golossus.test2.com"

	m, _ := byHost(h)
	matches, leaf := m(req)
	assert.True(t, matches)
	assert.True(t, leaf.hasParameters())

	req.Host = "app.1234.test2.com"
	m, _ = byHost(h)
	matches, leaf = m(req)
	assert.False(t, matches)
	assert.True(t, leaf.hasParameters())
}

func Test_byHost_ReturnsErrorWhenMalformedHost(t *testing.T) {
	h := "app.{subdomain:[a-z]+}{m}.test2.com"

	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = "app.golossus.test2.com"

	m, err := byHost(h)
	assert.Nil(t, m)
	assert.Error(t, err)
}