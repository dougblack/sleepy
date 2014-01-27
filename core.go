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

type API struct {
    mux *http.ServeMux
}

func (api *API) Abort(rw http.ResponseWriter, statusCode int) {
	rw.WriteHeader(statusCode)
	rw.Write([]byte(http.StatusText(statusCode)))
}

func (api *API) requestHandler(resource interface{}) http.HandlerFunc {
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

		content, err := json.Marshal(data)
		if err != nil {
			api.Abort(rw, 500)
		}
		rw.WriteHeader(code)
		rw.Write(content)
	}
}

func (api *API) AddResource(resource interface{}, path string) {
    if api.mux == nil {
        api.mux = http.NewServeMux()
    }
	api.mux.HandleFunc(path, api.requestHandler(resource))
}

func (api *API) Start(port int) error {
    if api.mux == nil {
        return &errorString{"You must add at last one resource to this API."}
    }
	portString := fmt.Sprintf(":%d", port)
	http.ListenAndServe(portString, api.mux)
    return nil
}
