package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/aidenfine/checkpoint"
)

func main() {

	serviceUrl := os.Getenv("SERVICE_URL")

	if serviceUrl == "" {
		panic("missing SERVICE_URL in env!")
	}
	backendURL, _ := url.Parse(serviceUrl)
	proxy := httputil.NewSingleHostReverseProxy(backendURL)

	limiter := checkpoint.NewTokenBucket(3, 100, 5)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)
		allowed, remaining := limiter.Allow(ip)
		if !allowed {
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprintf(w, "Rate limit exceeded! Try again later.\n")
			return
		}

		r.Header.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		fmt.Printf("[%s] allowed: %d tokens left\n", ip, remaining)

		proxy.ServeHTTP(w, r)
	})

	fmt.Println("Reverse proxy :8080 → forwarding to :8000")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		parts := strings.Split(ip, ",")
		return strings.TrimSpace(parts[0])
	}

	ip = strings.Split(r.RemoteAddr, ":")[0]
	return ip
}
