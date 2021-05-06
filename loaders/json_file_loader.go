package loaders

import (
	"encoding/json"
	"github.com/golossus/routing"
	"io/ioutil"
)

// JsonRoutes defines a json collection of routes
type JsonRoutes struct {
	Routes []JsonRoute `json:"routes"`
}

// JsonRoute defines a route json schema
type JsonRoute struct {
	Name    string `json:"name"`
	Method  string `json:"method"`
	Path  string `json:"path"`
	Host string `json:"host"`
	Handler string `json:"handler"`
	Schemas []string `json:"schemas"`
	Headers map[string]string `json:"headers"`
	QueryParams map[string]string `json:"queryParams"`
	CustomMatcher string `json:"customMatcher"`
}

// JsonFileLoader type loads routes from Json files
type JsonFileLoader struct {
	routes []routing.RouteDef
}

// Load implements routing.Loader interface
func (l *JsonFileLoader) Load() []routing.RouteDef {
	return l.routes
}

// FromFile loads a list of routes from one or many Json file paths
func (l *JsonFileLoader) FromFile(files ...string) error {

	for _, path := range files {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		var rs JsonRoutes
		err = json.Unmarshal(content, &rs)
		if err != nil {
			return err
		}

		fileRoutes := make([]routing.RouteDef, 0, len(rs.Routes))
		for _, r := range rs.Routes {
			fileRoutes = append(fileRoutes, routing.RouteDef{
				Method:  r.Method,
				Path:  r.Path,
				Handler: r.Handler,
				Options:  routing.RouteDefOptions{
					Name: r.Name,
					Host: r.Host,
					Schemas: r.Schemas,
					Headers: r.Headers,
					QueryParams: r.QueryParams,
					CustomMatcher: r.CustomMatcher,
				},
			})
		}

		l.routes = append(l.routes, fileRoutes...)

	}

	return nil
}