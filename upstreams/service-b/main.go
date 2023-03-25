package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/baz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unexpected HTTP method")
			return
		}
		fmt.Fprintf(w, "Hello from service B")
	})

	if err := http.ListenAndServe(":9101", nil); err != nil {
		panic(err)
	}
}
