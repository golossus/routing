package loaders

import (
	"encoding/json"
	. "github.com/golossus/routing"
	"io/ioutil"
)

// JsonRoutes defines a json collection of routes
type JsonRoutes struct {
	Routes []JsonRoute `json:"routes"`
}

// JsonRoute defines a route json schema
type JsonRoute struct {
	Name    string `json:"name"`
	Schema  string `json:"schema"`
	Method  string `json:"method"`
	Handler string `json:"handler"`
}

// JsonFileLoader type loads routes from Json files
type JsonFileLoader struct {
	routes []RouteDef
}

// Load implements routing.Loader interface
func (l *JsonFileLoader) Load() []RouteDef {
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

		fileRoutes := make([]RouteDef, 0, len(rs.Routes))
		for _, r := range rs.Routes {
			fileRoutes = append(fileRoutes, RouteDef{
				Method:  r.Method,
				Schema:  r.Schema,
				Handler: r.Handler,
				Name:    r.Name,
			})
		}

		l.routes = append(l.routes, fileRoutes...)

	}

	return nil
}