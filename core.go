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

type GetSupported interface {
	Get(values url.Values) (int, interface{})
}

type PostSupported interface {
	Post(values url.Values) (int, interface{})
}

type PutSupported interface {
	Put(values url.Values) (int, interface{})
}

type DeleteSupported interface {
	Delete(values url.Values) (int, interface{})
}

type API struct {
    mux *http.ServeMux
}

func NewAPI() *API {
    return &API{}
}

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
			if r, ok := resource.(GetSupported); ok {
				code, data = r.Get(values)
			} else {
                api.Abort(rw, 405)
                return
            }
		case POST:
			if r, ok := resource.(PostSupported); ok {
				code, data = r.Post(values)
			} else {
                api.Abort(rw, 405)
                return
            }
		case PUT:
			if r, ok := resource.(PutSupported); ok {
				code, data = r.Put(values)
			} else {
                api.Abort(rw, 405)
                return
            }
		case DELETE:
			if r, ok := resource.(DeleteSupported); ok {
				code, data = r.Delete(values)
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

func (api *API) AddResource(resource interface{}, path string) {
    if api.mux == nil {
        api.mux = http.NewServeMux()
    }
	api.mux.HandleFunc(path, api.requestHandler(resource))
}

func (api *API) Start(port int) error {
    if api.mux == nil {
        return &ErrorString{"You must add at least one resource to this API."}
    }
	portString := fmt.Sprintf(":%d", port)
	http.ListenAndServe(portString, api.mux)
    return nil
}
