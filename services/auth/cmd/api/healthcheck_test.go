package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthcheckHandler(t *testing.T) {
	rr := httptest.NewRecorder()
	app := application{
		config: config{
			env: "test",
		},
	}
	req := httptest.NewRequest(
		http.MethodGet,
		"/test/url",
		nil,
	)

	app.healthcheckHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	assert.Equal(
		t,
		`{"status":"available","system_info":{"environment":"test","version":"1.0.0"}}`,
		strings.TrimSpace(rr.Body.String()),
	)
}
