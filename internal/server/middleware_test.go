package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetSecurity_Headers(t *testing.T) {
	s := newTestServer(t, &fakeRepo{})
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	handler := s.setSecurity(inner)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	want := map[string]string{
		"X-Frame-Options":        "deny",
		"X-Content-Type-Options": "nosniff",
		"X-XSS-Protection":       "1; mode=block",
		"Referrer-Policy":        "strict-origin-when-cross-origin",
		"Cache-Control":          "no-cache, no-store, must-revalidate",
	}
	for header, wantVal := range want {
		if got := w.Header().Get(header); got != wantVal {
			t.Errorf("%s: got %q, want %q", header, got, wantVal)
		}
	}
}

func TestSetSecurity_PassesThrough(t *testing.T) {
	s := newTestServer(t, &fakeRepo{})
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	handler := s.setSecurity(inner)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTeapot {
		t.Errorf("got %d, want 418", w.Code)
	}
}

func TestSetCommon_ContentType(t *testing.T) {
	s := newTestServer(t, &fakeRepo{})
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	handler := s.setCommon(inner)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if ct := w.Header().Get("Content-Type"); ct != "text/html" {
		t.Errorf("Content-Type: got %q, want %q", ct, "text/html")
	}
}

func TestLogHandler_PassesThrough(t *testing.T) {
	s := newTestServer(t, &fakeRepo{})
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := s.logHandler("test-handler", inner)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want 200", w.Code)
	}
}
