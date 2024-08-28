package hg_middleware

import (
	"fmt"
	"net/http"

	"github.com/turnerbenjamin/heterogen-go/internal/cookies"
	"github.com/turnerbenjamin/heterogen-go/internal/httpErrors"
)

var PrintUserId Middleware = func(next httpErrors.ReqHandler) httpErrors.ReqHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		userId, ok := cookies.ParseAuthCookie(r)
		if !ok {
			next(w, r)
			return nil
		}

		fmt.Println(userId)
		return next(w, r)
	}
}
