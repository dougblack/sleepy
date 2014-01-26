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

type GetResource interface {
	Get(values ...url.Values) (int, interface{})
}

type PostResource interface {
	Post(values ...url.Values) (int, interface{})
}

type PutResource interface {
	Put(values ...url.Values) (int, interface{})
}

type DeleteResource interface {
	Delete(values ...url.Values) (int, interface{})
}


type Api struct{}

func (api *Api) Abort(rw http.ResponseWriter, statusCode int) {
	rw.WriteHeader(statusCode)
}

type HandleFunc func(http.ResponseWriter, *http.Request)

func (api *Api) requestHandler(resource interface{}) HandleFunc {
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
		    getter, ok := resource.(GetResource)
		    if ok {
    			code, data = getter.Get(values)
			} else {
			    code, data = 405, ""
		    }
		case POST:
		    poster, ok := resource.(PostResource)
		    if ok {
    			code, data = poster.Post(values)
			} else {
			    code, data = 405, ""
		    }
		case PUT:
		    putter, ok := resource.(PutResource)
		    if ok {
    			code, data = putter.Put(values)
			} else {
			    code, data = 405, ""
		    }
		case DELETE:
		    deleter, ok := resource.(DeleteResource)
		    if ok {
    			code, data = deleter.Delete(values)
			} else {
			    code, data = 405, ""
		    }
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

func (api *Api) AddResource(resource interface{}, path string) {
	http.HandleFunc(path, api.requestHandler(resource))
}

func (api *Api) Start(port int) {
	portString := fmt.Sprintf(":%d", port)
	http.ListenAndServe(portString, nil)
	fmt.Println("Hi.")
}
