package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/betasve/go-commerce/services/auth/internal/data"
	"github.com/julienschmidt/httprouter"
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
			http.MethodPost,
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

func TestShowUserHandler(t *testing.T) {
	tests := []struct {
		name                 string
		userId               string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{"Error on invalid id", "abc", http.StatusNotFound, `{"error":"the requested resource could not be found"}`},
		{"Error on missing user id", "0", http.StatusNotFound, `{"error":"the requested resource could not be found"}`},
		{"Finds the user", "1", http.StatusOK, `{"user":{"id":42,"name":"John Doe","email":"test_email@example.com","password":"[FILTERED]","created_at":"2025-03-26T15:04:05Z","updated_at":"2025-03-26T15:04:05Z"}}`},
	}

	for _, tc := range tests {
		rr := httptest.NewRecorder()

		app := application{
			models: data.NewMockModels(),
		}

		req := httptest.NewRequest(
			http.MethodGet,
			"/test/url",
			bytes.NewReader([]byte("")),
		)

		params := httprouter.Params{httprouter.Param{Key: "id", Value: tc.userId}}
		ctx := context.WithValue(req.Context(), httprouter.ParamsKey, params)
		req = req.WithContext(ctx)

		app.showUserHandler(rr, req)

		assert.Equal(t, tc.expectedStatusCode, rr.Result().StatusCode)
		assert.Equal(
			t,
			tc.expectedResponseBody,
			strings.TrimSpace(rr.Body.String()),
		)
	}
}

func TestUpdateUserHandler(t *testing.T) {
	tests := []struct {
		name                 string
		userId               string
		userBody             string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{"Error on invalid id", "abc", `{"name": "Johny Do","email":"test@example.com","password":"NewPass123"}`, http.StatusNotFound, `{"error":"the requested resource could not be found"}`},
		{"Error on missing user id", "0", `{"name": "Johny Do","email":"test@example.com","password":"NewPass123"}`, http.StatusNotFound, `{"error":"the requested resource could not be found"}`},
		{"Updates the user", "1", `{"name": "Johny Do","email":"test@example.com","password":"NewPass123"}`, http.StatusOK, `{"user":{"id":42,"name":"Johny Do","email":"test@example.com","password":"[FILTERED]","created_at":"2025-03-26T15:04:05Z","updated_at":"2025-03-26T15:04:05Z"}}`},
		{"Partially updates the user", "1", `{"name": "Johny Do"}`, http.StatusOK, `{"user":{"id":42,"name":"Johny Do","email":"test_email@example.com","password":"[FILTERED]","created_at":"2025-03-26T15:04:05Z","updated_at":"2025-03-26T15:04:05Z"}}`},
	}

	for _, tc := range tests {
		rr := httptest.NewRecorder()

		app := application{
			models: data.NewMockModels(),
		}

		req := httptest.NewRequest(
			http.MethodPatch,
			"/test/url",
			bytes.NewReader([]byte(tc.userBody)),
		)

		params := httprouter.Params{httprouter.Param{Key: "id", Value: tc.userId}}
		ctx := context.WithValue(req.Context(), httprouter.ParamsKey, params)
		req = req.WithContext(ctx)

		app.updateUserHandler(rr, req)

		assert.Equal(t, tc.expectedStatusCode, rr.Result().StatusCode)
		assert.Equal(
			t,
			tc.expectedResponseBody,
			strings.TrimSpace(rr.Body.String()),
		)
	}
}

func TestDeleteUserHandler(t *testing.T) {
	tests := []struct {
		name                 string
		userId               string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{"Error on invalid id", "abc", http.StatusNotFound, `{"error":"the requested resource could not be found"}`},
		{"Error on missing user id", "0", http.StatusNotFound, `{"error":"the requested resource could not be found"}`},
		{"Updates the user", "1", http.StatusNoContent, `{}`},
	}

	for _, tc := range tests {
		rr := httptest.NewRecorder()

		app := application{
			models: data.NewMockModels(),
		}

		req := httptest.NewRequest(
			http.MethodPut,
			"/test/url",
			bytes.NewReader([]byte("")),
		)

		params := httprouter.Params{httprouter.Param{Key: "id", Value: tc.userId}}
		ctx := context.WithValue(req.Context(), httprouter.ParamsKey, params)
		req = req.WithContext(ctx)

		app.deleteUserHandler(rr, req)

		assert.Equal(t, tc.expectedStatusCode, rr.Result().StatusCode)
		assert.Equal(
			t,
			tc.expectedResponseBody,
			strings.TrimSpace(rr.Body.String()),
		)
	}
}
