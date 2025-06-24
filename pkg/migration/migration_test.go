package migration_test

import (
	"testing"

	"github.com/ngoctrng/bookz/pkg/migration"
	"github.com/ngoctrng/bookz/pkg/testutil"

	"github.com/stretchr/testify/assert"
)

func TestMigration(t *testing.T) {
	dbName, dbUser, dbPass := "test1", "test1", "123456"
	db := testutil.CreateConnection(t, dbName, dbUser, dbPass)
	conn, _ := db.DB()

	_, err := migration.Run(conn)
	assert.NoError(t, err)
}
