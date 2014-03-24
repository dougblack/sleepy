package main

import (
	"github.com/akesling/sleepy"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

type Item struct{}

func (item Item) Get(values url.Values) (int, interface{}) {
	items := []string{"item1", "item2"}
	data := map[string][]string{"items": items}
	return 200, data
}

func TestBasicGet(t *testing.T) {

	item := new(Item)

	var api = sleepy.NewAPI()
	api.AddResource(item, "/items", "/bar", "/baz")
	go api.Start(3000)
	for path := range []string{"/items", "/bar", "/baz"} {
		resp, err := http.Get("http://localhost:3000" + path)
		if err != nil {
			t.Error(err)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) != `{"items":["item1","item2"]}` {
			t.Error("Not equal.")
		}
	}
}
