package htmlHandler

import (
	"net/http"

	"github.com/turnerbenjamin/heterogen-go/internal/hg_services"
	"github.com/turnerbenjamin/heterogen-go/internal/httpErrors"
	"github.com/turnerbenjamin/heterogen-go/internal/models"
	"github.com/turnerbenjamin/heterogen-go/internal/render"
)

type BusinessHandler struct {
	service hg_services.BusinessService
}

func NewBusinessesHandler(s hg_services.BusinessService) *BusinessHandler {
	return &BusinessHandler{
		service: s,
	}
}

//*REGISTER

func (bh *BusinessHandler) AddBusinessPage(w http.ResponseWriter, r *http.Request, m *models.ResponseModel) error {
	m.Validators = models.BusinessValidationHTMLAttributes()
	render.Page(w, r, "addBusiness", m, 200)
	return nil
}

func (bh *BusinessHandler) CreateBusiness(w http.ResponseWriter, r *http.Request, m *models.ResponseModel) error {
	business, err := models.BusinessFromForm(r)
	if err != nil {
		return err
	}

	ok, validationMessages := business.Validate()
	if !ok {
		return httpErrors.InvalidFormSubmission(validationMessages)
	}

	_, err = bh.service.Create(*business)
	if err != nil {
		return err
	}
	// w.Header().Add("Hx-Redirect", "/")
	return nil
}
