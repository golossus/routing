package routing

import "testing"

func assertBagSetting(t *testing.T, bag URLParameterBag, cap uint) {
	if bag.params != nil {
		t.Errorf("parameter bag params is not nil")
	}

	if bag.capacity != cap {
		t.Errorf("parameter bag capacity %d not equals to %d", bag.capacity, cap)
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

func TestURLParameterBag_EmptyValuesOnCreation(t *testing.T) {
	bag := URLParameterBag{}

	assertBagSetting(t, bag, 0)
}

func TestURLParameterBag_SetsRightValuesOnCreation(t *testing.T) {
	bag := newURLParameterBag(5)

	assertBagSetting(t, bag, 5)
}

func TestURLParameterBag_Add_Works(t *testing.T) {
	bag := URLParameterBag{}

	bag.add("param1", "v1")

	assertBagParameterContains(t, bag, "param1", "v1")
}

func TestURLParameterBag_GetByName(t *testing.T) {
	bag := URLParameterBag{}

	bag.add("param1", "v1")
	bag.add("param2", "v2")
	bag.add("param3", "v3")

	assertBagParameterContains(t, bag, "param1", "v1")
	assertBagParameterContains(t, bag, "param2", "v2")
	assertBagParameterContains(t, bag, "param3", "v3")
	assertBagParameterNotContains(t, bag, "param4")
}

func TestURLParameterBag_GetByIndex(t *testing.T) {
	bag := URLParameterBag{}

	bag.add("param1", "v1")
	bag.add("param2", "v2")
	bag.add("param3", "v3")

	assertBagParameterAtIndex(t, bag, 0, "v1")
	assertBagParameterAtIndex(t, bag, 1, "v2")
	assertBagParameterAtIndex(t, bag, 2, "v3")
	assertBagParameterNotAtIndex(t, bag, 3)
}

func TestURLParameterBag_merge(t *testing.T) {
	bag := URLParameterBag{}

	bag.add("param1", "v1")
	bag.add("param2", "v2")
	bag.add("param3", "v3")

	bag2 := URLParameterBag{}

	bag2.add("p1", "p1")
	bag2.add("p2", "p2")
	bag2.add("p3", "p3")

	merged := bag.merge(bag2)

	assertBagParameterAtIndex(t, merged, 0, "v1")
	assertBagParameterAtIndex(t, merged, 1, "v2")
	assertBagParameterAtIndex(t, merged, 2, "v3")
	assertBagParameterAtIndex(t, merged, 3, "p1")
	assertBagParameterAtIndex(t, merged, 4, "p2")
	assertBagParameterAtIndex(t, merged, 5, "p3")
	assertBagParameterNotAtIndex(t, bag, 6)
}
