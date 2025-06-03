package main

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/betasve/go-commerce/services/auth/internal/jsonlog"
	"github.com/stretchr/testify/assert"
)

func TestLogError(t *testing.T) {
	rec := bytes.NewBufferString("")
	app := application{
		logger: jsonlog.New(rec, jsonlog.LevelInfo),
	}
	req := httptest.NewRequest(
		http.MethodGet,
		"/test/url",
		nil,
	)

	app.logError(req, errors.New("test error"))

	assert.Contains(t, rec.String(), "test error")
}

func TestEditConflictResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	app := application{}
	req := httptest.NewRequest(
		http.MethodPatch,
		"/test/url",
		nil,
	)

	app.editConflictResponse(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Result().StatusCode)
	assert.Equal(t, `{"error":"unable to update the record due to an edit conflict, please try again"}`, strings.TrimSpace(rr.Body.String()))
}

func TestErrorResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	app := application{}
	req := httptest.NewRequest(
		http.MethodGet,
		"/test/url",
		nil,
	)

	app.errorResponse(rr, req, http.StatusOK, "some error occured")

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	assert.Equal(t, `{"error":"some error occured"}`, strings.TrimSpace(rr.Body.String()))
}

func TestServerErrorResponse(t *testing.T) {
	rec := bytes.NewBufferString("")
	rr := httptest.NewRecorder()
	app := application{
		logger: jsonlog.New(rec, jsonlog.LevelInfo),
	}
	req := httptest.NewRequest(
		http.MethodGet,
		"/test/url",
		nil,
	)

	app.serverErrorResponse(rr, req, errors.New("new test error"))

	assert.Equal(t, http.StatusInternalServerError, rr.Result().StatusCode)
	assert.Equal(
		t,
		`{"error":"the server encountered a problem and could not process your request"}`,
		strings.TrimSpace(rr.Body.String()),
	)
	assert.Contains(t, rec.String(), "new test error")
}

func TestNotFoundResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	app := application{}
	req := httptest.NewRequest(
		http.MethodGet,
		"/test/url",
		nil,
	)

	app.notFoundResponse(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Result().StatusCode)
	assert.Equal(
		t,
		`{"error":"the requested resource could not be found"}`,
		strings.TrimSpace(rr.Body.String()),
	)
}

func TestMethodNotAllowedResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	app := application{}
	req := httptest.NewRequest(
		http.MethodGet,
		"/test/url",
		nil,
	)

	app.methodNotAllowedResponse(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Result().StatusCode)
	assert.Equal(
		t,
		`{"error":"the GET method is not supported for this resource"}`,
		strings.TrimSpace(rr.Body.String()),
	)
}

func TestBadRequestResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	app := application{}
	req := httptest.NewRequest(
		http.MethodGet,
		"/test/url",
		nil,
	)

	app.badRequestResponse(rr, req, errors.New("bad request error"))

	assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
	assert.Equal(
		t,
		`{"error":"bad request error"}`,
		strings.TrimSpace(rr.Body.String()),
	)
}

func TestFailedValidationRequest(t *testing.T) {
	rr := httptest.NewRecorder()
	app := application{}
	req := httptest.NewRequest(
		http.MethodGet,
		"/test/url",
		nil,
	)

	app.failedValidationResponse(rr, req, map[string]string{"email": "duplicates", "name": "too short"})

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Result().StatusCode)
	assert.Equal(
		t,
		`{"error":{"email":"duplicates","name":"too short"}}`,
		strings.TrimSpace(rr.Body.String()),
	)
}

func TestRateLimitExceededResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	app := application{}
	req := httptest.NewRequest(
		http.MethodGet,
		"/test/url",
		nil,
	)

	app.rateLimitExceededResponse(rr, req)

	assert.Equal(t, http.StatusTooManyRequests, rr.Result().StatusCode)
	assert.Equal(
		t,
		`{"error":"rate limit exceeded"}`,
		strings.TrimSpace(rr.Body.String()),
	)
}
