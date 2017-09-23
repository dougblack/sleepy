package sleepy

import (
	"io/ioutil"
	"net/http"
	"testing"
)

type Item struct{}

func (item Item) Get(*http.Request) (int, interface{}, http.Header) {
	items := []string{"item1", "item2"}
	data := map[string][]string{"items": items}
	return 200, data, nil
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
