package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlers(t *testing.T) {
	app := application{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}

	ts := httptest.NewServer(app.routes())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/ping")
	if err != nil {
		t.Errorf("error getting response from / with error: %v", err)
	}

	data, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Errorf("error getting response from / with error: %v", err)
	}

	if !strings.Contains(string(data), "pong") {
		t.Errorf("unable to to find wanted data in / route")
	}
}
