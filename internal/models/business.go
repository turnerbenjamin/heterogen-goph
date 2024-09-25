package models

import (
	"net/http"

	"github.com/turnerbenjamin/heterogen-go/internal/httpErrors"
	"github.com/turnerbenjamin/heterogen-go/internal/uuid"
	"github.com/turnerbenjamin/heterogen-go/internal/validator"
)

type Location struct {
	Lt float64
	Ln float64
}

type Business struct {
	Id           string
	Reference    string
	CPH_Number   string
	TradingName  validator.ValidatedString
	Location     Location
	Postcode     validator.ValidatedString
	IsGrower     bool
	Logo         string
	About        string
	EmailAddress string
	Website      string
}

var BusinessValidationRules = map[string]*validator.ValidationRules{
	"TradingName": {
		Required:  true,
		MinLength: 3,
		MaxLength: 32,
	},
	"Address": {
		Required:  true,
		MinLength: 3,
		MaxLength: 32,
	},
	"Postcode": {
		Required:  true,
		MinLength: 3,
		MaxLength: 32,
	},
}

func BusinessFromForm(r *http.Request) (Business, error) {
	var business Business
	err := r.ParseForm()
	if err != nil {
		return business, err
	}

	business = Business{
		TradingName: validator.ValidatedString(r.PostFormValue("trading_name")),
		Postcode:    validator.ValidatedString(r.PostFormValue("postcode")),
		IsGrower:    r.PostFormValue("is_grower") == "true",
	}
	//!PARSE LOCATION

	business.Id, err = uuid.Generate()

	return business, err
}

func (b *Business) Validate() (bool, []httpErrors.ErrorMessage) {
	errorMessages := []httpErrors.ErrorMessage{}

	ok, err := b.TradingName.Validate("TradingName", UserValidationRules["TradingName"])
	if !ok {
		errorMessages = append(errorMessages, httpErrors.ErrorMessage(err))
	}

	ok, err = b.Postcode.Validate("Postcode", UserValidationRules["Postcode"])
	if !ok {
		errorMessages = append(errorMessages, httpErrors.ErrorMessage(err))
	}

	return len(errorMessages) == 0, errorMessages

}

func BusinessValidationHTMLAttributes() map[string][]string {
	htmlAttributes := map[string][]string{}

	for key, vr := range BusinessValidationRules {
		htmlAttributes[key] = vr.HtmlAttributes()
	}

	return htmlAttributes
}
