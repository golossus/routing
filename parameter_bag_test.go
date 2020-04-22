package routing

import "testing"

func assertBagSetting(t *testing.T, bag URLParameterBag, cap uint, reverse bool) {
	if bag.params != nil {
		t.Errorf("parameter bag params is not nil")
	}

	if bag.capacity != cap {
		t.Errorf("parameter bag capacity %d not equals to %d", bag.capacity, cap)
	}

	if bag.reverse != reverse {
		t.Errorf("parameter bag reverse mode %t not equals to %t", bag.reverse, reverse)
	}
}

func assertBagParameterContains(t *testing.T, bag URLParameterBag, name string, value string) {
	v, err := bag.GetByName(name)

	if err != nil {
		t.Errorf("parameter %s not found", name)
	}

	if v != value {
		t.Errorf("parameter %s with value %s not equals to %s", name, v, value)
	}
}

func assertBagParameterAtIndex(t *testing.T, bag URLParameterBag, index uint, value string) {
	v, err := bag.GetByIndex(index)

	if err != nil {
		t.Errorf("parameter at index %d not found", index)
	}

	if v != value {
		t.Errorf("parameter at index %d with value %s not equals to %s", index, v, value)
	}
}

func assertBagParameterNotContains(t *testing.T, bag URLParameterBag, name string) {
	_, err := bag.GetByName(name)

	if err == nil {
		t.Errorf("parameter %s found", name)
	}
}

func assertBagParameterNotAtIndex(t *testing.T, bag URLParameterBag, index uint) {
	_, err := bag.GetByIndex(index)

	if err == nil {
		t.Errorf("parameter at index %d found", index)
	}
}

func TestUrlParameterBagEmptyValuesOnCreation(t *testing.T) {
	bag := URLParameterBag{}

	assertBagSetting(t, bag, 0, false)
}

func TestNewUrlParameterBagSetsRightValuesOnCreation(t *testing.T) {
	bag := newURLParameterBag(5, true)

	assertBagSetting(t, bag, 5, true)
}

func TestUrlParameterBagAddsParameter(t *testing.T) {
	bag := URLParameterBag{}

	bag.add("param1", "v1")

	assertBagParameterContains(t, bag, "param1", "v1")
}

func TestUrlParameterBagAddsMultipleParameters(t *testing.T) {
	bag := URLParameterBag{}

	bag.add("param1", "v1")
	bag.add("param2", "v2")
	bag.add("param3", "v3")

	assertBagParameterContains(t, bag, "param1", "v1")
	assertBagParameterContains(t, bag, "param2", "v2")
	assertBagParameterContains(t, bag, "param3", "v3")
	assertBagParameterNotContains(t, bag, "param4")
}

func TestUrlParameterBagGetByNameInReverseMode(t *testing.T) {
	bag := URLParameterBag{reverse: true}

	bag.add("param1", "v1")
	bag.add("param2", "v2")
	bag.add("param3", "v3")
	bag.add("param1", "v4")

	assertBagParameterContains(t, bag, "param1", "v4")
	assertBagParameterContains(t, bag, "param2", "v2")
	assertBagParameterContains(t, bag, "param3", "v3")
	assertBagParameterNotContains(t, bag, "param4")
}

func TestUrlParameterBagGetByIndex(t *testing.T) {
	bag := URLParameterBag{}

	bag.add("param1", "v1")
	bag.add("param2", "v2")
	bag.add("param3", "v3")

	assertBagParameterAtIndex(t, bag, 0, "v1")
	assertBagParameterAtIndex(t, bag, 1, "v2")
	assertBagParameterAtIndex(t, bag, 2, "v3")
	assertBagParameterNotAtIndex(t, bag, 3)
}

func TestUrlParameterBagGetByIndexInReverseMode(t *testing.T) {
	bag := URLParameterBag{reverse: true}

	bag.add("param3", "v3")
	bag.add("param2", "v2")
	bag.add("param1", "v1")

	assertBagParameterAtIndex(t, bag, 0, "v1")
	assertBagParameterAtIndex(t, bag, 1, "v2")
	assertBagParameterAtIndex(t, bag, 2, "v3")
	assertBagParameterNotAtIndex(t, bag, 3)
}
