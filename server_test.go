package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
)

var (
	expectedJSON = `{"alive":true}`
)

func TestHealthCheckHandler(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/health-check")
	e.NewContext(req, rec)

	HealthCheckHandler(c)

	if rec.Body.String() != expectedJSON {
		t.Errorf("response does not match: got %v want %v", rec.Code, http.StatusOK)
	}
}
