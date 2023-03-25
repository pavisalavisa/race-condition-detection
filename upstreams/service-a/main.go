package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/v1/foo/bar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unexpected HTTP method")
			return
		}
		fmt.Fprintf(w, "Hello from service A")
	})

	if err := http.ListenAndServe(":9202", nil); err != nil {
		panic(err)
	}
}
