package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aidenfine/checkpoint"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hello world")
	})
	rlMw := checkpoint.LimitByIp(15, 30*time.Second)
	handler := rlMw(mux)

	http.ListenAndServe(":8080", handler)
}
