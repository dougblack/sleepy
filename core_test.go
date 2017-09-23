package sleepy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

type Item struct{}

func (item Item) Get(values url.Values, headers http.Header) (int, interface{}, http.Header) {
	items := []string{"item1", "item2"}
	data := map[string][]string{"items": items}
	return 200, data, nil
}

func (item Item) Post(values url.Values, headers http.Header) (int, interface{}, http.Header) {
	data := fmt.Sprintf("You sent: %s", values.Get("hello"))
	return 200, data, http.Header{"Content-Type": {"text/plain"}}
}

func TestBasicGet(t *testing.T) {

	item := new(Item)

	var api = NewAPI()
	api.AddResource(item, "/items", "/bar", "/baz")
	go api.Start(3000)
	resp, err := http.Get("http://localhost:3000/items")
	if err != nil {
		t.Error(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) != "{\n  \"items\": [\n    \"item1\",\n    \"item2\"\n  ]\n}" {
		t.Error("Not equal.")
	}
}

func TestBasicPostWithTextPlainResponse(t *testing.T) {

	item := new(Item)

	var api = NewAPI()
	api.AddResource(item, "/items", "/bar", "/baz")
	go api.Start(3000)
	resp, err := http.PostForm("http://localhost:3000/items", url.Values{"hello": {"sleepy"}})
	if err != nil {
		t.Error(err)
	}
	if resp.Header.Get("Content-Type") != "text/plain" {
		t.Error("Content-Type wrong.")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) != "You sent: sleepy" {
		t.Error("Not equal.")
	}

}
