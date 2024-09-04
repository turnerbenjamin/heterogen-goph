package hg_services

import (
	"database/sql"

	"github.com/turnerbenjamin/heterogen-go/internal/httpErrors"
	"github.com/turnerbenjamin/heterogen-go/internal/models"
)

type UsersService interface {
	GetAll() ([]models.User, error)
}

type usersService struct {
	db *sql.DB
}

func NewUsersService(database *sql.DB) UsersService {
	return &usersService{
		db: database,
	}
}

func (userSvc *usersService) GetAll() ([]models.User, error) {
	baseQuery := `
	SELECT id, email_address, password, first_name, last_name, business, permissions FROM users;
	`
	//Select users
	rows, err := userSvc.db.Query(baseQuery)
	if err != nil {
		return nil, httpErrors.ServerFail()
	}

	users := []models.User{}
	for rows.Next() {
		user := models.User{}
		err := rows.Scan(&user.Id, &user.EmailAddress, &user.Password, &user.FirstName, &user.LastName, &user.Business, &user.Permissions)
		if err != nil {
			return nil, httpErrors.ServerFail()
		}
		users = append(users, user)
	}

	return users, nil
}
