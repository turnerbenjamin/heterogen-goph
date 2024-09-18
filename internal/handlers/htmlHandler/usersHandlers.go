package htmlHandler

import (
	"net/http"
	"strings"

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

	//column config
	columnQuery := r.URL.Query().Get("columns")
	var columns []string
	if columnQuery != "" {
		columns = strings.Split(columnQuery, ",")
	}
	columnConfig, err := models.GetColumnConfig(columns)
	if err != nil {
		return err
	}

	//sorting config
	var tableSortConfig *models.TableSortConfig
	if sortingQuery := r.URL.Query().Get("sort"); sortingQuery != "" {
		tableSortConfig = columnConfig.ApplySortingQuery(sortingQuery)
	}

	//Make Query
	users, err := uh.service.GetAll(tableSortConfig)
	if err != nil {
		return err
	}

	m.Reports.TableData = models.GetUserTableData(users, columnConfig)
	render.Component(w, r, "usersTable", m, 200)
	return nil
}
