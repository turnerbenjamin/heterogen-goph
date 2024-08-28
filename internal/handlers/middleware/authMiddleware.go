package middleware

import (
	"fmt"
	"net/http"

	"github.com/turnerbenjamin/heterogen-go/internal/cookies"
	"github.com/turnerbenjamin/heterogen-go/internal/models"
	"github.com/turnerbenjamin/heterogen-go/internal/router"
)

var PrintUserId router.Middleware = func(next router.ReqHandler) router.ReqHandler {
	return func(w http.ResponseWriter, r *http.Request, m *models.ResponseModal) error {
		userId, ok := cookies.ParseAuthCookie(r)
		if !ok {
			next(w, r, m)
			return nil
		}

		fmt.Println(userId)
		return next(w, r, m)
	}
}
