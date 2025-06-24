package account_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ngoctrng/bookz/internal/account"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	t.Run("should create user with given values", func(t *testing.T) {
		id := uuid.New()
		username := "testuser"
		email := "test@example.com"
		password := "password123"

		user := account.NewUser(id, username, email, password)

		assert.Equal(t, id, user.ID)
		assert.Equal(t, username, user.Username)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, password, user.Password)
	})
}
