package repository_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ngoctrng/bookz/internal/account"
	"github.com/ngoctrng/bookz/internal/account/repository"
	"github.com/ngoctrng/bookz/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSaveUser(t *testing.T) {
	dbName, dbUser, dbPass := "test2", "test2", "123456"
	db := testutil.CreateConnection(t, dbName, dbUser, dbPass)
	conn, _ := db.DB()
	testutil.MigrateTestDatabase(t, conn)

	r := repository.New(db)

	t.Run("should successfully inserted into users table", func(t *testing.T) {
		u := account.NewUser(uuid.New(), "ngoctest", "john.doe@example.com", "password123")
		err := r.Save(u)
		assert.NoError(t, err)
		verifyInsertedUser(t, db, u)
	})
}

func TestFindByEmail(t *testing.T) {
	dbName, dbUser, dbPass := "test2", "test2", "123456"
	db := testutil.CreateConnection(t, dbName, dbUser, dbPass)
	conn, _ := db.DB()
	testutil.MigrateTestDatabase(t, conn)

	r := repository.New(db)

	t.Run("should return user by email", func(t *testing.T) {
		schema := setupExistingUser(t, db)
		email := "john.doe@example.com"

		user, err := r.FindByEmail(email)

		assert.NoError(t, err)
		assertUserByEmail(t, user, schema, email)
	})

	t.Run("should return nil if not found", func(t *testing.T) {
		email := "notfound@example.com"
		user, err := r.FindByEmail(email)
		assert.NoError(t, err)
		assert.Nil(t, user)
	})
}

func assertUserByEmail(t *testing.T, user *account.User, schema repository.UserSchema, email string) {
	assert.NotNil(t, user)
	assert.Equal(t, schema.ID, user.ID.String())
	assert.Equal(t, email, user.Email)
	assert.Equal(t, schema.Username, user.Username)
}

func setupExistingUser(t *testing.T, db *gorm.DB) repository.UserSchema {
	u := repository.UserSchema{
		ID:           uuid.New().String(),
		Username:     "ngoctest",
		Email:        "john.doe@example.com",
		PasswordHash: "password123",
	}
	err := db.Table("users").Create(&u).Error
	assert.NoError(t, err)
	return u
}

func verifyInsertedUser(t *testing.T, db *gorm.DB, u *account.User) {
	var id string
	db.Table("users").Select("id").Where("id = ?", u.ID).Scan(&id)
	assert.Equal(t, u.ID.String(), id)
}
