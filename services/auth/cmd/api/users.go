package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/betasve/go-commerce/services/auth/internal/data"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "placeholder for creating a user")
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
