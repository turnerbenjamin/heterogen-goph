package hg_middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/turnerbenjamin/heterogen-go/internal/httpErrors"
)

var Logger Middleware = func(next httpErrors.ReqHandler) httpErrors.ReqHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		start := time.Now()
		defer func() { log.Printf("[%s] %s%s\t%s\n", r.Method, r.Host, r.URL.Path, time.Since(start)) }()
		return next(w, r)
	}
}
