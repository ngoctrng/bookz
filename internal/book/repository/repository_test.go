package repository_test

import (
	"testing"

	"github.com/ngoctrng/bookz/internal/book"
	"github.com/ngoctrng/bookz/internal/book/repository"
	"github.com/ngoctrng/bookz/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestBookRepository(t *testing.T) {
	dbName, dbUser, dbPass := "test2", "test2", "123456"
	db := testutil.CreateConnection(t, dbName, dbUser, dbPass)
	conn, _ := db.DB()
	testutil.MigrateTestDatabase(t, conn)

	r := repository.New(db)

	t.Run("should save and find book by ISBN", func(t *testing.T) {
		b := &book.Book{
			OwnerID:     "owner1",
			ISBN:        "isbn-001",
			Title:       "Book Title",
			Description: "A book description",
			BriefReview: "A brief review",
			Author:      "Author Name",
			Year:        2024,
		}
		err := r.Save(b)
		assert.NoError(t, err)

		found, err := r.FindByISBN("isbn-001")
		assert.NoError(t, err)
		assertBookEqual(t, b, found)
	})

	t.Run("should update book", func(t *testing.T) {
		b := &book.Book{
			OwnerID:     "owner1",
			ISBN:        "isbn-001",
			Title:       "Updated Title",
			Description: "Updated description",
			BriefReview: "Updated review",
			Author:      "Updated Author",
			Year:        2025,
		}
		err := r.Update(b)
		assert.NoError(t, err)

		found, err := r.FindByISBN("isbn-001")
		assert.NoError(t, err)
		assertBookEqual(t, b, found)
	})

	t.Run("should list books", func(t *testing.T) {
		books, err := r.List()
		assert.NoError(t, err)
		assert.True(t, len(books) >= 1)
	})

	t.Run("should delete book", func(t *testing.T) {
		err := r.Delete("isbn-001")
		assert.NoError(t, err)

		found, err := r.FindByISBN("isbn-001")
		assert.NoError(t, err)
		assert.Nil(t, found)
	})

	t.Run("should return nil if not found", func(t *testing.T) {
		found, err := r.FindByISBN("not-exist")
		assert.NoError(t, err)
		assert.Nil(t, found)
	})
}

func assertBookEqual(t *testing.T, expected, actual *book.Book) {
	assert.NotNil(t, actual)
	assert.Equal(t, expected.OwnerID, actual.OwnerID)
	assert.Equal(t, expected.ISBN, actual.ISBN)
	assert.Equal(t, expected.Title, actual.Title)
	assert.Equal(t, expected.Description, actual.Description)
	assert.Equal(t, expected.BriefReview, actual.BriefReview)
	assert.Equal(t, expected.Author, actual.Author)
	assert.Equal(t, expected.Year, actual.Year)
}
