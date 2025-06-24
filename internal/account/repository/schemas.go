package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/ngoctrng/bookz/internal/account"
)

type UserSchema struct {
	ID           string
	Email        string
	Username     string
	PasswordHash string
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func DomainToUserSchema(user *account.User) *UserSchema {
	return &UserSchema{
		ID:           user.ID.String(),
		Email:        user.Email,
		Username:     user.Username,
		PasswordHash: user.Password,
	}
}

func UserSchemaToDomain(user *UserSchema) *account.User {
	return account.NewUser(
		uuid.MustParse(user.ID),
		user.Username,
		user.Email,
		user.PasswordHash,
	)
}
