package htmlHandler

import (
	"net/http"

	"github.com/turnerbenjamin/heterogen-go/internal/hg_services"
	"github.com/turnerbenjamin/heterogen-go/internal/models"
	"github.com/turnerbenjamin/heterogen-go/internal/render"
)

type UsersHandler struct {
	service hg_services.UsersService
}

func NewUsersHandler(s hg_services.UsersService) *UsersHandler {
	return &UsersHandler{
		service: s,
	}
}

//*REGISTER

func (uh *UsersHandler) UsersPage(w http.ResponseWriter, r *http.Request, m *models.ResponseModel) error {
	render.Page(w, r, "users", m, 200)
	return nil
}

func (uh *UsersHandler) UsersTable(w http.ResponseWriter, r *http.Request, m *models.ResponseModel) error {
	users, err := uh.service.GetAll()
	if err != nil {
		return err
	}

	m.Reports.Users = models.GetUserTableData(users, models.NameUC, models.BusinessUS, models.AdminUC, models.EmailUC)
	render.Component(w, r, "usersTable", m, 200)
	return nil
}
