package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelloEndpoint(t *testing.T) {
	router := SetupRouter() // Assuming SetupRouter initializes your router with all the routes

	req, _ := http.NewRequest(http.MethodGet, "/hello", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "Hello, World!\n", resp.Body.String())
}

func TestUserCardEndpoint(t *testing.T) {
	router := SetupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/user-card/1", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// assert.Equal(t, http.StatusOK, resp.Code)
	// assert.Equal(t, expectedBody, resp.Body.String())
}

func TestUsersEndpoint(t *testing.T) {
	router := SetupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// assert.Equal(t, http.StatusOK, resp.Code)
	// assert.Equal(t, expectedBody, resp.Body.String())
}
