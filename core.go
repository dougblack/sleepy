package sleepy

import (
	"encoding/json"
	"errors"
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

// GetSupported is the interface that provides the Get
// method a resource must support to receive HTTP GETs.
type GetSupported interface {
	Get(values url.Values) (int, interface{})
}

// PostSupported is the interface that provides the Post
// method a resource must support to receive HTTP POST.
type PostSupported interface {
	Post(values url.Values) (int, interface{})
}

// PutSupported is the interface that provides the Put
// method a resource must support to receive HTTP PUT.
type PutSupported interface {
	Put(values url.Values) (int, interface{})
}

// DeleteSupported is the interface that provides the Delete
// method a resource must support to receive HTTP DELETE.
type DeleteSupported interface {
	Delete(values url.Values) (int, interface{})
}

// An API manages a group of resources by routing to requests
// to the correct method on a matching resource.
//
// You can instantiate multiple APIs on separate ports. Each API
// will manage its own set of resources.
type API struct {
	mux *http.ServeMux
}

// NewAPI allocates and returns a new API.
func NewAPI() *API {
	return &API{}
}

// Abort responds to a given request with the status text associated
// with the passed in HTTP status code.
func (api *API) Abort(rw http.ResponseWriter, statusCode int) {
	rw.WriteHeader(statusCode)
	rw.Write([]byte(http.StatusText(statusCode)))
}

func (api *API) requestHandler(resource interface{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {

		var err error
		var data interface{} = ""
		var code int = 405

		method := request.Method
		if request.ParseForm() != nil {
			api.Abort(rw, 400)
			return
		}
		values := request.Form

		switch method {
		case GET:
			if resource, ok := resource.(GetSupported); ok {
				code, data = resource.Get(values)
			} else {
				api.Abort(rw, 405)
				return
			}
		case POST:
			if resource, ok := resource.(PostSupported); ok {
				code, data = resource.Post(values)
			} else {
				api.Abort(rw, 405)
				return
			}
		case PUT:
			if resource, ok := resource.(PutSupported); ok {
				code, data = resource.Put(values)
			} else {
				api.Abort(rw, 405)
				return
			}
		case DELETE:
			if resource, ok := resource.(DeleteSupported); ok {
				code, data = resource.Delete(values)
			} else {
				api.Abort(rw, 405)
				return
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

// AddResource adds a new resource to an API. The API will route
// request matching the path to the correct HTTP method on the
// resource.
func (api *API) AddResource(resource interface{}, path string) {
	if api.mux == nil {
		api.mux = http.NewServeMux()
	}
	api.mux.HandleFunc(path, api.requestHandler(resource))
}

// Start causes the API to begin serving requests on the given port.
func (api *API) Start(port int) error {
	if api.mux == nil {
		return errors.New("You must add at least one resource to this API.")
	}
	portString := fmt.Sprintf(":%d", port)
	http.ListenAndServe(portString, api.mux)
	return nil
}
