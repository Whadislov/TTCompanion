package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsApiReady(t *testing.T) {
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(IsApiReadyHandler)
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestLoginHandler(t *testing.T) {
	creds := map[string]string{
		"username": "testuser",
		"password": "testpassword",
	}
	body, _ := json.Marshal(creds)

	req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(LoginHandler)
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusUnauthorized {
		t.Errorf("expected unauthorized, got %v", status)
	}
}

func TestSignUpHandler(t *testing.T) {
	signUpData := map[string]string{
		"username": "newuser",
		"password": "newpassword",
		"email":    "newuser@example.com",
	}
	body, _ := json.Marshal(signUpData)

	req, err := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(SignUpHandler)
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusInternalServerError {
		t.Errorf("expected error due to missing DB, got %v", status)
	}
}

func TestCorsMiddleware(t *testing.T) {
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := CorsMiddleware(mockHandler)
	req, _ := http.NewRequest("OPTIONS", "/api", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected 200, got %v", recorder.Code)
	}
}
