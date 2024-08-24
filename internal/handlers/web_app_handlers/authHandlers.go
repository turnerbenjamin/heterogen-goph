package web_app_handlers

import (
	"fmt"
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

	ok, validationMessages := user.Validate()
	if !ok {
		w.WriteHeader(400)
		model := map[string][]string{"Errors": validationMessages}
		render.Template(w, "errorMessage.component.tmpl", model)
		return
	}

	_, errorMessage := ac.service.Create(user)

	if errorMessage != "" {
		w.WriteHeader(400)
		model := map[string][]string{"Errors": {errorMessage}}
		render.Template(w, "errorMessage.component.tmpl", model)
		return
	}

	render.Template(w, "successMessage.component.tmpl", "<p>Account created successfully. <a href=\"log-in\">Log-in</a></p>")

}

func (ac *AuthController) RegistrationPage(w http.ResponseWriter, r *http.Request) {

	validators := db_models.UserValidationHTMLAttributes()
	render.Template(w, "registration.page.tmpl", validators)
}

func (ac *AuthController) LogInPage(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "logInForm.component.tmpl", nil)
}

func (ac *AuthController) LogIn(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(500)
		model := map[string][]string{"Errors": {"Server error"}}
		render.Template(w, "errorMessage.component.tmpl", model)
		return
	}

	emailAddress := r.PostFormValue("email_address")
	password := r.PostFormValue("password")

	_, errMessage := ac.service.SignIn(emailAddress, password)
	if errMessage != "" {
		fmt.Println(errMessage)
		w.WriteHeader(401)
		model := map[string][]string{"Errors": {errMessage}}
		render.Template(w, "errorMessage.component.tmpl", model)
		return
	}
	w.Header().Set("HX-Refresh", "true")

	// render.Template(w, "logInForm.component.tmpl", nil)
}
