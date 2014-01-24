package sleepy

type Resource interface {
    Get() string
    // Post
    // Put 
    // Delete
}

type Route struct {
    resource Resource
    path string
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

func (api *Api) HandleRequest(path string) string {
    resource := api.matchResource(path)
    return resource.Get()
}


func (api *Api) AddResource(resource Resource, path string) {
    api.routes = append(api.routes, Route{resource, path})
}
