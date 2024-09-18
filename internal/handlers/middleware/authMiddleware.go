package middleware

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/turnerbenjamin/heterogen-go/internal/cookies"
	"github.com/turnerbenjamin/heterogen-go/internal/hg_services"
	"github.com/turnerbenjamin/heterogen-go/internal/models"
	"github.com/turnerbenjamin/heterogen-go/internal/render"
	"github.com/turnerbenjamin/heterogen-go/internal/router"
)

func GetUserAuthenticator(authService hg_services.HgAuthService) router.Middleware {
	return func(next router.ReqHandler) router.ReqHandler {
		return func(w http.ResponseWriter, r *http.Request, m *models.ResponseModel) error {
			fmt.Println("GET USER")
			m.Location = r.URL.Path
			if m.Location == "" {
				m.Location = "/"
			}
			userId, ok := cookies.ParseAuthCookie(r)
			if !ok {
				return next(w, r, m)
			}
			user, err := authService.GetById(userId)

			if err != nil {
				return next(w, r, m)
			}
			m.IsLoggedIn = true
			m.User = user
			return next(w, r, m)
		}
	}
}

func RequireAuthentication() router.Middleware {
	return func(next router.ReqHandler) router.ReqHandler {
		return func(w http.ResponseWriter, r *http.Request, m *models.ResponseModel) error {
			fmt.Println("AUTH HANDLER")
			fmt.Println("CHECK USER: ", m.User)
			if m.User == nil {
				return render.Page(w, r, "notAuthenticated", m, http.StatusAccepted)
			}
			return next(w, r, m)
		}
	}
}

func RequireAdmin() router.Middleware {
	return func(next router.ReqHandler) router.ReqHandler {
		return func(w http.ResponseWriter, r *http.Request, m *models.ResponseModel) error {
			fmt.Println("CHECK ADMIN: ", m.User)
			if m.User == nil || !slices.Contains(m.User.Permissions, "admin") {
				return render.Page(w, r, "forbidden", m, http.StatusAccepted)
			}
			return next(w, r, m)
		}
	}
}
