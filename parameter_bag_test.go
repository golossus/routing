package http_router

import "testing"

func TestUrlParameterBagEmptyOnCreation(t *testing.T) {
	bag := NewUrlParameterBag()

	if bag.params != nil {
		t.Errorf("")
	}
}

func TestUrlParameterBagAddsParameter(t *testing.T) {
	bag := NewUrlParameterBag()

	bag.addParameter(urlParameter{name: "param1", value: "v1"})

	if bag.params == nil {
		t.Errorf("")
	}

	if len(bag.params) == 0 {
		t.Errorf("")
	}

	if bag.params[0].name != "param1" {
		t.Errorf("")
	}
}

func TestUrlParameterBagGetByName(t *testing.T) {
	bag := NewUrlParameterBag()

	bag.addParameter(urlParameter{name: "param1", value: "v1"})
	bag.addParameter(urlParameter{name: "param2", value: "v2"})
	bag.addParameter(urlParameter{name: "param3", value: "v3"})

	if bag.params == nil {
		t.Errorf("")
	}

	if len(bag.params) != 3 {
		t.Errorf("")
	}

	if bag.GetByName("param1", "") != "v1" {
		t.Errorf("")
	}

	if bag.GetByName("param2", "") != "v2" {
		t.Errorf("")
	}

	if bag.GetByName("param3", "") != "v3" {
		t.Errorf("")
	}

	if bag.GetByName("param4", "v4") != "v4" {
		t.Errorf("")
	}
}

func TestUrlParameterBagGetByIndex(t *testing.T) {
	bag := NewUrlParameterBag()

	bag.addParameter(urlParameter{name: "param1", value: "v1"})
	bag.addParameter(urlParameter{name: "param2", value: "v2"})
	bag.addParameter(urlParameter{name: "param3", value: "v3"})

	if bag.params == nil {
		t.Errorf("")
	}

	if len(bag.params) != 3 {
		t.Errorf("")
	}

	if bag.GetByIndex(0, "") != "v1" {
		t.Errorf("")
	}

	if bag.GetByIndex(1, "") != "v2" {
		t.Errorf("")
	}

	if bag.GetByIndex(2, "") != "v3" {
		t.Errorf("")
	}

	if bag.GetByIndex(3, "v4") != "v4" {
		t.Errorf("")
	}
}
