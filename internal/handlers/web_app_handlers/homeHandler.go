package web_app_handlers

import (
	"net/http"

	"github.com/turnerbenjamin/heterogen-go/internal/render"
)

var HomeHandler = func(w http.ResponseWriter, r *http.Request) error {
	return render.Template(w, r, "home.page.go.tmpl", nil)
}
