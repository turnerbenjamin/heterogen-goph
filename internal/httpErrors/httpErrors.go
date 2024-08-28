package httpErrors

import (
	"net/http"

	"github.com/turnerbenjamin/heterogen-go/internal/render"
)

type StatusCode int
type ErrorMessage string
type errorMessages []ErrorMessage

type HttpError struct {
	statusCode StatusCode
	msgs       errorMessages
}

type ReqHandler func(w http.ResponseWriter, r *http.Request) error

func HandleErrors(reqHandler ReqHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := reqHandler(w, r)
		if err == nil {
			return
		}

		httpError, ok := err.(HttpError)
		if !ok {
			httpError = ServerFail()
		}

		w.WriteHeader(400)
		model := map[string][]ErrorMessage{"Errors": httpError.msgs}

		if err = render.Template(w, r, "errorMessage.component.go.tmpl", model); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
	}
}

func Make(statusCode StatusCode, errorMessages []ErrorMessage) HttpError {
	return HttpError{
		statusCode: statusCode,
		msgs:       errorMessages,
	}
}

// 401
func Unauthorised() HttpError {
	return HttpError{
		statusCode: http.StatusUnauthorized,
		msgs:       errorMessages{"Incorrect credentials"},
	}
}

// 422
func InvalidFormSubmission(errorMessages []ErrorMessage) HttpError {
	return HttpError{
		statusCode: http.StatusUnprocessableEntity,
		msgs:       errorMessages,
	}
}

// 500
func ServerFail() HttpError {
	return HttpError{
		statusCode: http.StatusInternalServerError,
		msgs:       errorMessages{"Server Error"},
	}
}

func (e HttpError) Error() string {
	msg := ""
	for _, em := range e.msgs {
		msg += (string(em) + "\n")
	}
	return msg + "\n"
}
