package db_models

import (
	"net/http"

	"github.com/lib/pq"
	"github.com/turnerbenjamin/heterogen-go/internal/uuid"
	"github.com/turnerbenjamin/heterogen-go/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id             string
	EmailAddress   validator.ValidatedString
	FirstName      validator.ValidatedString
	LastName       validator.ValidatedString
	Business       validator.ValidatedString
	Password       validator.ValidatedString
	HashedPassword string
	Permissions    pq.StringArray
}

var UserValidationRules = map[string]*validator.ValidationRules{
	"EmailAddress": {
		Required: true,
		IsEmail:  true,
	},
	"FirstName": {
		Required:  true,
		MinLength: 3,
		MaxLength: 32,
	},
	"LastName": {
		Required:  true,
		MinLength: 3,
		MaxLength: 32,
	},
	"Business": {
		Required:  true,
		MinLength: 3,
		MaxLength: 32,
	},
	"Password": {
		Required:           true,
		MinLength:          8,
		MaxLength:          32,
		RequireDigit:       true,
		RequireSpecialChar: true,
		Pattern: &validator.ValidationPattern{
			RegXStr: "^.*[!£$%^&*#~].*$",
			Message: "Password must contain at least one digit and special character (!£$%^&*#~)",
		},
	},
}

func UserFromForm(r *http.Request) (User, error) {
	var user User
	err := r.ParseForm()
	if err != nil {
		return user, err
	}

	user = User{
		EmailAddress: validator.ValidatedString(r.PostFormValue("email_address")),
		FirstName:    validator.ValidatedString(r.PostFormValue("first_name")),
		LastName:     validator.ValidatedString(r.PostFormValue("last_name")),
		Password:     validator.ValidatedString(r.PostFormValue("password")),
		Business:     validator.ValidatedString(r.PostFormValue("business")),
		Permissions:  pq.StringArray{},
	}

	err = user.hashPassword()
	if err != nil {
		return user, err
	}

	user.Id, err = uuid.Generate()

	return user, err
}

func (u *User) Validate() (bool, []string) {
	errorMessages := []string{}

	ok, err := u.EmailAddress.Validate("EmailAddress", UserValidationRules["EmailAddress"])
	if !ok {
		errorMessages = append(errorMessages, err)
	}

	ok, err = u.FirstName.Validate("FirstName", UserValidationRules["FirstName"])
	if !ok {
		errorMessages = append(errorMessages, err)
	}

	ok, err = u.LastName.Validate("LastName", UserValidationRules["LastName"])
	if !ok {
		errorMessages = append(errorMessages, err)
	}

	ok, err = u.Business.Validate("Business", UserValidationRules["Business"])
	if !ok {
		errorMessages = append(errorMessages, err)
	}

	ok, err = u.Password.Validate("Password", UserValidationRules["Password"])
	if !ok {
		errorMessages = append(errorMessages, err)
	}

	return len(errorMessages) == 0, errorMessages

}

func UserValidationHTMLAttributes() map[string][]string {
	htmlAttributes := map[string][]string{}

	for key, vr := range UserValidationRules {
		htmlAttributes[key] = vr.HtmlAttributes()
	}

	return htmlAttributes
}

func (u *User) hashPassword() error {
	passwordBS := []byte(u.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBS, 12)
	if err != nil {
		return err
	}
	u.HashedPassword = string(hashedPassword)
	return nil
}
