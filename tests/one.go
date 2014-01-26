package main

import (
	"net/http"
)

func response(rw http.ResponseWriter, request *http.Request) {
	rw.Write([]byte("Hello world."))
}

func main() {
	http.HandleFunc("/", response)
	http.ListenAndServe(":3000", nil)
}
