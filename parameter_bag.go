package http_router

const (
	ParamsBagKey = "urlParameters"
)

type urlParameter struct {
	name  string
	value string
}

type UrlParameterBag struct {
	params []urlParameter
}

func (u *UrlParameterBag) add(name, value string) {
	if u.params == nil {
		u.params = make([]urlParameter, 0, 5)
	}

	u.params = append(u.params, urlParameter{name, value})
}

func (u *UrlParameterBag) GetByName(name string, def string) string {
	for _, item := range u.params {
		if item.name == name {
			return item.value
		}
	}

	return def
}

func (u *UrlParameterBag) GetByIndex(index uint, def string) string {
	i := int(index)
	if len(u.params) <= i {
		return def
	}

	return u.params[i].value
}

func NewUrlParameterBag() UrlParameterBag {
	return UrlParameterBag{}
}