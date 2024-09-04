package models

import (
	"github.com/turnerbenjamin/heterogen-go/internal/httpErrors"
)

type Reports struct {
	Users UserTableData
}

type ResponseModel struct {
	IsProduction bool
	Location     string
	IsLoggedIn   bool
	Errors       []httpErrors.ErrorMessage
	ToastMessage string
	User         *User
	Validators   map[string][]string
	Reports
}
