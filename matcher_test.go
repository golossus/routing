package routing

import (
	"net/http"
	"reflect"
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
	assertTrue(t, matches)
	assertFalse(t, leaf.hasParameters())

	matches2, leaf2 := m2(req)
	assertFalse(t, matches2)
	assertFalse(t, leaf2.hasParameters())
}

func Test_byHost_WithDynamicHosts(t *testing.T) {
	h := "app.{subdomain:[a-z]+}.test2.com"

	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = "app.golossus.test2.com"

	m, _ := byHost(h)
	matches, leaf := m(req)
	assertTrue(t, matches)
	assertTrue(t, leaf.hasParameters())

	req.Host = "app.1234.test2.com"
	m, _ = byHost(h)
	matches, leaf = m(req)
	assertFalse(t, matches)
	assertTrue(t, leaf.hasParameters())
}

func Test_byHost_ReturnsErrorWhenMalformedHost(t *testing.T) {
	h := "app.{subdomain:[a-z]+}{m}.test2.com"

	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = "app.golossus.test2.com"

	m, err := byHost(h)
	assertNil(t, m)
	assertNotNil(t, err)
}

func assertTrue(t *testing.T, value bool) {
	if !value{
		t.Errorf("%v is not true", value)
	}
}

func assertFalse(t *testing.T, value bool) {
	if value{
		t.Errorf("%v is not false", value)
	}
}

func assertNil(t *testing.T, value interface{}) {
	reflectedValue := reflect.ValueOf(value)
	if !reflectedValue.IsNil(){
		t.Errorf("%v is not nil", value)
	}
}

func assertNotNil(t *testing.T, value interface{}) {
	if value == nil{
		t.Errorf("%v is nil", value)
	}
}