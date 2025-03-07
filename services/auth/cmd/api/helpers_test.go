package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type UserReq struct {
	testStruct
}

func TestReadIDParam(t *testing.T) {
	tests := []struct {
		name      string
		urlParam  string
		expectID  int64
		expectErr error
	}{
		{"Valid ID", "5", 5, nil},
		{"Invalid ID (negative)", "-1", 0, errors.New("invalid ID parameter")},
		{"Invalid ID (non-numeric)", "abc", 0, errors.New("invalid ID parameter")},
		{"Empty ID", "", 0, errors.New("invalid ID parameter")},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test/"+tc.urlParam, nil)

			params := httprouter.Params{httprouter.Param{Key: "id", Value: tc.urlParam}}
			ctx := context.WithValue(req.Context(), httprouter.ParamsKey, params)
			req = req.WithContext(ctx)

			app := application{}
			id, err := app.readIDParam(req)

			if id != tc.expectID {
				t.Errorf("Expected ID %d, got %d", tc.expectID, id)
			}

			if (err == nil && tc.expectErr != nil) || (err != nil && tc.expectErr == nil) || (err != nil && err.Error() != tc.expectErr.Error()) {
				t.Errorf("Expected error '%v', got '%v'", tc.expectErr, err)
			}
		})
	}
}

func TestWriteJSON(t *testing.T) {
	user := testStruct{
		Email: "test@example.com",
		Name:  "john doe",
	}

	rr := httptest.NewRecorder()

	app := application{}
	app.writeJSON(
		rr,
		200,
		envelope{"user": user},
		make(http.Header),
	)

	assert.Equal(t, 200, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.Equal(
		t,
		`{"user":{"email":"test@example.com","name":"john doe"}}`,
		strings.TrimSpace(rr.Body.String()),
	)
}

func TestReadJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	tests := []struct {
		name       string
		reqBody    string
		expectBody string
		expectErr  error
	}{
		{"Valid Request Body", `{"user":{"email":"test@example.com","name":"john doe"}}`, strings.TrimSpace(rr.Body.String()), nil},
		{"Invalid JSON", `{"user": "test", name = john}`, "{}", errors.New("the body contains badly-formed JSON (at character 18)")},
		{"Invalid JSON type field", `{"notuser":{}}`, "{}", errors.New("the body contains unknown key \"notuser\"")},
		{"Invalid JSON empty body", ``, "{}", errors.New("the body must not be empty")},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			outputUser := UserReq{}
			req := httptest.NewRequest(
				http.MethodGet,
				"/test/url",
				bytes.NewReader([]byte(tc.reqBody)),
			)

			app := application{}
			err := app.readJSON(rr, req, &outputUser)
			if tc.expectErr != nil {
				assert.Equal(t, tc.expectErr, err)
			}

			if err == nil {
				jsn, err := json.Marshal(&outputUser)
				assert.Equal(t, tc.expectErr, err)
				assert.Equal(t, strings.TrimSpace(string(jsn)), tc.expectBody)
			}
		})
	}
}
