package main

import (
	"net/url"

	"github.com/kid0m4n/sleepy"
)

type Item struct {
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
