package htmlHandler

import (
	"net/http"

	"github.com/turnerbenjamin/heterogen-go/internal/models"
	"github.com/turnerbenjamin/heterogen-go/internal/render"
)

var DashboardHandler = func(w http.ResponseWriter, r *http.Request, m *models.ResponseModel) error {
	return render.Page(w, r, "dashboard", m, http.StatusOK)
}
