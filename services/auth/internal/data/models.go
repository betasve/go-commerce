package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Users interface {
		Insert(user *User) error
		Get(id int64) (*User, error)
		GetAll(email, name string, filters Filters) ([]*User, Metadata, error)
		Update(user *User) error
		Delete(id int64) error
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users: UserModel{DB: db},
	}
}

func NewMockModels() Models {
	return Models{
		Users: MockUserModel{},
	}
}
