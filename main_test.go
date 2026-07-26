package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func get(t *testing.T, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	NewRouter().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}

func TestHealth(t *testing.T) {
	rec := get(t, "/health")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("bad json: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf(`expected status "ok", got %q`, body["status"])
	}
}

func TestVersionHasMetadata(t *testing.T) {
	rec := get(t, "/version")
	var body map[string]string
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	if _, ok := body["version"]; !ok {
		t.Fatal("version field missing")
	}
	if _, ok := body["environment"]; !ok {
		t.Fatal("environment field missing")
	}
}

func TestOrders(t *testing.T) {
	rec := get(t, "/orders")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var body struct {
		Count int `json:"count"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	if body.Count != 3 {
		t.Fatalf("expected 3 orders, got %d", body.Count)
	}
}

func TestUnknownRoute(t *testing.T) {
	if rec := get(t, "/definitely-not-a-route"); rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}
