package httpapi_test

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"sezzle-calculator/internal/httpapi"
)

const jsonContentType = "application/json"

type apiResponse struct {
	Operation string   `json:"operation"`
	Result    *float64 `json:"result"`
	Error     *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func newTestRouter(t *testing.T, staticDir string) http.Handler {
	t.Helper()
	return httpapi.NewRouter(slog.New(slog.NewTextHandler(io.Discard, nil)), staticDir)
}

func do(t *testing.T, h http.Handler, method, path, contentType, body string) (*httptest.ResponseRecorder, apiResponse) {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	var resp apiResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("%s %s: response is not JSON: %v\n%s", method, path, err, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Errorf("%s %s: Content-Type = %q, want JSON", method, path, got)
	}
	return rec, resp
}

func TestCalculateSuccess(t *testing.T) {
	tests := []struct {
		op   string
		body string
		want float64
	}{
		{op: "add", body: `{"a": 2, "b": 3}`, want: 5},
		{op: "subtract", body: `{"a": 4, "b": 10}`, want: -6},
		{op: "multiply", body: `{"a": -2, "b": 0.5}`, want: -1},
		{op: "divide", body: `{"a": 10, "b": 4}`, want: 2.5},
		{op: "power", body: `{"a": 2, "b": 10}`, want: 1024},
		{op: "sqrt", body: `{"a": 16}`, want: 4},
		{op: "sqrt", body: `{"a": 16, "b": null}`, want: 4},
		{op: "percentage", body: `{"a": 15, "b": 200}`, want: 30},
		{op: "add", body: `{"a": 1e3, "b": -0.5}`, want: 999.5},
		{op: "add", body: `{"a": 0, "b": 0}`, want: 0},
		// Values below the smallest float64 decode to 0 (documented assumption).
		{op: "add", body: `{"a": 1e-400, "b": 1}`, want: 1},
	}

	h := newTestRouter(t, "")
	for _, tt := range tests {
		t.Run(tt.op+" "+tt.body, func(t *testing.T) {
			rec, resp := do(t, h, http.MethodPost, "/api/v1/"+tt.op, jsonContentType, tt.body)
			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body %s", rec.Code, rec.Body)
			}
			if resp.Operation != tt.op || resp.Result == nil || *resp.Result != tt.want {
				t.Errorf("body = %s, want operation %q and result %v", rec.Body, tt.op, tt.want)
			}
		})
	}
}

func TestCalculateErrors(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		path        string
		contentType string
		body        string
		wantStatus  int
		wantCode    string
		wantMessage string // substring; empty means not checked
	}{
		// 400: malformed request
		{name: "empty body", body: ``, wantStatus: 400, wantCode: "INVALID_REQUEST", wantMessage: "empty"},
		{name: "truncated JSON", body: `{"a": 1,`, wantStatus: 400, wantCode: "INVALID_REQUEST", wantMessage: "valid JSON"},
		{name: "invalid JSON", body: `{a: 1}`, wantStatus: 400, wantCode: "INVALID_REQUEST", wantMessage: "valid JSON"},
		{name: "array body", body: `[1, 2]`, wantStatus: 400, wantCode: "INVALID_REQUEST", wantMessage: "JSON object"},
		{name: "string body", body: `"1 + 2"`, wantStatus: 400, wantCode: "INVALID_REQUEST", wantMessage: "JSON object"},
		{name: "numeric string", body: `{"a": "10", "b": 1}`, wantStatus: 400, wantCode: "INVALID_REQUEST", wantMessage: `"a"`},
		{name: "number out of float64 range", body: `{"a": 1e400, "b": 1}`, wantStatus: 400, wantCode: "INVALID_REQUEST", wantMessage: "float64 range"},
		{name: "unknown field", body: `{"a": 1, "b": 2, "c": 3}`, wantStatus: 400, wantCode: "INVALID_REQUEST", wantMessage: `"c"`},
		{name: "concatenated objects", body: `{"a": 1, "b": 2}{"a": 3}`, wantStatus: 400, wantCode: "INVALID_REQUEST", wantMessage: "single JSON object"},
		{name: "null body", body: `null`, wantStatus: 400, wantCode: "MISSING_OPERAND"},
		{name: "missing a", body: `{"b": 2}`, wantStatus: 400, wantCode: "MISSING_OPERAND", wantMessage: `"a"`},
		{name: "null a", body: `{"a": null, "b": 2}`, wantStatus: 400, wantCode: "MISSING_OPERAND", wantMessage: `"a"`},
		{name: "missing b", body: `{"a": 1}`, wantStatus: 400, wantCode: "MISSING_OPERAND", wantMessage: `"b"`},
		{name: "b sent to sqrt", path: "/api/v1/sqrt", body: `{"a": 4, "b": 2}`, wantStatus: 400, wantCode: "UNEXPECTED_OPERAND"},

		// 404, 405, 413, 415
		{name: "unsupported operation", path: "/api/v1/modulo", body: `{"a": 1, "b": 2}`, wantStatus: 404, wantCode: "UNSUPPORTED_OPERATION"},
		{name: "operation names are case-sensitive", path: "/api/v1/ADD", body: `{"a": 1, "b": 2}`, wantStatus: 404, wantCode: "UNSUPPORTED_OPERATION"},
		{name: "unknown operation with GET", method: http.MethodGet, path: "/api/v1/modulo", wantStatus: 404, wantCode: "UNSUPPORTED_OPERATION"},
		{name: "GET", method: http.MethodGet, wantStatus: 405, wantCode: "METHOD_NOT_ALLOWED"},
		{name: "PUT", method: http.MethodPut, body: `{"a": 1, "b": 2}`, wantStatus: 405, wantCode: "METHOD_NOT_ALLOWED"},
		{name: "body over 1 KiB", body: `{"a": 1, "b": 2}` + strings.Repeat("x", 2048), wantStatus: 413, wantCode: "PAYLOAD_TOO_LARGE"},
		{name: "missing Content-Type", contentType: "-", body: `{"a": 1, "b": 2}`, wantStatus: 415, wantCode: "UNSUPPORTED_MEDIA_TYPE", wantMessage: "Content-Type: application/json"},
		{name: "form Content-Type (curl -d default)", contentType: "application/x-www-form-urlencoded", body: `{"a": 1, "b": 2}`, wantStatus: 415, wantCode: "UNSUPPORTED_MEDIA_TYPE"},

		// 422: well-formed but mathematically invalid
		{name: "division by zero", path: "/api/v1/divide", body: `{"a": 1, "b": 0}`, wantStatus: 422, wantCode: "DIVISION_BY_ZERO"},
		{name: "division by negative zero", path: "/api/v1/divide", body: `{"a": 1, "b": -0}`, wantStatus: 422, wantCode: "DIVISION_BY_ZERO"},
		{name: "sqrt of negative", path: "/api/v1/sqrt", body: `{"a": -1}`, wantStatus: 422, wantCode: "INVALID_OPERAND"},
		{name: "fractional power of negative", path: "/api/v1/power", body: `{"a": -8, "b": 0.5}`, wantStatus: 422, wantCode: "INVALID_OPERAND"},
		{name: "overflow", path: "/api/v1/multiply", body: `{"a": 1e308, "b": 10}`, wantStatus: 422, wantCode: "RESULT_OUT_OF_RANGE"},
	}

	h := newTestRouter(t, "")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method, path, contentType := tt.method, tt.path, tt.contentType
			if method == "" {
				method = http.MethodPost
			}
			if path == "" {
				path = "/api/v1/add"
			}
			switch contentType {
			case "":
				contentType = jsonContentType
			case "-":
				contentType = ""
			}

			rec, resp := do(t, h, method, path, contentType, tt.body)

			if rec.Code != tt.wantStatus || resp.Error == nil || resp.Error.Code != tt.wantCode {
				t.Fatalf("got %d %s, want %d %s", rec.Code, rec.Body, tt.wantStatus, tt.wantCode)
			}
			if !strings.Contains(resp.Error.Message, tt.wantMessage) {
				t.Errorf("message = %q, want it to mention %q", resp.Error.Message, tt.wantMessage)
			}
			if tt.wantStatus == http.StatusMethodNotAllowed && rec.Header().Get("Allow") != http.MethodPost {
				t.Errorf("Allow = %q, want POST", rec.Header().Get("Allow"))
			}
		})
	}
}

func TestContentTypeWithCharset(t *testing.T) {
	rec, _ := do(t, newTestRouter(t, ""), http.MethodPost, "/api/v1/add",
		"application/json; charset=utf-8", `{"a": 1, "b": 2}`)
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
}
