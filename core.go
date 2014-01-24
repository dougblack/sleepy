package sleepy

import (
	"net/http"
)

type Resource interface {
	Get(map[string][]string) string
	Post(map[string][]string) string
	Put(map[string][]string) string
	Delete(map[string][]string) string
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

func (api *Api) dispatchRequest(request *http.Request, resource Resource) string {
	method := request.Method

	switch method {
	case "GET":
		return resource.Get(nil)
	case "POST":
		return resource.Post(nil)
	case "PUT":
		return resource.Put(nil)
	case "DELETE":
		return resource.Delete(nil)
	}
	return "Not implemented!"
}

func (api *Api) HandleRequest(request *http.Request) string {
	resource := api.matchResource(request.URL.Path)
	return api.dispatchRequest(request, resource)
}

func (api *Api) AddResource(resource Resource, path string) {
	api.routes = append(api.routes, Route{resource, path})
}
