package validator

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"log"
	"net/mail"
	"regexp"
	"strings"
)

type ValidationPattern struct {
	RegXStr string
	Message string
}

type ValidationRules struct {
	Required           bool
	MinLength          int
	MaxLength          int
	Pattern            *ValidationPattern
	RequireDigit       bool
	RequireSpecialChar bool
	IsEmail            bool
	IsUKLocation       bool
}

func (vr *ValidationRules) HtmlAttributes() []string {
	attributes := []string{}

	if vr.Required {
		attributes = append(attributes, "required")
	}

	if vr.MinLength > 0 {
		attributes = append(attributes, fmt.Sprintf("minlength=\"%d\"", vr.MinLength))
	}

	if vr.MaxLength > 0 {
		attributes = append(attributes, fmt.Sprintf("maxLength=\"%d\"", vr.MaxLength))
	}

	if vr.Pattern != nil {
		attributes = append(attributes, fmt.Sprintf("pattern=\"%s\"", vr.Pattern.RegXStr))
		msgAttribute := fmt.Sprintf("oninvalid=this.setCustomValidity(\"%s\") oninput=\"setCustomValidity('')\"", vr.Pattern.Message)
		msgAttribute = strings.ReplaceAll(msgAttribute, " ", "&nbsp;")
		attributes = append(attributes, msgAttribute)
	}

	return attributes
}

type ValidatedString string

func (s ValidatedString) Validate(fieldName string, vr *ValidationRules) (bool, string) {
	if vr.Required && len(s) == 0 {
		return false, fmt.Sprintf("%s is required", fieldName)
	}

	if len(s) != 0 && vr.MinLength > 0 && len(s) < vr.MinLength {
		return false, fmt.Sprintf("%s must be at least %d characters long", fieldName, vr.MinLength)
	}

	if len(s) != 0 && vr.MaxLength > 0 && len(s) > vr.MaxLength {
		return false, fmt.Sprintf("%s must not exceed %d characters", fieldName, vr.MaxLength)
	}

	if len(s) != 0 && vr.Pattern != nil && !regexp.MustCompile(vr.Pattern.RegXStr).Match([]byte(s)) {
		return false, vr.Pattern.Message
	}

	if len(s) != 0 && vr.IsEmail {
		_, err := mail.ParseAddress(string(s))
		if err != nil {
			return false, "Invalid email address"
		}
	}

	return true, ""
}

type ValidatedLocation struct {
	Lt float64
	Ln float64
}

func (l ValidatedLocation) Validate(fieldName string, vr *ValidationRules) (bool, string) {
	log.Println(l)
	if vr.Required && (l.Ln == 0.0 || l.Lt == 0.0) {
		return false, fmt.Sprintf("%s is required", fieldName)
	}

	if vr.IsUKLocation && (l.Ln < -14.0 || l.Ln > 4.0) {
		return false, "Invalid Longitude"
	}

	if vr.IsUKLocation && (l.Lt < 49 || l.Lt > 63) {
		return false, "Invalid Latitude"
	}

	return true, ""
}

func (l ValidatedLocation) Value() (driver.Value, error) {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "POINT(%f %f)", l.Lt, l.Ln)
	return buf.Bytes(), nil
}
