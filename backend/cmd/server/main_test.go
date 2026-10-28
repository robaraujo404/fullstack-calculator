package main

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"sezzle-calculator/internal/httpapi"
)

func discardLogger() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }

func TestLoadConfig(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("STATIC_DIR", "")
	if got := loadConfig(); got != (config{port: "8080"}) {
		t.Errorf("defaults = %+v, want port 8080 and no static dir", got)
	}

	t.Setenv("PORT", "9000")
	t.Setenv("STATIC_DIR", "/srv/web")
	if got := loadConfig(); got != (config{port: "9000", staticDir: "/srv/web"}) {
		t.Errorf("from env = %+v", got)
	}
}

func TestRunRejectsMissingStaticDir(t *testing.T) {
	err := run(context.Background(), config{port: "0", staticDir: "/does/not/exist"}, discardLogger())
	if err == nil || !strings.Contains(err.Error(), "STATIC_DIR") {
		t.Fatalf("run() error = %v, want a STATIC_DIR error", err)
	}
}

func TestRunReportsPortInUse(t *testing.T) {
	taken, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	defer taken.Close()
	_, port, _ := net.SplitHostPort(taken.Addr().String())

	var logs bytes.Buffer
	err = run(context.Background(), config{port: port}, slog.New(slog.NewTextHandler(&logs, nil)))
	if err == nil {
		t.Fatal("run() succeeded on a port that is already in use")
	}
	if strings.Contains(logs.String(), "listening") {
		t.Errorf("logged %q before failing to bind", logs.String())
	}
}

func TestServeAndShutdown(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- serve(ctx, ln, httpapi.NewRouter(discardLogger(), ""), discardLogger()) }()

	resp, err := http.Get("http://" + ln.Addr().String() + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET /healthz = %d, want 200", resp.StatusCode)
	}

	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("serve() after cancel = %v, want nil", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("serve() did not return after the context was cancelled")
	}
}
