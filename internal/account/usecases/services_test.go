package usecases_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ngoctrng/bookz/internal/account"
	"github.com/ngoctrng/bookz/internal/account/mocks"
	"github.com/ngoctrng/bookz/internal/account/usecases"
	"github.com/ngoctrng/bookz/pkg/hasher"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulUserRegistration(t *testing.T) {
	u := account.NewUser(uuid.New(), "testuser", "test@example.com", "securepassword")
	r := new(mocks.MockRepository)
	svc := usecases.NewService(r)

	t.Run("should register the new user successfully", func(t *testing.T) {
		r.EXPECT().Save(u).Return(nil).Once()
		r.EXPECT().FindByEmail(u.Email).Return(nil, nil).Once()

		err := svc.Register(u.ID, u.Username, u.Email, u.Password)

		assert.NoError(t, err)
		r.AssertExpectations(t)
	})
}

func TestDuplicateEmailRegistration(t *testing.T) {
	u := account.NewUser(uuid.New(), "testuser", "duplicate@example.com", "securepassword")
	r := new(mocks.MockRepository)
	svc := usecases.NewService(r)

	t.Run("should reject the duplicated email", func(t *testing.T) {
		r.EXPECT().FindByEmail(u.Email).Return(u, nil).Once()

		err := svc.Register(u.ID, u.Username, u.Email, u.Password)

		assert.Error(t, err)
		assert.ErrorAs(t, err, &account.ErrEmailAlreadyExists)
		r.AssertExpectations(t)
	})
}

func TestUserRegistrationFailure(t *testing.T) {
	u := account.NewUser(uuid.New(), "testuser", "test@example.com", "securepassword")
	r := new(mocks.MockRepository)
	svc := usecases.NewService(r)
	unexpectedErr := errors.New("unexpected error")

	t.Run("should reject when unexpected error happen", func(t *testing.T) {
		r.EXPECT().FindByEmail(u.Email).Return(nil, nil).Once()
		r.EXPECT().Save(u).Return(unexpectedErr).Once()

		err := svc.Register(u.ID, u.Username, u.Email, u.Password)

		assert.Error(t, err)
		assert.ErrorAs(t, err, &unexpectedErr)
		r.AssertExpectations(t)
	})
}

func TestLoginSuccess(t *testing.T) {
	hashed, err := hasher.HashPassword("123456")
	assert.NoError(t, err)
	u := account.NewUser(uuid.New(), "testuser", "test@example.com", hashed)
	r := new(mocks.MockRepository)
	svc := usecases.NewService(r)

	t.Run("should login successfully with correct credentials", func(t *testing.T) {
		r.EXPECT().FindByEmail(u.Email).Return(u, nil).Once()

		got, err := svc.Login(u.Email, "123456")

		assert.NoError(t, err)
		assert.Equal(t, u, got)
		r.AssertExpectations(t)
	})
}

func TestLoginUserNotFound(t *testing.T) {
	r := new(mocks.MockRepository)
	svc := usecases.NewService(r)

	t.Run("should reject login if user not found", func(t *testing.T) {
		r.EXPECT().FindByEmail("notfound@example.com").Return(nil, nil).Once()

		got, err := svc.Login("notfound@example.com", "any")

		assert.Error(t, err)
		assert.Nil(t, got)
		assert.ErrorAs(t, err, &account.ErrUserNotFound)
		r.AssertExpectations(t)
	})
}

func TestLoginFailure(t *testing.T) {
	u := account.NewUser(uuid.New(), "testuser", "test@example.com", "securepassword")
	r := new(mocks.MockRepository)
	svc := usecases.NewService(r)

	t.Run("should reject login if unexpected error happen", func(t *testing.T) {
		repoErr := errors.New("db error")
		r.EXPECT().FindByEmail(u.Email).Return(nil, repoErr).Once()

		got, err := svc.Login(u.Email, "any")

		assert.Error(t, err)
		assert.Nil(t, got)
		assert.Contains(t, err.Error(), "can not find user")
		r.AssertExpectations(t)
	})
}

func TestLoginInvalidPassword(t *testing.T) {
	u := account.NewUser(uuid.New(), "testuser", "test@example.com", "securepassword")
	r := new(mocks.MockRepository)
	svc := usecases.NewService(r)

	t.Run("should reject login if password is invalid", func(t *testing.T) {
		r.EXPECT().FindByEmail(u.Email).Return(u, nil).Once()

		got, err := svc.Login(u.Email, "wrongpassword")

		assert.Error(t, err)
		assert.Nil(t, got)
		r.AssertExpectations(t)
	})
}
