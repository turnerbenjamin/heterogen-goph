package hg_services

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/turnerbenjamin/heterogen-go/internal/httpErrors"
	"github.com/turnerbenjamin/heterogen-go/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type HgAuthService interface {
	Create(models.User) (*models.User, error)
	SignIn(string, string) (*models.User, error)
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
func (authSvc *authService) Create(usr models.User) (*models.User, error) {
	baseQuery := `
	INSERT INTO users
	(id, email_address, first_name, last_name, business, password, permissions)
	VALUES($1,$2,$3,$4,$5,$6,$7)
	RETURNING id, email_address, first_name, last_name, business, permissions
	;
	`
	var newUser models.User
	rows, err := authSvc.db.Query(baseQuery, usr.Id, usr.EmailAddress, usr.FirstName, usr.LastName, usr.Business, usr.HashedPassword, usr.Permissions)

	if err != nil {
		if strings.Contains(err.Error(), "users_email_address_key") {
			return nil, httpErrors.Make(401, []httpErrors.ErrorMessage{"Email address is already associated with an account"})
		}
		return nil, httpErrors.ServerFail()
	}

	if rows.Next() {
		err := rows.Scan(&newUser.Id, &newUser.EmailAddress, &newUser.FirstName, &newUser.LastName, &newUser.Business, &newUser.Permissions)
		if err != nil {
			return nil, httpErrors.ServerFail()
		}
	}
	return &newUser, nil
}

func (authSvc *authService) SignIn(emailAddress string, password string) (*models.User, error) {
	baseQuery := `
	SELECT id, email_address, password, first_name, last_name, business, permissions FROM users
	WHERE email_address=$1
	;
	`

	//Select user
	var user models.User
	row := authSvc.db.QueryRow(baseQuery, emailAddress)

	//Scan row to user struct
	err := row.Scan(&user.Id, &user.EmailAddress, &user.Password, &user.FirstName, &user.LastName, &user.Business, &user.Permissions)
	if err != nil {
		return nil, httpErrors.Unauthorised()
	}

	//Validate password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	fmt.Println(user.Password, password)
	if err != nil {
		return nil, httpErrors.Unauthorised()
	}

	user.Password = ""

	return &user, nil
}
