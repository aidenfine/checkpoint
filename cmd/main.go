package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/aidenfine/checkpoint"
)

func main() {

	checkpoint.LoadConfig()
	cfg := checkpoint.GetConfig()

	if cfg.ServiceUrl == "" {
		panic("missing SERVICE_URL env")
	}
	if cfg.Port == "" {
		panic("missing PORT env")
	}

	backendURL, _ := url.Parse(cfg.ServiceUrl)
	proxy := httputil.NewSingleHostReverseProxy(backendURL)

	limiter := checkpoint.NewTokenBucket(3, 75, 5)

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

	fmt.Printf("Reverse proxy :%s → forwarding to :%s \n", cfg.Port, cfg.ServiceUrl)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
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
