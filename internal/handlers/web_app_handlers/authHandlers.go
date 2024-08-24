package web_app_handlers

import (
	"net/http"

	"github.com/turnerbenjamin/heterogen-go/internal/db_models"
	"github.com/turnerbenjamin/heterogen-go/internal/hg_services"
	"github.com/turnerbenjamin/heterogen-go/internal/render"
)

type AuthController struct {
	service hg_services.HgAuthService
}

func NewAuthController(s hg_services.HgAuthService) *AuthController {
	return &AuthController{
		service: s,
	}
}

func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	user, err := db_models.UserFromForm(r)
	if err != nil {
		println(err)
		return
	}
	_, errorMessage := ac.service.Create(user)
	if errorMessage != "" {
		w.WriteHeader(400)
		model := map[string][]string{"Errors": {errorMessage}}
		render.Template(w, "errorMessage.component.tmpl", model)
	}

	render.Template(w, "successMessage.component.tmpl", "<p>Account created successfully. <a href=\"log-in\">Log-in</a></p>")

}

func (ac *AuthController) RegistrationPage(w http.ResponseWriter, r *http.Request) {

	validators := db_models.UserValidationHTMLAttributes()
	render.Template(w, "registration.page.tmpl", validators)
}

func (ac *AuthController) LogInPage(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "log-in.page.tmpl", nil)
}
