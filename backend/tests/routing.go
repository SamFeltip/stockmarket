package main

import (
	"net/http"
	"net/http/httptest"
	"stockmarket/database"
	"stockmarket/router"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelloEndpoint(t *testing.T) {
	router := router.SetupRoutes() // Assuming SetupRouter initializes your router with all the routes

	req, _ := http.NewRequest(http.MethodGet, "/hello", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "Hello, World!\n", resp.Body.String())
}

func TestUserCardEndpoint(t *testing.T) {
	database.SetupDb()
	r := router.SetupRoutes()

	req, _ := http.NewRequest(http.MethodGet, "/users/card/1", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	// assert.Equal(t, http.StatusOK, resp.Code)
	// assert.Equal(t, expectedBody, resp.Body.String())
}

func TestUsersEndpoint(t *testing.T) {
	database.SetupDb()
	r := router.SetupRoutes()

	req, _ := http.NewRequest(http.MethodGet, "/users/show/1", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	// assert.Equal(t, http.StatusOK, resp.Code)
	// assert.Equal(t, expectedBody, resp.Body.String())
}
