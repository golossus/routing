package routing

import (
	"fmt"
)

type urlParameter struct {
	name  string
	value string
}

// URLParameterBag is a structure where the URL parameters are saved
type URLParameterBag struct {
	params   []urlParameter
	capacity uint
}

func (u *URLParameterBag) add(name, value string) {
	if u.params == nil {
		u.params = make([]urlParameter, 0, u.capacity)
	}

	u.params = append(u.params, urlParameter{name, value})
}

// GetByName is a method to retrieve a dynamic parameter of the URL using a name. For example 'userId' in /users/{userId}
func (u *URLParameterBag) GetByName(name string) (string, error) {
	for i := range u.params {
		if u.params[i].name == name {
			return u.params[i].value, nil
		}
	}

	return "", fmt.Errorf("url parameter with name %s does not exist", name)
}

// GetByIndex is a method to retrieve a dynamic parameter of the URL using a index. For example 1 to obtain the 'userId' in /users/{userId}/file/{fileId}
func (u *URLParameterBag) GetByIndex(index uint) (string, error) {
	i := int(index)
	if len(u.params) <= i {
		return "", fmt.Errorf("url parameter at index %d does not exist", i)
	}

	return u.params[i].value, nil
}

func newURLParameterBag(capacity uint) URLParameterBag {
	return URLParameterBag{
		capacity: capacity,
	}
}
