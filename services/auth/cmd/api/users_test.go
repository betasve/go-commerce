package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
	}

	for _, tc := range tests {
		rr := httptest.NewRecorder()
		app := application{}
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
