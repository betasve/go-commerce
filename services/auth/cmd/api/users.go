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

	user := &data.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	v := validator.New()

	if data.ValidateUser(v, user); v.Invalid() {
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
