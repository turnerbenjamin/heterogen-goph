package models

import (
	"log"
	"net/http"
	"strconv"

	"github.com/turnerbenjamin/heterogen-go/internal/httpErrors"
	"github.com/turnerbenjamin/heterogen-go/internal/uuid"
	"github.com/turnerbenjamin/heterogen-go/internal/validator"
)

type Business struct {
	Id           string
	Reference    string
	TradingName  validator.ValidatedString
	EmailAddress string
	Address      []validator.ValidatedString
	Postcode     validator.ValidatedString
	Location     validator.ValidatedLocation
	IsGrower     bool
	CphNumber    validator.ValidatedString
}

var BusinessValidationRules = map[string]*validator.ValidationRules{
	"TradingName": {
		Required:  true,
		MinLength: 3,
		MaxLength: 64,
	},
	"EmailAddress": {
		Required:  true,
		MinLength: 3,
		MaxLength: 64,
	},
	"AddressLine1": {
		Required:  true,
		MinLength: 3,
		MaxLength: 128,
	},
	"AddressLine2": {
		Required:  false,
		MinLength: 3,
		MaxLength: 128,
	},
	"AddressLine3": {
		Required:  false,
		MinLength: 3,
		MaxLength: 128,
	},
	"Postcode": {
		Required:  true,
		MinLength: 5,
		MaxLength: 7,
	},
	"Location": {
		Required:     true,
		IsUKLocation: true,
	},
	"CphNumber": {
		Required: false,
		Pattern: &validator.ValidationPattern{
			RegXStr: `^\d{2}\/\d{3}\/\d{4}$`,
			Message: "Invalid CPH number, please use format XX/XXX/XXXX",
		},
	},
}

func BusinessFromForm(r *http.Request) (*Business, error) {
	var business Business

	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	//Parse Address
	l1 := validator.ValidatedString(r.PostFormValue("address-line-1"))
	l2 := validator.ValidatedString(r.PostFormValue("address-line-2"))
	l3 := validator.ValidatedString(r.PostFormValue("address-line-3"))
	address := []validator.ValidatedString{l1, l2, l3}

	//Parse Location
	latitude, err := strconv.ParseFloat(r.PostFormValue("latitude"), 64)
	if err != nil {
		return nil, httpErrors.Make(http.StatusUnprocessableEntity, []httpErrors.ErrorMessage{"Invalid location"})
	}
	longitude, err := strconv.ParseFloat(r.PostFormValue("longitude"), 64)
	if err != nil {
		return nil, httpErrors.Make(http.StatusUnprocessableEntity, []httpErrors.ErrorMessage{"Invalid location"})
	}
	Location := validator.ValidatedLocation{Lt: latitude, Ln: longitude}

	//Generate Reference
	postcode := validator.ValidatedString(r.PostFormValue("postcode"))
	ok, errorMessage := postcode.Validate("Postcode", BusinessValidationRules["Postcode"])
	if !ok {
		return nil, httpErrors.Make(http.StatusUnprocessableEntity, []httpErrors.ErrorMessage{httpErrors.ErrorMessage(errorMessage)})
	}
	reference := "test"

	//Generate ID
	id, err := uuid.Generate()
	if err != nil {
		return nil, err
	}

	//Construct business struct
	business = Business{
		Id:          id,
		Reference:   reference,
		TradingName: validator.ValidatedString(r.PostFormValue("trading_name")),
		Postcode:    postcode,
		IsGrower:    r.PostFormValue("is_grower") == "true",
		Address:     address,
		Location:    Location,
		CphNumber:   validator.ValidatedString(r.PostFormValue("cph_number")),
	}

	log.Println(business)

	return &business, err
}

func (b *Business) Validate() (bool, []httpErrors.ErrorMessage) {

	errorMessages := []httpErrors.ErrorMessage{}

	ok, err := b.TradingName.Validate("Trading name", BusinessValidationRules["TradingName"])
	if !ok {
		errorMessages = append(errorMessages, httpErrors.ErrorMessage(err))
	}

	if len(b.Address) > 0 {
		ok, err = b.Address[0].Validate("Address line 1", BusinessValidationRules["AddressLine1"])
		if !ok {
			errorMessages = append(errorMessages, httpErrors.ErrorMessage(err))
		}
	}

	if len(b.Address) > 1 {
		ok, err = b.Address[1].Validate("Address line 2", BusinessValidationRules["AddressLine2"])
		if !ok {
			errorMessages = append(errorMessages, httpErrors.ErrorMessage(err))
		}
	}

	if len(b.Address) > 2 {
		ok, err = b.Address[2].Validate("Address line 3", BusinessValidationRules["AddressLine3"])
		if !ok {
			errorMessages = append(errorMessages, httpErrors.ErrorMessage(err))
		}
	}

	ok, err = b.Postcode.Validate("Postcode", BusinessValidationRules["Postcode"])
	if !ok {
		errorMessages = append(errorMessages, httpErrors.ErrorMessage(err))
	}

	ok, err = b.Location.Validate("Location", BusinessValidationRules["Location"])
	if !ok {
		errorMessages = append(errorMessages, httpErrors.ErrorMessage(err))
	}

	ok, err = b.CphNumber.Validate("CPH Number", BusinessValidationRules["CphNumber"])
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
