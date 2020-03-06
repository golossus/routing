package http_router

import (
	"fmt"
)

type urlParameter struct {
	name  string
	value string
}

type UrlParameterBag struct {
	params   []urlParameter
	capacity uint
	reverse  bool
}

func (u *UrlParameterBag) Add(name, value string) {
	if u.params == nil {
		u.params = make([]urlParameter, 0, u.capacity)
	}

	u.params = append(u.params, urlParameter{name, value})
}

func (u *UrlParameterBag) GetByName(name string) (string, error) {
	for i := range u.params {
		if u.reverse {
			i = len(u.params) - 1 - i
		}
		if u.params[i].name == name {
			return u.params[i].value, nil
		}
	}

	return "", fmt.Errorf("url parameter with name %s does not exist", name)
}

func (u *UrlParameterBag) GetByIndex(index uint) (string, error) {
	i := int(index)
	if len(u.params) <= i {
		return "", fmt.Errorf("url parameter at index %d does not exist", i)
	}

	if u.reverse {
		i = len(u.params) - 1 - i
	}

	return u.params[i].value, nil
}

func NewUrlParameterBag(capacity uint, reverse bool) UrlParameterBag {
	return UrlParameterBag{
		capacity: capacity,
		reverse:  reverse,
	}
}
