package main

import (
	"net/http"
)

func main() {

	proxy, err := NewProxy("http://service-a:9202", "http://service-b:9101")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", proxy.ServeHTTP)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
