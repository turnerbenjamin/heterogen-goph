package httpErrors

import (
	"net/http"
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
	msg := ""
	for _, em := range e.Msgs {
		msg += (string(em) + "\n")
	}
	return msg + "\n"
}
