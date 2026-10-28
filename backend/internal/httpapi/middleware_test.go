package httpapi

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPanicBecomesJSON500AndIsLogged(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	boom := http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("boom") })

	rec := httptest.NewRecorder()
	logRequests(logger, recoverPanics(logger, boom)).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/add", nil))

	if rec.Code != http.StatusInternalServerError || !strings.Contains(rec.Body.String(), `"code":"INTERNAL_ERROR"`) {
		t.Errorf("got %d %s, want 500 INTERNAL_ERROR", rec.Code, rec.Body)
	}
	if !strings.Contains(logs.String(), `"msg":"panic"`) || !strings.Contains(logs.String(), `"status":500`) {
		t.Errorf("logs missing panic or access entry with status 500:\n%s", logs.String())
	}
}

func TestUnexpectedErrorIsInternal(t *testing.T) {
	status, code, _ := mapError(errors.New("something unexpected"))
	if status != http.StatusInternalServerError || code != "INTERNAL_ERROR" {
		t.Errorf("mapError = %d %s, want 500 INTERNAL_ERROR", status, code)
	}
}
