## Checkpoint
Small go rate limiter that can be used inside middleware.

## Example
run `go get github.com/aidenfine/checkpoint`

```go
package main

import (
	"checkpoint"
	"fmt"
	"net/http"
	"time"
)
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	})

    // limits requests by ip user can only make 15 requests during a 30 second window
	rlMiddleware := checkpoint.LimitByEndpoint(15, 30 * time.Second)
	rlHandler := rlMiddleware(mux)

	http.ListenAndServe(":8080", rlHandler)
}
```

I use chi for my go servers and use checkpoint like this.
```go
package main

import (
	"time"
	"github.com/aidenfine/checkpoint"
	"github.com/go-chi/chi/v5"
)

func main(){
    r := chi.NewRouter()
    r.Use(checkpoint.LimitByIp(15, 30 * time.Minute))
    ...
    log.Fatal(http.ListenAndServe(":8080"), r)
}
```