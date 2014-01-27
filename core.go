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

type geter interface {
	Get(values url.Values) (int, interface{})
}

type poster interface {
	Post(values url.Values) (int, interface{})
}

type puter interface {
	Put(values url.Values) (int, interface{})
}

type deleter interface {
	Delete(values url.Values) (int, interface{})
}

type Api struct{}

func (api *Api) Abort(rw http.ResponseWriter, statusCode int) {
	rw.WriteHeader(statusCode)
}

func (api *Api) requestHandler(resource interface{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		method := request.Method
		if request.ParseForm() != nil {
			api.Abort(rw, 400)
			return
		}
		values := request.Form

		var data interface{} = ""
		var code int = 405

		switch method {
		case GET:
			if r, ok := resource.(geter); ok {
				code, data = r.Get(values)
			}
		case POST:
			if r, ok := resource.(poster); ok {
				code, data = r.Post(values)
			}
		case PUT:
			if r, ok := resource.(puter); ok {
				code, data = r.Put(values)
			}
		case DELETE:
			if r, ok := resource.(deleter); ok {
				code, data = r.Delete(values)
			}
		default:
			api.Abort(rw, 405)
			return
		}

		responseWriter := json.NewEncoder(rw)
		rw.WriteHeader(code)
		if responseWriter.Encode(data) != nil {
			api.Abort(rw, 500)
			return
		}
	}
}

func (api *Api) AddResource(resource interface{}, path string) {
	http.HandleFunc(path, api.requestHandler(resource))
}

func (api *Api) Start(port int) {
	portString := fmt.Sprintf(":%d", port)
	http.ListenAndServe(portString, nil)
}
