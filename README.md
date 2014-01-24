## Sleepy

#### A RESTful framework for Go

Sleepy is not done yet.  Here is a potential target API.

```go

import (
    "net/http"
    "sleepy"
)

type Item struct { }

func (item *Item) Get(foo string, bar int) (interface{}, int, http.Headers) {
    data := map[string]int {
        foo : bar
    }
    return data, 200, nil
}

func main() {

    item = new(Item)

    var api = new(sleepy.Api)
    api.AddResource(item, "/item")

    request, _ := http.NewRequest("GET", "/item", "foo=thing&bar=5")
    fmt.Println(api.HandleRequest(request))

}
```

With a response of

``javascript
{
    "thing": 5
}
```

Stay tuned.
