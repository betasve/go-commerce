package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/betasve/go-commerce/services/auth/internal/data"
	"github.com/betasve/go-commerce/services/auth/internal/validator"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)

		return
	}

	v := validator.New()

	v.Check(input.Name != "", "name", "can't be blank")
	v.Check(len(input.Name) > 5, "name", "can't be less than 5 characters")
	v.Check(input.Email != "", "email", "can't be blank")
	// TODO: Add requirements for stronger password
	v.Check(input.Password != "", "password", "can't be blank")

	v.Check(
		validator.Matches(
			input.Name,
			validator.NameRX,
		),
		"name",
		"does not look like a valid name",
	)

	v.Check(
		validator.Matches(
			input.Email,
			validator.EmailRX,
		),
		"email",
		"does not look like a valid email",
	)

	if v.Invalid() {
		app.failedValidationResponse(w, r, v.Errors)

		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)

		return
	}

	user := data.User{
		ID:        id,
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		Password:  "[HIDDEN]",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
