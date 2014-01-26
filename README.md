## Sleepy

#### A RESTful framework for Go

Sleepy is a micro-framework for building RESTful APIs.

```go
package main

import (
    "net/url"
    "github.com/dougblack/sleepy"
)

type Item struct {
    sleepy.PostNotSupported
    sleepy.PutNotSupported
    sleepy.DeleteNotSupported
}

func (item Item) Get(values url.Values) (int, interface{}) {
    items := []string{"item1", "item2"}
    data := map[string][]string{"items": items}

    return 200, data
}

func main() {

    item := new(Item)

    var api = new(sleepy.Api)
    api.AddResource(item, "/items")
    api.Start(3000)

}
```

Now if we curl that endpoint:

```bash
curl localhost:3000/items
{"items": ["item1", "item2"]}
```

Stay tuned.
