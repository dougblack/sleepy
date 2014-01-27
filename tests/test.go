package main

import (
	"net/url"
	"github.com/dougblack/sleepy"
)

type Item struct {}

func (item Item) Get(values url.Values) (int, interface{}) {
	items := []string{"item1", "item2"}
	data := map[string][]string{"items": items}

	return 200, data
}

func main() {

	item := new(Item)

	var api = sleepy.NewAPI()
	api.AddResource(item, "/items")
	api.Start(3000)

}
