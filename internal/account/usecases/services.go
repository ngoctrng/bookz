package usecases

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ngoctrng/bookz/internal/account"
	"github.com/ngoctrng/bookz/pkg/hasher"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Register(id uuid.UUID, username string, email string, password string) error {
	u, err := s.repo.FindByEmail(email)
	if err != nil {
		return errors.Join(fmt.Errorf("can not find user '%s'", email), err)
	}
	if u != nil {
		return errors.Join(account.ErrEmailAlreadyExists, err)
	}

	err = s.repo.Save(account.NewUser(id, username, email, password))
	if err != nil {
		return errors.Join(fmt.Errorf("can not save user '%s'", email), err)
	}

	return nil
}

func (s *Service) Login(email string, password string) (*account.User, error) {
	u, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("can not find user '%s'", email), err)
	}
	if u == nil {
		return nil, errors.Join(account.ErrUserNotFound, fmt.Errorf("user '%s' not found", email))
	}

	if err := hasher.VerifyPassword(u.Password, password); err != nil {
		return nil, errors.Join(fmt.Errorf("password for user '%s' is incorrect", email), err)
	}

	return u, nil
}
