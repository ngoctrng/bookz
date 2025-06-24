package usecases

import (
	"github.com/google/uuid"
	"github.com/ngoctrng/bookz/internal/account"
)

type Repository interface {
	Save(user *account.User) error
	FindByEmail(email string) (*account.User, error)
}

type Usecase interface {
	Register(id uuid.UUID, username string, email string, password string) error
	Login(email string, password string) (*account.User, error)
}
