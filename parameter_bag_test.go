package http_router

import "testing"

func TestUrlParameterBagEmptyValuesOnCreation(t *testing.T) {
	bag := UrlParameterBag{}

	if bag.params != nil {
		t.Errorf("")
	}

	if bag.capacity != 0 {
		t.Errorf("")
	}

	if bag.reverse != false {
		t.Errorf("")
	}
}

func TestNewUrlParameterBagSetsRightValuesOnCreation(t *testing.T) {
	bag := NewUrlParameterBag(5, true)

	if bag.params != nil {
		t.Errorf("")
	}

	if bag.capacity != 5 {
		t.Errorf("")
	}

	if bag.reverse != true {
		t.Errorf("")
	}
}

func TestUrlParameterBagAddsParameter(t *testing.T) {
	bag := UrlParameterBag{}

	bag.Add("param1", "v1")

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

func TestUrlParameterBagAddsMultipleParameters(t *testing.T) {
	bag := UrlParameterBag{}

	bag.Add("param1", "v1")
	bag.Add("param2", "v2")

	if bag.params == nil {
		t.Errorf("")
	}

	if len(bag.params) != 2 {
		t.Errorf("")
	}

	if bag.params[0].name != "param1" {
		t.Errorf("")
	}
	if bag.params[0].value != "v1" {
		t.Errorf("")
	}

	if bag.params[1].name != "param2" {
		t.Errorf("")
	}

	if bag.params[1].value != "v2" {
		t.Errorf("")
	}
}

func TestUrlParameterBagGetByName(t *testing.T) {
	bag := UrlParameterBag{}

	bag.Add("param1", "v1")
	bag.Add("param2", "v2")
	bag.Add("param3", "v3")

	if bag.params == nil {
		t.Errorf("")
	}

	if len(bag.params) != 3 {
		t.Errorf("")
	}

	v, err := bag.GetByName("param1")
	if v != "v1" || err != nil {
		t.Errorf("")
	}

	v, _ = bag.GetByName("param2")
	if v != "v2" || err != nil {
		t.Errorf("")
	}

	v, _ = bag.GetByName("param3")
	if v != "v3" || err != nil {
		t.Errorf("")
	}

	v, err = bag.GetByName("param4")

	if v != "" || err == nil {
		t.Errorf("")
	}
}

func TestUrlParameterBagGetByNameInReverseMode(t *testing.T) {
	bag := UrlParameterBag{reverse: true}

	bag.Add("param1", "v1")
	bag.Add("param2", "v2")
	bag.Add("param3", "v3")
	bag.Add("param1", "v4")

	if bag.params == nil {
		t.Errorf("")
	}

	if len(bag.params) != 4 {
		t.Errorf("")
	}

	v, err := bag.GetByName("param1")
	if v != "v4" || err != nil {
		t.Errorf("")
	}

	v, _ = bag.GetByName("param2")
	if v != "v2" || err != nil {
		t.Errorf("")
	}

	v, _ = bag.GetByName("param3")
	if v != "v3" || err != nil {
		t.Errorf("")
	}

	v, err = bag.GetByName("param4")

	if v != "" || err == nil {
		t.Errorf("")
	}
}

func TestUrlParameterBagGetByIndex(t *testing.T) {
	bag := UrlParameterBag{}

	bag.Add("param1", "v1")
	bag.Add("param2", "v2")
	bag.Add("param3", "v3")

	if bag.params == nil {
		t.Errorf("")
	}

	if len(bag.params) != 3 {
		t.Errorf("")
	}

	v, err := bag.GetByIndex(0)
	if v != "v1" || err != nil {
		t.Errorf("")
	}

	v, err = bag.GetByIndex(1)
	if v != "v2" || err != nil {
		t.Errorf("")
	}

	v, err = bag.GetByIndex(2)
	if v != "v3" || err != nil {
		t.Errorf("")
	}

	v, err = bag.GetByIndex(3)
	if v != "" || err == nil {
		t.Errorf("")
	}
}

func TestUrlParameterBagGetByIndexInReverseMode(t *testing.T) {
	bag := UrlParameterBag{reverse: true}

	bag.Add("param3", "v3")
	bag.Add("param2", "v2")
	bag.Add("param1", "v1")

	if bag.params == nil {
		t.Errorf("")
	}

	if len(bag.params) != 3 {
		t.Errorf("")
	}

	v, err := bag.GetByIndex(0)
	if v != "v1" || err != nil {
		t.Errorf("")
	}

	v, err = bag.GetByIndex(1)
	if v != "v2" || err != nil {
		t.Errorf("")
	}

	v, err = bag.GetByIndex(2)
	if v != "v3" || err != nil {
		t.Errorf("")
	}

	v, err = bag.GetByIndex(3)
	if v != "" || err == nil {
		t.Errorf("")
	}
}
