package models

import "github.com/turnerbenjamin/heterogen-go/internal/httpErrors"

type ResponseModal struct {
	IsProduction bool
	IsLoggedIn   bool
	UserId       string
	Errors       []httpErrors.ErrorMessage
	ToastMessage string
}
