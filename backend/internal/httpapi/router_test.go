package httpapi_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The routing cases run with and without a static directory because
// registering the "/" file server changes how ServeMux resolves /api paths.
func TestRouting(t *testing.T) {
	staticDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(staticDir, "index.html"), []byte("<h1>calculator</h1>"), 0o644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		method     string
		path       string
		wantStatus int
		wantCode   string
	}{
		{method: http.MethodGet, path: "/api/v1/add", wantStatus: 405, wantCode: "METHOD_NOT_ALLOWED"},
		{method: http.MethodGet, path: "/api/v1/unknown", wantStatus: 404, wantCode: "UNSUPPORTED_OPERATION"},
		{method: http.MethodPost, path: "/api", wantStatus: 404, wantCode: "NOT_FOUND"},
		{method: http.MethodPost, path: "/api/nothing", wantStatus: 404, wantCode: "NOT_FOUND"},
		{method: http.MethodPost, path: "/api/v1/add/extra", wantStatus: 404, wantCode: "NOT_FOUND"},
		{method: http.MethodPost, path: "/api/v2/add", wantStatus: 404, wantCode: "NOT_FOUND"},
	}

	for _, mode := range []struct {
		name      string
		staticDir string
	}{
		{name: "api only", staticDir: ""},
		{name: "with static files", staticDir: staticDir},
	} {
		t.Run(mode.name, func(t *testing.T) {
			h := newTestRouter(t, mode.staticDir)
			for _, tt := range tests {
				rec, resp := do(t, h, tt.method, tt.path, jsonContentType, "")
				if rec.Code != tt.wantStatus || resp.Error == nil || resp.Error.Code != tt.wantCode {
					t.Errorf("%s %s = %d %s, want %d %s", tt.method, tt.path, rec.Code, rec.Body, tt.wantStatus, tt.wantCode)
				}
			}

			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
			if rec.Code != http.StatusOK || strings.TrimSpace(rec.Body.String()) != `{"status":"ok"}` {
				t.Errorf("GET /healthz = %d %s", rec.Code, rec.Body)
			}
		})
	}

	t.Run("serves the frontend build", func(t *testing.T) {
		rec := httptest.NewRecorder()
		newTestRouter(t, staticDir).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
		body, _ := io.ReadAll(rec.Body)
		if rec.Code != http.StatusOK || !strings.Contains(string(body), "calculator") {
			t.Errorf("GET / = %d %s", rec.Code, body)
		}
	})

	t.Run("does not list directories", func(t *testing.T) {
		assets := filepath.Join(staticDir, "assets")
		if err := os.Mkdir(assets, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(assets, "app.js"), []byte("console.log(1)"), 0o644); err != nil {
			t.Fatal(err)
		}
		h := newTestRouter(t, staticDir)

		for path, want := range map[string]int{"/assets/": http.StatusNotFound, "/assets/app.js": http.StatusOK} {
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
			if rec.Code != want {
				t.Errorf("GET %s = %d, want %d", path, rec.Code, want)
			}
		}
	})
}
