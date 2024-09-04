package htmlHandler

import (
	"log"
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
func (ac *AuthHandler) RegistrationPage(w http.ResponseWriter, r *http.Request, m *models.ResponseModel) error {
	m.Validators = models.UserValidationHTMLAttributes()
	render.Page(w, r, "registration", m, 200)
	return nil
}

// Creates a new user.
func (ac *AuthHandler) Register(w http.ResponseWriter, r *http.Request, m *models.ResponseModel) error {
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
	w.Header().Add("Hx-Redirect", "/")
	return nil
}

// Set auth cookie and refresh the page
func (ac *AuthHandler) LogIn(w http.ResponseWriter, r *http.Request, m *models.ResponseModel) error {

	err := r.ParseForm()
	if err != nil {
		return httpErrors.ServerFail()
	}

	emailAddress := r.PostFormValue("email_address")
	password := r.PostFormValue("password")

	user, err := ac.service.SignIn(emailAddress, password)
	if err != nil {
		log.Println(err)
		return err
	}

	http.SetCookie(w, cookies.NewAuthCookie(user.Id, cookies.Day*7))
	w.Header().Add("Hx-Redirect", "/dashboard")
	return nil
}

// *Log-Out
// Returns the log in modal
func (ac *AuthHandler) LogOut(w http.ResponseWriter, r *http.Request, m *models.ResponseModel) error {
	http.SetCookie(w, cookies.UnsetAuthCookie())
	w.Header().Add("Hx-Redirect", "/")
	return nil
	// return httpErrors.ServerFail()
}
