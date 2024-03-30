package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHomeHandler(t *testing.T) {
	t.Skip("Skipping test for now")
	// Switch to test mode so you don't get such noisy output
	gin.SetMode(gin.TestMode)

	router := NewHttpServer(HTTPServerConfig{
		Host: "localhost",
		Port: 8080,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected to receive status %d, but received %d", http.StatusOK, w.Code)
	}

	// Check the response body is what we expect.
	expected := "Expected response body here"
	if w.Body.String() != expected {
		t.Fatalf("Expected to receive body %s, but received %s", expected, w.Body.String())
	}
}
