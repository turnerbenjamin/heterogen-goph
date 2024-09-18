package hg_services

import (
	"database/sql"
	"fmt"

	"github.com/turnerbenjamin/heterogen-go/internal/httpErrors"
	"github.com/turnerbenjamin/heterogen-go/internal/models"
)

type UsersService interface {
	GetAll(*models.TableSortConfig) ([]models.User, error)
}

type usersService struct {
	db *sql.DB
}

func NewUsersService(database *sql.DB) UsersService {
	return &usersService{
		db: database,
	}
}

func (userSvc *usersService) GetAll(tableSortConfig *models.TableSortConfig) ([]models.User, error) {
	fmt.Println(tableSortConfig)
	baseQuery := `
	SELECT id, email_address, password, first_name, last_name, business, permissions 
	FROM (
		SELECT id, email_address, password, first_name, last_name, business, permissions, is_admin, CONCAT(last_name, ' ', first_name) AS full_name 
		FROM USERS
	) As fullTable
	`

	if tableSortConfig != nil {
		baseQuery += fmt.Sprintf("ORDER BY %s %s", tableSortConfig.Fieldname, tableSortConfig.Direction)
	}

	//Select users
	rows, err := userSvc.db.Query(baseQuery)
	if err != nil {
		fmt.Println(err)
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
