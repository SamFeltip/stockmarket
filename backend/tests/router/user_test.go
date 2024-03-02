package router

import (
	"net/http"
	"net/http/httptest"
	"stockmarket/database"
	"stockmarket/router"
	"testing"
)

func TestUserCardEndpoint(t *testing.T) {
	database.SetupDevDb()
	r := router.SetupRoutes()

	req, _ := http.NewRequest(http.MethodGet, "/users/card/1", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	// assert.Equal(t, http.StatusOK, resp.Code)
	// assert.Equal(t, expectedBody, resp.Body.String())
}

func TestUsersEndpoint(t *testing.T) {
	database.SetupDevDb()
	r := router.SetupRoutes()

	req, _ := http.NewRequest(http.MethodGet, "/users/show/1", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	// assert.Equal(t, http.StatusOK, resp.Code)
	// assert.Equal(t, expectedBody, resp.Body.String())
}
