package token_test

import (
	"testing"
	"time"

	"github.com/ngoctrng/bookz/pkg/token"
	"github.com/stretchr/testify/assert"
)

func TestSignAndVerify(t *testing.T) {
	secret := "testsecret"
	userID := "user-123"
	duration := time.Minute

	t.Run("valid token", func(t *testing.T) {
		tokenStr, err := token.Sign(userID, secret, duration)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenStr)

		claims, err := token.Verify(tokenStr, secret)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.WithinDuration(t, time.Now(), claims.IssuedAt.Time, time.Second)
		assert.WithinDuration(t, time.Now().Add(duration), claims.ExpiresAt.Time, time.Second*2)
	})

	t.Run("invalid secret", func(t *testing.T) {
		tokenStr, err := token.Sign(userID, secret, duration)
		assert.NoError(t, err)

		claims, err := token.Verify(tokenStr, "wrongsecret")
		assert.ErrorIs(t, err, token.ErrInvalidToken)
		assert.Nil(t, claims)
	})

	t.Run("expired token", func(t *testing.T) {
		tokenStr, err := token.Sign(userID, secret, -time.Minute)
		assert.NoError(t, err)

		claims, err := token.Verify(tokenStr, secret)
		assert.ErrorIs(t, err, token.ErrInvalidToken)
		assert.Nil(t, claims)
	})

	t.Run("malformed token", func(t *testing.T) {
		claims, err := token.Verify("not.a.jwt", secret)
		assert.ErrorIs(t, err, token.ErrInvalidToken)
		assert.Nil(t, claims)
	})
}
