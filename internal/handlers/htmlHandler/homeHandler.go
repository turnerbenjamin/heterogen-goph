package htmlHandler

import (
	"net/http"

	"github.com/turnerbenjamin/heterogen-go/internal/models"
	"github.com/turnerbenjamin/heterogen-go/internal/render"
)

var HomeHandler = func(w http.ResponseWriter, r *http.Request, m *models.ResponseModel) error {
	if r.URL.Path != "/" {
		return render.Page(w, r, "notFound", m, http.StatusOK)
	}
	return render.Page(w, r, "home", m, http.StatusOK)
}
