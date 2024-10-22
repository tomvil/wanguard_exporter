package wgc

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type Response struct {
	Test string `json:"test"`
}

func TestNewClient(t *testing.T) {
	client := NewClient("http://127.0.0.1", "u", "p")

	if client.apiAddress != "http://127.0.0.1" {
		t.Errorf("Expected http://127.0.0.1, got %s", client.apiAddress)
	}

	if client.apiUsername != "u" {
		t.Errorf("Expected user, got %s", client.apiUsername)
	}

	if client.apiPassword != "p" {
		t.Errorf("Expected password, got %s", client.apiPassword)
	}

}

func TestBasicAuth(t *testing.T) {
	auth := basicAuth("u", "p")
	expectedAuth := "dTpw" // Expected base64 encoded "u:p"

	if auth != expectedAuth {
		t.Errorf("Expected %s, got %s", expectedAuth, auth)
	}

}

func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		expectedAuth := "Basic " + basicAuth("u", "p")

		if authHeader != expectedAuth {
			t.Errorf("Expected %s, got %s", expectedAuth, authHeader)
		}

		w.Write([]byte("test"))
	}))
	defer server.Close()

	client := NewClient(server.URL, "u", "p")
	body, err := client.Get("/test")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if string(body) != "test" {
		t.Errorf("Expected test, got %s", string(body))
	}
}

func TestGetParsed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"test": "success"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "u", "p")

	var response Response
	err := client.GetParsed("/test", &response)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if response.Test != "success" {
		t.Errorf("Expected success, got %s", response.Test)
	}
}

func BenchmarkGetParsed(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"test": "success"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "u", "p")

	for i := 0; i < b.N; i++ {
		var response Response
		client.GetParsed("/test", &response)
	}
}
