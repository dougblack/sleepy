package main

import (
	"github.com/dougblack/sleepy"
	"net/url"
)

type Item struct{}

func (item Item) Get(values url.Values) (int, interface{}) {
	items := []string{"item1", "item2"}
	data := map[string][]string{"items": items}

	return 200, data
}

func main() {

	item := new(Item)

	var api = sleepy.NewAPI()
	api.AddResource(item, "/items", "/bar", "/baz")
	api.Start(3000)

}
