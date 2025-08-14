package checkpoint

import (
	"log"
	"net/http"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = forwarded
		}

		log.Printf("Request from IP: %s | Method: %s | Path: %s", ip, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
