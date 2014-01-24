package sleepy

import (
	"net/http"
)

type Resource interface {
	Get() string
}

type Route struct {
	resource Resource
	path     string
}

func (route *Route) pathMatch(path string) bool {
	return route.path == path
}

type Api struct {
	routes []Route
}

func (api *Api) matchResource(path string) Resource {
	for _, route := range api.routes {
		if route.pathMatch(path) {
			return route.resource
		}
	}
	return nil
}

func (api *Api) HandleRequest(request *http.Request) string {
	resource := api.matchResource(request.URL.Path)
	return resource.Get()
}

func (api *Api) AddResource(resource Resource, path string) {
	api.routes = append(api.routes, Route{resource, path})
}
