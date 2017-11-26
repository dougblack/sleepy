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
	HEAD   = "HEAD"
	PATCH  = "PATCH"
)

type Resource interface {
	Get(url.Values, http.Header) (int, interface{}, http.Header)
	Post(url.Values, http.Header) (int, interface{}, http.Header)
	Put(url.Values, http.Header) (int, interface{}, http.Header)
	Delete(url.Values, http.Header) (int, interface{}, http.Header)
	Head(url.Values, http.Header) (int, interface{}, http.Header)
	Patch(url.Values, http.Header) (int, interface{}, http.Header)
}

type NotSupported struct {}

func (NotSupported) Get(url.Values, http.Header) (int, interface{}, http.Header) {
				return http.StatusMethodNotAllowed, nil, nil
}

func (NotSupported) Post(url.Values, http.Header) (int, interface{}, http.Header) {
				return http.StatusMethodNotAllowed, nil, nil
}

func (NotSupported) Put(url.Values, http.Header) (int, interface{}, http.Header) {
				return http.StatusMethodNotAllowed, nil, nil
}

func (NotSupported) Delete(url.Values, http.Header) (int, interface{}, http.Header) {
				return http.StatusMethodNotAllowed, nil, nil
}

func (NotSupported) Patch(url.Values, http.Header) (int, interface{}, http.Header) {
				return http.StatusMethodNotAllowed, nil, nil
}

func (NotSupported) Head(url.Values, http.Header) (int, interface{}, http.Header) {
				return http.StatusMethodNotAllowed, nil, nil
}


// An API manages a group of resources by routing requests
// to the correct method on a matching resource and marshalling
// the returned data to JSON for the HTTP response.
//
// You can instantiate multiple APIs on separate ports. Each API
// will manage its own set of resources.
type API struct {
	mux     *http.ServeMux
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

		var handler func(url.Values, http.Header) (int, interface{}, http.Header)

		switch request.Method {
		case GET:
				handler = resource.Get
		case POST:
				handler = resource.Post
		case PUT:
				handler = resource.Put
		case DELETE:
				handler = resource.Delete
		case HEAD:
				handler = resource.Head
		case PATCH:
				handler = resource.Patch
		}

		code, data, header := handler(request.Form, request.Header)

		if code < http.StatusOK || code > 300  {
			rw.WriteHeader(code)
			return
		}

		content, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		for name, values := range header {
			for _, value := range values {
				rw.Header().Add(name, value)
			}
		}
		rw.WriteHeader(code)
		rw.Write(content)
	}
}

// Mux returns the http.ServeMux used by an API. If a ServeMux has
// does not yet exist, a new one will be created and returned.
func (api *API) Mux() *http.ServeMux {
	if api.muxInitialized {
		return api.mux
	} else {
		api.mux = http.NewServeMux()
		api.muxInitialized = true
		return api.mux
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

// AddResourceWithWrapper behaves exactly like AddResource but wraps
// the generated handler function with a give wrapper function to allow
// to hook in Gzip support and similar.
func (api *API) AddResourceWithWrapper(resource Resource, wrapper func(handler http.HandlerFunc) http.HandlerFunc, paths ...string) {
	for _, path := range paths {
		api.Mux().HandleFunc(path, wrapper(api.requestHandler(resource)))
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
