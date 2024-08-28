package htmlHandler

import (
	"net/http"

	"github.com/turnerbenjamin/heterogen-go/internal/cookies"
	"github.com/turnerbenjamin/heterogen-go/internal/hg_services"
	"github.com/turnerbenjamin/heterogen-go/internal/httpErrors"
	"github.com/turnerbenjamin/heterogen-go/internal/models"
	"github.com/turnerbenjamin/heterogen-go/internal/render"
)

type AuthHandler struct {
	service hg_services.HgAuthService
}

func NewAuthHandler(s hg_services.HgAuthService) *AuthHandler {
	return &AuthHandler{
		service: s,
	}
}

//*REGISTER

// Returns registration page
func (ac *AuthHandler) RegistrationPage(w http.ResponseWriter, r *http.Request, m *models.ResponseModal) error {
	validators := models.UserValidationHTMLAttributes()
	render.Template(w, r, "registration.page.go.tmpl", validators)
	return nil
}

// Creates a new user. Returns a success modal
func (ac *AuthHandler) Register(w http.ResponseWriter, r *http.Request, m *models.ResponseModal) error {
	user, err := models.UserFromForm(r)
	if err != nil {
		return httpErrors.ServerFail()
	}

	ok, validationMessages := user.Validate()
	if !ok {
		return httpErrors.InvalidFormSubmission(validationMessages)
	}

	_, err = ac.service.Create(user)
	if err != nil {
		return err
	}

	render.Template(w, r, "successMessage.component.go.tmpl", "<p>Account created successfully. <a href=\"log-in\">Log-in</a></p>")
	return nil
}

//*LOG IN

// Returns the log in modal
func (ac *AuthHandler) LogInPage(w http.ResponseWriter, r *http.Request, m *models.ResponseModal) error {
	render.Template(w, r, "logInForm.component.go.tmpl", nil)
	return nil
}

// Set auth cookie and refresh the page
func (ac *AuthHandler) LogIn(w http.ResponseWriter, r *http.Request, m *models.ResponseModal) error {

	err := r.ParseForm()
	if err != nil {
		return httpErrors.ServerFail()
	}

	emailAddress := r.PostFormValue("email_address")
	password := r.PostFormValue("password")

	user, err := ac.service.SignIn(emailAddress, password)
	if err != nil {
		return err
	}

	http.SetCookie(w, cookies.NewAuthCookie(user.Id, cookies.Day*7))
	w.Header().Set("HX-Refresh", "true")
	return nil
}

// *Log-Out
// Returns the log in modal
func (ac *AuthHandler) LogOut(w http.ResponseWriter, r *http.Request, m *models.ResponseModal) error {
	http.SetCookie(w, cookies.UnsetAuthCookie())
	w.Header().Set("HX-Refresh", "true")
	return nil
}
