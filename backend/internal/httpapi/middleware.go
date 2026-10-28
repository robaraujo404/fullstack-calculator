package httpapi

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

// logRequests is the outermost middleware, so it also records the 500s
// produced by recoverPanics.
func logRequests(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		logger.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration", time.Since(start).String(),
		)
	})
}

func recoverPanics(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			v := recover()
			if v == nil {
				return
			}
			if v == http.ErrAbortHandler {
				panic(v) // deliberate abort; let net/http handle it
			}
			logger.Error("panic", "error", v, "stack", string(debug.Stack()))
			writeJSON(w, http.StatusInternalServerError, errorResponse{
				Error: errorDetail{Code: "INTERNAL_ERROR", Message: "internal server error"},
			})
		}()
		next.ServeHTTP(w, r)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
