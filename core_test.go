package sleepy_test

import (
	"github.com/akesling/sleepy"
	"io/ioutil"
	"net/http"
	"testing"
)

type Item struct{}

func (item Item) Get(req *http.Request) (int, interface{}) {
	items := []string{"item1", "item2"}
	data := map[string][]string{"items": items}
	return 200, data
}

func TestBasicGet(t *testing.T) {
	api := sleepy.NewAPI()
	api.AddResource(new(Item), "/items", "/bar", "/baz")
	go api.Start(3000)
	paths := []string{"/items", "/bar", "/baz"}
	for i := range paths {
		resp, err := http.Get("http://localhost:3000" + paths[i])
		if err != nil {
			t.Error(err)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		actual := string(body)
		expected := `{"items":["item1","item2"]}`
		if actual != expected {
			t.Error("\nActual:\n" + actual +
				"Does not equal expected:\n" + expected)
		}
	}
}
