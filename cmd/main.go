package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	checkpoint "github.com/aidenfine/checkpoint"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	})

	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "logs page")
	})

	// CONFIG METHOD ==
	config := checkpoint.Config{
		IgnorePaths:     []string{"/logs"},
		MaxTokens:       25,
		RefillRate:      1,
		TokensPerRefill: 1,
		LimitMethod:     checkpoint.LimitByIp,
	}

	// Quick method
	// rlMiddleware := checkpoint.LimitByIp(25, 1 , 1)

	rlMiddleware := checkpoint.WithConfig(config)
	rlHandler := rlMiddleware(http.DefaultServeMux)

	http.ListenAndServe(":8080", rlHandler)
}
