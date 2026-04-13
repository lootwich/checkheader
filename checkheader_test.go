package checkheader_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lootwich/checkheader"
)

func TestApproved(t *testing.T) {
	cfg := checkheader.CreateConfig()
	cfg.GroupHeaderName = "X-Demo"
	cfg.NeededGroups = []string{"test", "example"}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := checkheader.New(ctx, next, cfg, "checkheader")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Set headers which should work
	req.Header.Set("X-Demo", "test,example, useless, example2")

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
}

func TestCustomSeparator(t *testing.T) {
	cfg := checkheader.CreateConfig()
	cfg.GroupHeaderName = "X-Demo"
	cfg.NeededGroups = []string{"test", "example"}
	cfg.GroupSeparator = "|"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := checkheader.New(ctx, next, cfg, "checkheader")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Demo", "test|example|useless")

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
}

func TestMissingGroups(t *testing.T) {
	cfg := checkheader.CreateConfig()
	cfg.GroupHeaderName = "X-Demo"
	cfg.NeededGroups = []string{"test", "example"}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := checkheader.New(ctx, next, cfg, "checkheader")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Set headers which should work
	req.Header.Set("X-Demo", "example, useless, example2")

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, recorder.Code)
	}
}

func TestMissingHeader(t *testing.T) {
	cfg := checkheader.CreateConfig()
	cfg.GroupHeaderName = "X-Demo"
	cfg.NeededGroups = []string{"test", "example"}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := checkheader.New(ctx, next, cfg, "checkheader")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, recorder.Code)
	}
}
