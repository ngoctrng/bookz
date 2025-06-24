package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/ngoctrng/bookz/pkg/migration"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func MigrateTestDatabase(t testing.TB, db *sql.DB) {
	t.Helper()

	_, err := migration.Run(db)
	assert.NoError(t, err)
}

func CreateConnection(t testing.TB, dbName string, dbUser string, dbPass string) *gorm.DB {
	cont := SetupPostgresContainer(t, dbName, dbUser, dbPass)
	host, _ := cont.Host(context.Background())
	port, _ := cont.MappedPort(context.Background(), "5432")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		host, dbUser, dbPass, dbName, port.Int())
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	return db
}
