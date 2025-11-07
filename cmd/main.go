package main

import (
	"fmt"
	"net/http"

	checkpointmiddleware "github.com/aidenfine/checkpoint-middleware"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	})

	rlMiddleware := checkpointmiddleware.LimitByIp(100, 1, 1)
	rlHandler := rlMiddleware(mux)

	http.ListenAndServe(":8080", rlHandler)
}
