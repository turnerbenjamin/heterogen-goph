package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/turnerbenjamin/heterogen-go/internal/models"
	"github.com/turnerbenjamin/heterogen-go/internal/router"
)

var Logger router.Middleware = func(next router.ReqHandler) router.ReqHandler {
	return func(w http.ResponseWriter, r *http.Request, m *models.ResponseModal) error {
		start := time.Now()
		defer func() { log.Printf("[%s] %s%s\t%s\n", r.Method, r.Host, r.URL.Path, time.Since(start)) }()
		return next(w, r, m)
	}
}
