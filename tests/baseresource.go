package main

import (
	"net/url"
	"github.com/rakoo/sleepy"
)

type Bar sleepy.BaseResource

func (b Bar) Get(values ...url.Values) (int, interface{}) {
	return 200, map[string]string{"hello": "goodbye"}
}

func main() {
	bar := new(Bar)

	var api = new(sleepy.Api)
	api.AddResource(bar, "/bar")
	api.Start(3000)
}
