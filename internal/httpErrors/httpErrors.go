package httpErrors

import (
	"net/http"
	"strings"
)

type StatusCode int
type ErrorMessage string
type errorMessages []ErrorMessage

type HttpError struct {
	StatusCode StatusCode
	Msgs       errorMessages
}

func Make(statusCode StatusCode, errorMessages []ErrorMessage) HttpError {
	return HttpError{
		StatusCode: statusCode,
		Msgs:       errorMessages,
	}
}

func Unauthorised() HttpError {
	return HttpError{
		StatusCode: http.StatusUnauthorized, //401
		Msgs:       errorMessages{"Incorrect credentials"},
	}
}

func InvalidFormSubmission(errorMessages []ErrorMessage) HttpError {
	return HttpError{
		StatusCode: http.StatusUnprocessableEntity, //422
		Msgs:       errorMessages,
	}
}

func ServerFail() HttpError {
	return HttpError{
		StatusCode: http.StatusInternalServerError, //500
		Msgs:       errorMessages{"Server Error"},
	}
}

func (e HttpError) Error() string {
	lines := []string{}
	for _, em := range e.Msgs {
		lines = append(lines, string(em))
	}
	return strings.Join(lines, "\n")
}
