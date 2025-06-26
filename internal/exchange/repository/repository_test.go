package repository_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/ngoctrng/bookz/internal/exchange"
	"github.com/ngoctrng/bookz/internal/exchange/repository"
	"github.com/ngoctrng/bookz/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestExchangeRepository(t *testing.T) {
	dbName, dbUser, dbPass := "test2", "test2", "123456"
	db := testutil.CreateConnection(t, dbName, dbUser, dbPass)
	conn, _ := db.DB()
	testutil.MigrateTestDatabase(t, conn)

	r := repository.New(db)

	requestBy := uuid.New()
	requestedID := 1

	t.Run("should save and get proposal by ID", func(t *testing.T) {
		p := &exchange.Proposal{
			ID:            1,
			RequestBy:     requestBy,
			RequestedID:   exchange.BookID(requestedID),
			ForExchangeID: 42,
			Message:       "Let's trade!",
			Status:        "pending",
		}
		err := r.Save(p)
		assert.NoError(t, err)

		found, err := r.GetByID(1)
		assert.NoError(t, err)
		assertProposalEqual(t, p, found)
	})

	t.Run("should get all proposals for user", func(t *testing.T) {
		proposals, err := r.GetAll(requestBy)
		assert.NoError(t, err)
		assert.True(t, len(proposals) >= 1)
		found := false
		for _, prop := range proposals {
			if prop.ID == 1 {
				found = true
				assertProposalEqual(t, &exchange.Proposal{
					ID:            1,
					RequestBy:     requestBy,
					RequestedID:   exchange.BookID(requestedID),
					ForExchangeID: 42,
					Message:       "Let's trade!",
					Status:        "pending",
				}, prop)
			}
		}
		assert.True(t, found, "inserted proposal should be in the list")
	})

	t.Run("should return nil if not found", func(t *testing.T) {
		found, err := r.GetByID(999)
		assert.NoError(t, err)
		assert.Nil(t, found)
	})

	t.Run("should fetch requested book owner", func(t *testing.T) {
		requestBy := uuid.New()
		requestedID := 1

		givenBookInRepository(t, db, requestedID, requestBy)

		ownerID, err := r.FetchRequestedBookOwner(requestedID)

		assert.NoError(t, err)
		assert.Equal(t, requestBy, ownerID)
	})

	t.Run("should return nil uuid if not found", func(t *testing.T) {
		ownerID, err := r.FetchRequestedBookOwner(-1)
		assert.NoError(t, err)
		assert.Equal(t, uuid.Nil, ownerID)
	})
}

func assertProposalEqual(t *testing.T, expected, actual *exchange.Proposal) {
	assert.NotNil(t, actual)
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.RequestBy, actual.RequestBy)
	assert.Equal(t, expected.RequestedID, actual.RequestedID)
	assert.Equal(t, expected.ForExchangeID, actual.ForExchangeID)
	assert.Equal(t, expected.Message, actual.Message)
	assert.Equal(t, expected.Status, actual.Status)
}

func givenBookInRepository(t *testing.T, db *gorm.DB, bookID int, ownerID uuid.UUID) {
	t.Helper()

	sql := fmt.Sprintf(`
		INSERT INTO books (id, isbn, owner_id, title, author) 
		VALUES (%d, '9780143127741', '%s', 'The Martian', 'Andy Weir')`, bookID, ownerID.String())

	inserted := db.Exec(sql)

	assert.NoError(t, inserted.Error, "should insert book into repository")
	assert.Equal(t, inserted.RowsAffected, int64(1))
}
