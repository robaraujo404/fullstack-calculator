package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"strings"

	"sezzle-calculator/internal/calculator"
)

const maxBodyBytes = 1 << 10

type api struct {
	logger *slog.Logger
}

// calculateRequest uses pointers so that a missing or null operand can be told
// apart from zero.
type calculateRequest struct {
	A *float64 `json:"a"`
	B *float64 `json:"b"`
}

type calculateResponse struct {
	Operation string  `json:"operation"`
	Result    float64 `json:"result"`
}

func (a *api) handleCalculate(w http.ResponseWriter, r *http.Request) {
	resp, err := calculate(w, r)
	if err != nil {
		a.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (a *api) handleNotFound(w http.ResponseWriter, r *http.Request) {
	a.writeError(w, &requestError{
		status:  http.StatusNotFound,
		code:    "NOT_FOUND",
		message: fmt.Sprintf("no API route for %s", r.URL.Path),
	})
}

// calculate runs the checks in a fixed order (operation, method, media type,
// body, operands, math) so that each request maps to exactly one error.
func calculate(w http.ResponseWriter, r *http.Request) (calculateResponse, error) {
	name := r.PathValue("operation")
	op, ok := calculator.Lookup(name)
	if !ok {
		return calculateResponse{}, &requestError{
			status:  http.StatusNotFound,
			code:    "UNSUPPORTED_OPERATION",
			message: fmt.Sprintf("unsupported operation %q", name),
		}
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		return calculateResponse{}, &requestError{
			status:  http.StatusMethodNotAllowed,
			code:    "METHOD_NOT_ALLOWED",
			message: fmt.Sprintf("method %s is not allowed; use POST", r.Method),
		}
	}

	if mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type")); err != nil || mediaType != "application/json" {
		return calculateResponse{}, &requestError{
			status:  http.StatusUnsupportedMediaType,
			code:    "UNSUPPORTED_MEDIA_TYPE",
			message: "send Content-Type: application/json",
		}
	}

	// Reading the whole body first makes the size limit independent of its
	// content: anything over maxBodyBytes is a 413, even if it starts with a
	// valid JSON object.
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxBodyBytes))
	if err != nil {
		return calculateResponse{}, err
	}

	req, err := decodeRequest(body)
	if err != nil {
		return calculateResponse{}, err
	}

	operands, err := operandsFor(op, req)
	if err != nil {
		return calculateResponse{}, err
	}

	result, err := op.Apply(operands...)
	if err != nil {
		return calculateResponse{}, err
	}
	return calculateResponse{Operation: op.Name, Result: result}, nil
}

// decodeRequest follows encoding/json in matching field names
// case-insensitively and keeping the last of duplicate keys; the README
// documents both.
func decodeRequest(body []byte) (calculateRequest, error) {
	// Decoding into a pointer exposes a null body, which would otherwise leave
	// a struct untouched and surface later as a missing operand.
	var req *calculateRequest
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		return calculateRequest{}, invalidRequest(decodeErrorMessage(err))
	}
	if req == nil {
		return calculateRequest{}, invalidRequest("request body must be a JSON object")
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return calculateRequest{}, invalidRequest("request body must contain a single JSON object")
	}
	return *req, nil
}

func decodeErrorMessage(err error) string {
	var typeErr *json.UnmarshalTypeError
	var syntaxErr *json.SyntaxError
	switch {
	case errors.Is(err, io.EOF):
		return "request body is empty"
	case errors.As(err, &typeErr) && typeErr.Field == "":
		return "request body must be a JSON object"
	case errors.As(err, &typeErr) && strings.HasPrefix(typeErr.Value, "number"):
		// encoding/json reports a number that overflows float64 ("a": 1e400)
		// as a type error whose Value is "number 1e400".
		return fmt.Sprintf("field %q is outside the float64 range", typeErr.Field)
	case errors.As(err, &typeErr):
		return fmt.Sprintf("field %q must be a number", typeErr.Field)
	case errors.As(err, &syntaxErr), errors.Is(err, io.ErrUnexpectedEOF):
		return "request body must be valid JSON"
	default:
		// For example `json: unknown field "c"`, which has no dedicated type.
		return strings.TrimPrefix(err.Error(), "json: ")
	}
}

func operandsFor(op calculator.Operation, req calculateRequest) ([]float64, error) {
	if req.A == nil {
		return nil, &requestError{
			status:  http.StatusBadRequest,
			code:    "MISSING_OPERAND",
			message: `operand "a" is required`,
		}
	}

	if op.Arity == 1 {
		if req.B != nil {
			return nil, &requestError{
				status:  http.StatusBadRequest,
				code:    "UNEXPECTED_OPERAND",
				message: fmt.Sprintf(`%s takes only operand "a"`, op.Name),
			}
		}
		return []float64{*req.A}, nil
	}

	if req.B == nil {
		return nil, &requestError{
			status:  http.StatusBadRequest,
			code:    "MISSING_OPERAND",
			message: fmt.Sprintf(`operand "b" is required for %s`, op.Name),
		}
	}
	return []float64{*req.A, *req.B}, nil
}
