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
	Get(values url.Values) (int, interface{})
	Post(values url.Values) (int, interface{})
	Put(values url.Values) (int, interface{})
	Delete(values url.Values) (int, interface{})
}

type API struct {
    mux *http.ServeMux
}

func (api *API) Abort(rw http.ResponseWriter, statusCode int) {
	rw.WriteHeader(statusCode)
	rw.Write([]byte(http.StatusText(statusCode)))
}


func (api *API) requestHandler(resource Resource) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {

		var data interface{}
		var code int

		method := request.Method
		if request.ParseForm() != nil {
			api.Abort(rw, 400)
			return
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

func (api *API) AddResource(resource Resource, path string) {
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
	fmt.Println("Hi.")
    return nil
}
