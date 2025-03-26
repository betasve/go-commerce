package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/betasve/go-commerce/services/auth/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserHandler(t *testing.T) {
	tests := []struct {
		name                 string
		reqBody              string
		expectedStatusCode   int
		expectedLocation     string
		expectedResponseBody string
	}{
		{"Error on empty request body", "", http.StatusBadRequest, "", `{"error":"the body must not be empty"}`},
		{"Error on empty email in user body", `{"email":"","password": "pass123","name":"Test User"}`, http.StatusUnprocessableEntity, "", `{"error":{"email":"can't be blank"}}`},
		{"Error on empty password in user body", `{"email":"test_email@example.com","password": "","name":"Test User"}`, http.StatusUnprocessableEntity, "", `{"error":{"password":"can't be blank"}}`},
		{"Error on empty name in user body", `{"email":"test_email@example.com","password": "pass123","name":""}`, http.StatusUnprocessableEntity, "", `{"error":{"name":"can't be blank"}}`},
		{"Successfully created the user", `{"email":"test_email@example.com","password": "pass123","name":"John Doe"}`, http.StatusCreated, "/v1/users/42", `{"user":{"id":42,"name":"John Doe","email":"test_email@example.com","password":"[FILTERED]","created_at":"2025-03-26T15:04:05Z","updated_at":"2025-03-26T15:04:05Z"}}`},
	}

	for _, tc := range tests {
		rr := httptest.NewRecorder()

		app := application{
			models: data.NewMockModels(),
		}

		req := httptest.NewRequest(
			http.MethodGet,
			"/test/url",
			bytes.NewReader([]byte(tc.reqBody)),
		)

		app.createUserHandler(rr, req)

		assert.Equal(t, tc.expectedStatusCode, rr.Result().StatusCode)
		assert.Equal(t, tc.expectedLocation, rr.Header().Get("Location"))
		assert.Equal(
			t,
			tc.expectedResponseBody,
			strings.TrimSpace(rr.Body.String()),
		)
	}
}
