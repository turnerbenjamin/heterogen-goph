package hg_services

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/turnerbenjamin/heterogen-go/internal/db_models"
)

type HgAuthService interface {
	Create(db_models.User) (*db_models.User, string)
}

type authService struct {
	db *sql.DB
}

func NewAuthService(database *sql.DB) HgAuthService {
	return &authService{
		db: database,
	}
}

/*
REGISTER A NEW USER
*/
func (authSvc *authService) Create(usr db_models.User) (*db_models.User, string) {
	baseQuery := `
	INSERT INTO users
	(id, email_address, first_name, last_name, business, password, permissions)
	VALUES($1,$2,$3,$4,$5,$6,$7)
	RETURNING id, email_address, first_name, last_name, business, permissions
	;
	`

	fmt.Println(usr)
	var newUser db_models.User
	rows, err := authSvc.db.Query(baseQuery, usr.Id, usr.EmailAddress, usr.FirstName, usr.LastName, usr.Business, usr.Password, usr.Permissions)
	if err != nil {
		msg := "Server error"
		if strings.Contains(err.Error(), "users_email_address_key") {
			msg = "Email address is already associated with an account"
		}
		return &newUser, msg
	}

	if rows.Next() {
		err := rows.Scan(&newUser.Id, &newUser.EmailAddress, &newUser.FirstName, &newUser.LastName, &newUser.Business, &newUser.Permissions)
		if err != nil {
			msg := "Server error"
			fmt.Println(err)
			return &newUser, msg
		}
	}
	return &newUser, ""
}
