package sleepy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type Resource interface {
	// returns pointer of an instance of a model
	ModelInstance() interface{}

	// returns the length of resource path
	// example: len("/books/")
	PathLength() int
}

// GetSupported is the interface that provides the Get
// method a resource must support to receive HTTP GETs.
type GetSupported interface {
	Get(instance interface{}, id string) (int, interface{})
}

// PostSupported is the interface that provides the Post
// method a resource must support to receive HTTP POSTs.
type PostSupported interface {
	Post(instance interface{}, id string) (int, interface{})
}

// PutSupported is the interface that provides the Put
// method a resource must support to receive HTTP PUTs.
type PutSupported interface {
	Put(instance interface{}, id string) (int, interface{})
}

// DeleteSupported is the interface that provides the Delete
// method a resource must support to receive HTTP DELETEs.
type DeleteSupported interface {
	Delete(instance interface{}, id string) (int, interface{})
}

// An API manages a group of resources by routing to requests
// to the correct method on a matching resource and marshalling
// the returned data to JSON for the HTTP response.
//
// You can instantiate multiple APIs on separate ports. Each API
// will manage its own set of resources.
type API struct {
	muxPointer     *http.ServeMux
	muxInitialized bool
}

// NewAPI allocates and returns a new API.
func NewAPI() *API {
	return &API{}
}

func (api *API) requestHandler(resource Resource) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {

		if request.ParseForm() != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		var handler func(instance interface{}, id string) (int, interface{})

		switch request.Method {
		case GET:
			if resource, ok := resource.(GetSupported); ok {
				handler = resource.Get
			}
		case POST:
			if resource, ok := resource.(PostSupported); ok {
				handler = resource.Post
			}
		case PUT:
			if resource, ok := resource.(PutSupported); ok {
				handler = resource.Put
			}
		case DELETE:
			if resource, ok := resource.(DeleteSupported); ok {
				handler = resource.Delete
			}
		}

		if handler == nil {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			// error handling
		}

		instance := resource.ModelInstance()
		json.Unmarshal(body, instance)

		id := request.URL.Path[resource.PathLength():]
		code, data := handler(instance, id)

		content, err := json.Marshal(data)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(code)
		rw.Write(content)
	}
}

// singleton mux
func (api *API) Mux() *http.ServeMux {
	if api.muxInitialized {
		return api.muxPointer
	} else {
		api.muxPointer = http.NewServeMux()
		api.muxInitialized = true
		return api.muxPointer
	}
}

// AddResource adds a new resource to an API. The API will route
// requests that match one of the given paths to the matching HTTP
// method on the resource.
func (api *API) AddResource(resource Resource, paths ...string) {
	for _, path := range paths {
		api.Mux().HandleFunc(path, api.requestHandler(resource))
	}
}

// Start causes the API to begin serving requests on the given port.
func (api *API) Start(port int) error {
	if !api.muxInitialized {
		return errors.New("You must add at least one resource to this API.")
	}
	portString := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(portString, api.Mux())
}
