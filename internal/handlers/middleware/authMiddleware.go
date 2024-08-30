package middleware

import (
	"net/http"

	"github.com/turnerbenjamin/heterogen-go/internal/cookies"
	"github.com/turnerbenjamin/heterogen-go/internal/models"
	"github.com/turnerbenjamin/heterogen-go/internal/router"
)

var ParseAuthCookie router.Middleware = func(next router.ReqHandler) router.ReqHandler {
	return func(w http.ResponseWriter, r *http.Request, m *models.ResponseModal) error {
		if userId, ok := cookies.ParseAuthCookie(r); ok {
			m.UserId = userId
			m.IsLoggedIn = true
		}
		return next(w, r, m)
	}
}
