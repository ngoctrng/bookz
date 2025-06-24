package hasher_test

import (
	"testing"

	"github.com/ngoctrng/bookz/pkg/hasher"
	"github.com/stretchr/testify/assert"
)

func TestHashAndVerifyPassword(t *testing.T) {
	password := "supersecret"

	t.Run("should hash and verify password successfully", func(t *testing.T) {
		hashed, err := hasher.HashPassword(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hashed)

		err = hasher.VerifyPassword(hashed, password)
		assert.NoError(t, err)
	})

	t.Run("should fail verification with wrong password", func(t *testing.T) {
		hashed, err := hasher.HashPassword(password)
		assert.NoError(t, err)

		err = hasher.VerifyPassword(hashed, "wrongpassword")
		assert.Error(t, err)
	})
}
