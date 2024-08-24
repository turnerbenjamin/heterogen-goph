package hg_middleware

import (
	"log"
	"net/http"
	"time"
)

var Logger Middleware = func(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() { log.Printf("[%s] %s%s\t%s\n", r.Method, r.Host, r.URL.Path, time.Since(start)) }()
		next(w, r)
	}
}
