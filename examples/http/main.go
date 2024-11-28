package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.ListenAndServe("127.0.0.1:8000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, r.Host, r.URL)
		w.WriteHeader(http.StatusOK)
	}))
}
