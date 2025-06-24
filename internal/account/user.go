package account

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
)

type User struct {
	ID       uuid.UUID
	Username string
	Email    string
	Password string
}

func NewUser(id uuid.UUID, username, email, password string) *User {
	return &User{
		ID:       id,
		Username: username,
		Email:    email,
		Password: password,
	}
}
