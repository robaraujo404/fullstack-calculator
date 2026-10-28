// Package httpapi exposes the calculator over HTTP. Every response under /api/
// is JSON, including 404, 405 and 500.
package httpapi

import (
	"log/slog"
	"net/http"
	"strings"
)

// NewRouter wires the API routes and, when staticDir is set, the frontend build.
//
// The registrations below work around two ServeMux behaviours that would break
// the "every /api response is JSON" contract:
//
//   - Method-specific patterns ("POST /api/...") make ServeMux answer 405 in
//     plain text, and only when no other pattern matches. Once the "/" file
//     server is registered, GET /api/v1/add would fall through to it and get a
//     404 instead. So the calculation route accepts any method and the handler
//     checks it.
//   - "/api/" catches every other API path and answers 404 in JSON. "/api" is
//     registered too, or ServeMux would redirect it to "/api/" with an empty
//     307 response.
//
// The file server is registered as "/" without a method on purpose: "GET /"
// would overlap "/api/" with neither pattern more specific, and ServeMux panics
// on conflicting registrations.
func NewRouter(logger *slog.Logger, staticDir string) http.Handler {
	a := &api{logger: logger}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/{operation}", a.handleCalculate)
	mux.HandleFunc("/api/", a.handleNotFound)
	mux.HandleFunc("/api", a.handleNotFound)
	mux.HandleFunc("GET /healthz", handleHealth)
	if staticDir != "" {
		mux.Handle("/", noDirectoryListing(http.FileServer(http.Dir(staticDir))))
	}

	return logRequests(logger, recoverPanics(logger, mux))
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// noDirectoryListing answers 404 for directory paths other than the root,
// which http.FileServer would otherwise list.
func noDirectoryListing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
