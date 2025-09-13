package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var (
	// allow localhost (any port, http/https)
	localhostRe = regexp.MustCompile(`^https?://localhost(:\d+)?$`)
	// allow keepwondering.io and any subdomain (http/https)
	kwRe = regexp.MustCompile(`^https?://([a-z0-9-]+\.)*keepwondering\.io(:\d+)?$`)
)

// allowedOrigin returns the origin if it is allowed, otherwise empty string.
func allowedOrigin(origin string) string {
	origin = strings.TrimSpace(origin)
	if origin == "" {
		return ""
	}
	if localhostRe.MatchString(origin) || kwRe.MatchString(origin) {
		return origin
	}
	return ""
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowed := allowedOrigin(origin)

		// Always vary on these so caches don't poison responses across origins.
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")
		w.Header().Add("Vary", "Access-Control-Request-Headers")

		if allowed != "" {
			w.Header().Set("Access-Control-Allow-Origin", allowed)
			// Set to true only if you actually use cookies/auth headers across origins.
			// If not needed, you can omit this.
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Handle preflight
		if r.Method == http.MethodOptions {
			if allowed != "" {
				reqMethod := r.Header.Get("Access-Control-Request-Method")
				reqHeaders := r.Header.Get("Access-Control-Request-Headers")

				if reqMethod != "" {
					w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
				}
				// Echo requested headers (safer with credentials than using '*')
				if reqHeaders != "" {
					w.Header().Set("Access-Control-Allow-Headers", reqHeaders)
				} else {
					w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, X-Requested-With")
				}
				// cache preflight for an hour
				w.Header().Set("Access-Control-Max-Age", "3600")
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Normal request flow
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, `{ "message": "hello world!" }`)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("listening on %s", port)
	log.Fatal(http.ListenAndServe(":"+port, cors(mux)))
}
