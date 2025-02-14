package data

import (
	"time"

	"github.com/betasve/go-commerce/services/auth/internal/validator"
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ValidateUser(v *validator.Validator, u *User) {
	v.Check(u.Name != "", "name", "can't be blank")
	v.Check(len(u.Name) > 5, "name", "can't be less than 5 characters")
	v.Check(u.Email != "", "email", "can't be blank")
	// TODO: Add requirements for stronger password
	v.Check(u.Password != "", "password", "can't be blank")

	v.Check(
		validator.Matches(
			u.Name,
			validator.NameRX,
		),
		"name",
		"does not look like a valid name",
	)

	v.Check(
		validator.Matches(
			u.Email,
			validator.EmailRX,
		),
		"email",
		"does not look like a valid email",
	)
}
