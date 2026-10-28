package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"sezzle-calculator/internal/calculator"
)

type errorResponse struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// requestError is a client error detected before the operation runs. It
// already carries its HTTP status and error code.
type requestError struct {
	status  int
	code    string
	message string
}

func (e *requestError) Error() string { return e.message }

func invalidRequest(message string) *requestError {
	return &requestError{status: http.StatusBadRequest, code: "INVALID_REQUEST", message: message}
}

// mapError translates an error into the HTTP status, code and message sent to
// the client. 400 means the request is malformed; 422 means it is well formed
// but mathematically invalid.
func mapError(err error) (status int, code, message string) {
	var reqErr *requestError
	var tooLarge *http.MaxBytesError
	switch {
	case errors.As(err, &reqErr):
		return reqErr.status, reqErr.code, reqErr.message
	case errors.As(err, &tooLarge):
		return http.StatusRequestEntityTooLarge, "PAYLOAD_TOO_LARGE",
			fmt.Sprintf("request body must not exceed %d bytes", tooLarge.Limit)
	case errors.Is(err, calculator.ErrDivisionByZero):
		return http.StatusUnprocessableEntity, "DIVISION_BY_ZERO", err.Error()
	case errors.Is(err, calculator.ErrInvalidOperand):
		return http.StatusUnprocessableEntity, "INVALID_OPERAND", err.Error()
	case errors.Is(err, calculator.ErrResultOutOfRange):
		return http.StatusUnprocessableEntity, "RESULT_OUT_OF_RANGE", err.Error()
	default:
		return http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error"
	}
}

func (a *api) writeError(w http.ResponseWriter, err error) {
	status, code, message := mapError(err)
	if status == http.StatusInternalServerError {
		a.logger.Error("request failed", "error", err)
	}
	writeJSON(w, status, errorResponse{Error: errorDetail{Code: code, Message: message}})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	// Values are finite numbers and strings, so encoding cannot fail; the only
	// possible error is a failed write to a client that went away.
	_ = json.NewEncoder(w).Encode(v)
}
