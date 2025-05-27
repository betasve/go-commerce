package data

import (
	"context"
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return u.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt)
}

func (u UserModel) Get(id int64) (*User, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT pg_sleep(10), id, created_at, updated_at, name, email
		FROM users
		WHERE id = $1
	`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, id).Scan(
		&[]byte{},
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

func (u UserModel) GetAll(email, name string, filters Filters) ([]*User, error) {
	query := `
		SELECT id, email, name, created_at, updated_at
		FROM users
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (LOWER(email) = LOWER($2) OR $2 = '')
		ORDER BY id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := u.DB.QueryContext(ctx, query, name, email)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []*User{}

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (u UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, password = $3
		WHERE id = $4 AND updated_at = $5
		RETURNING updated_at
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
		user.UpdatedAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = u.DB.QueryRowContext(ctx, query, args...).Scan(&user.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := u.DB.ExecContext(ctx, query, id)
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
	user.ID = 42
	user.Name = "John Doe"
	user.Email = "test_email@example.com"

	t, err := time.Parse("2006-01-02 15:04:05", "2025-03-26 15:04:05")
	if err != nil {
		return err
	}
	user.CreatedAt = t
	user.UpdatedAt = t

	return nil
}

func (u MockUserModel) Get(id int64) (*User, error) {
	user := &User{}
	user.ID = 42
	user.Name = "John Doe"
	user.Email = "test_email@example.com"
	user.Password = "TestPassword321"

	t, err := time.Parse("2006-01-02 15:04:05", "2025-03-26 15:04:05")
	if err != nil {
		return nil, err
	}
	user.CreatedAt = t
	user.UpdatedAt = t

	return user, nil
}

func (u MockUserModel) GetAll(email, name string, filters Filters) ([]*User, error) {
	t, err := time.Parse("2006-01-02 15:04:05", "2025-03-26 15:04:05")
	if err != nil {
		return nil, err
	}

	return []*User{
		{ID: 1, Email: "test@example.com", Name: "John Doe", CreatedAt: t, UpdatedAt: t},
		{ID: 1, Email: "test2@example.com", Name: "Jill Doe", CreatedAt: t, UpdatedAt: t},
	}, nil
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
