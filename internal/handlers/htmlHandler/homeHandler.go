package htmlHandler

import (
	"net/http"

	"github.com/turnerbenjamin/heterogen-go/internal/models"
	"github.com/turnerbenjamin/heterogen-go/internal/render"
)

var HomeHandler = func(w http.ResponseWriter, r *http.Request, m *models.ResponseModal) error {
	return render.Template(w, r, "home.page.go.tmpl", nil)
}
