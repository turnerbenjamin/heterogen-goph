package web_app_handlers

import (
	"net/http"

	"github.com/turnerbenjamin/heterogen-go/internal/render"
)

var HomeHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "home.page.tmpl", nil)
})
