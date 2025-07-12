package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/betasve/go-commerce/services/auth/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicated email")
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

type UserModel struct {
	DB *sql.DB
}

func NewUser(name, email, pwd string) *User {
	return &User{
		Name:     name,
		Email:    email,
		Password: password{plaintext: &pwd, hash: []byte{}},
	}

}

func (u UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := user.Password.Set(*user.Password.plaintext)
	if err != nil {
		return err
	}

	args := []interface{}{
		user.Name,
		user.Email,
		user.Password.hash,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = u.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, id).Scan(
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

func (u UserModel) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, created_at, updated_at, name, email, password
		FROM users
		WHERE email = $1
	`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
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

func (u UserModel) GetAll(email, name string, filters Filters) ([]*User, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, email, name, created_at, updated_at
		FROM users
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (LOWER(email) = LOWER($2) OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := u.DB.QueryContext(ctx, query, name, email, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	users := []*User{}

	for rows.Next() {
		var user User

		err := rows.Scan(
			&totalRecords,
			&user.ID,
			&user.Email,
			&user.Name,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return users, metadata, nil
}

func (u UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, password = $3
		WHERE id = $4 AND updated_at = $5
		RETURNING updated_at
	`
	err := user.Password.Set(*user.Password.plaintext)
	if err != nil {
		return err
	}

	args := []interface{}{
		user.Name,
		user.Email,
		user.Password.hash,
		user.ID,
		user.UpdatedAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = u.DB.QueryRowContext(ctx, query, args...).Scan(&user.UpdatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
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
	plaintextPassword := "TestPassword321"
	user.Password = password{&plaintextPassword, []byte{}}

	t, err := time.Parse("2006-01-02 15:04:05", "2025-03-26 15:04:05")
	if err != nil {
		return nil, err
	}
	user.CreatedAt = t
	user.UpdatedAt = t

	return user, nil
}

func (u MockUserModel) GetAll(email, name string, filters Filters) ([]*User, Metadata, error) {
	t, err := time.Parse("2006-01-02 15:04:05", "2025-03-26 15:04:05")
	if err != nil {
		return nil, Metadata{}, err
	}

	return []*User{
		{ID: 1, Email: "test@example.com", Name: "John Doe", CreatedAt: t, UpdatedAt: t},
		{ID: 1, Email: "test2@example.com", Name: "Jill Doe", CreatedAt: t, UpdatedAt: t},
	}, Metadata{1, 20, 1, 2, 40}, nil
}

func (u MockUserModel) Update(user *User) error {
	return nil
}

func (u MockUserModel) Delete(id int64) error {
	return nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "can't be blank")
	v.Check(
		validator.Matches(
			email,
			validator.EmailRX,
		),
		"email",
		"does not look like a valid email",
	)
}

func ValidatePasswordPlaintext(v *validator.Validator, plaintext string) {
	// TODO: Add requirements for stronger password
	v.Check(plaintext != "", "password", "can't be blank")
	v.Check(len(plaintext) >= 8, "password", "too short")
	v.Check(len(plaintext) <= 72, "password", "too long")
}

func ValidateUser(v *validator.Validator, u *User) {
	v.Check(u.Name != "", "name", "can't be blank")
	v.Check(len(u.Name) > 5, "name", "can't be less than 5 characters")
	v.Check(
		validator.Matches(
			u.Name,
			validator.NameRX,
		),
		"name",
		"does not look like a valid name",
	)

	ValidateEmail(v, u.Email)

	if u.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *u.Password.plaintext)
	}

	if u.Password.hash == nil {
		panic("missing password hash")
	}
}
