package main

import (
    "fmt"
    "sleepy"
    "net/http"
)

type Bar struct {
    sleepy.PostNotSupported
    sleepy.PutNotSupported
    sleepy.DeleteNotSupported
}

func (b Bar) Get(map[string][]string) string {
    return "Hello"
}

type Baz struct {
    sleepy.PostNotSupported
    sleepy.PutNotSupported
    sleepy.DeleteNotSupported
}

func (b Baz) Get(map[string][]string) string {
    return "Goodbye"
}

func main() {
    bar := new(Bar)
    baz := new(Baz)

    var api = new(sleepy.Api)
    api.AddResource(bar, "/bar")
    api.AddResource(baz, "/baz")

    request1, _ := http.NewRequest("GET", "https://dougblack.io/bar", nil)
    request2, _ := http.NewRequest("GET", "https://dougblack.io/baz", nil)
    fmt.Println(api.HandleRequest(request1))
    fmt.Println(api.HandleRequest(request2))

}

