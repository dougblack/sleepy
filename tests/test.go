package main

import (
    "sleepy"
    "net/url"
)

type Bar struct {
    sleepy.PostNotSupported
    sleepy.PutNotSupported
    sleepy.DeleteNotSupported
}

func (b Bar) Get(values ...url.Values) (int, interface{}) {
    return 200, map[string]string{"hello": "goodbye"}
}

type Baz struct {
    sleepy.PostNotSupported
    sleepy.PutNotSupported
    sleepy.DeleteNotSupported
}

func main() {
    bar := new(Bar)

    var api = new(sleepy.Api)
    api.AddResource(bar, "/bar")
    api.Start(3000)
}

