package sleepy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type Resource interface {
	Get(values ...url.Values) (int, interface{})
	Post(values ...url.Values) (int, interface{})
	Put(values ...url.Values) (int, interface{})
	Delete(values ...url.Values) (int, interface{})
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

func (api *Api) Abort(statusCode int) (int, interface{}) {
	return statusCode, map[string]string{"error": "Aborted."}
}

type HandleFunc func(http.ResponseWriter, *http.Request)

func (api *Api) requestHandler(resource Resource) HandleFunc {
	return func(rw http.ResponseWriter, request *http.Request) {

		var code int
		var data interface{}
		var content []byte

		method := request.Method

		if request.ParseForm() == nil {
			code, data = api.Abort(500)
		}

		values := request.Form

		switch method {
		case GET:
			code, data = resource.Get(values)
		case POST:
			code, data = resource.Post(values)
		case PUT:
			code, data = resource.Put(values)
		case DELETE:
			code, data = resource.Delete(values)
		default:
			code, data = 405, map[string]string{"error": "Not implemented!"}
		}

		content, err := json.Marshal(data)
		if err != nil {
			content, _ = json.Marshal(map[string]string{"error": "Bad response."})
		}

		rw.WriteHeader(code)
		rw.Write(content)
	}
}

func (api *Api) AddResource(resource Resource, path string) {
	http.HandleFunc(path, api.requestHandler(resource))
}

func (api *Api) Start(port int) {
	portString := fmt.Sprintf(":%d", port)
	http.ListenAndServe(portString, nil)
}
