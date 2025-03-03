package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/betasve/go-commerce/services/auth/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserModel struct {
	DB *sql.DB
}

func (u UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return err
	}

	args := []interface{}{
		user.Name,
		user.Email,
		hashedPassword,
	}

	return u.DB.QueryRow(query, args...).Scan(&user.ID, &user.CreatedAt)
}

func (u UserModel) Get(id int64) (*User, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, updated_at, name, email
		FROM users
		WHERE id = $1
	`

	var user User

	err := u.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Name,
		&user.Email,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (u UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, password = $3
		WHERE id = $4
	`

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return err
	}

	args := []interface{}{
		user.Name,
		user.Email,
		hashedPassword,
		user.ID,
	}

	_, err = u.DB.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (u UserModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM users
		WHERE id = $1
	`

	result, err := u.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

type MockUserModel struct {
	DB *sql.DB
}

func (u MockUserModel) Insert(user *User) error {
	return nil
}

func (u MockUserModel) Get(id int64) (*User, error) {
	return nil, nil
}

func (u MockUserModel) Update(user *User) error {
	return nil
}

func (u MockUserModel) Delete(id int64) error {
	return nil
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

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}
