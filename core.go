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

type BaseResource struct {
  GetNotSupported
  PostNotSupported
  PutNotSupported
  DeleteNotSupported
}

type Api struct{}

func (api *Api) Abort(rw http.ResponseWriter, statusCode int) {
	rw.WriteHeader(statusCode)
}

type HandleFunc func(http.ResponseWriter, *http.Request)

func (api *Api) requestHandler(resource Resource) HandleFunc {
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
			return
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
	fmt.Println("Hi.")
}
