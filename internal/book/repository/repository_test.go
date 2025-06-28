package repository_test

import (
	"testing"

	"github.com/google/uuid"
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

	const ownerID = "c2ee4e49-2c77-452d-9a1b-5765cea6f10e"
	db.Raw(`INSERT INTO users (id, username, email, password) VALUES 
		('c2ee4e49-2c77-452d-9a1b-5765cea6f10e', 'owner1', 'owner1@example.com', 'password')`)
	assert.NoError(t, db.Error)

	t.Run("should save and find book by ISBN", func(t *testing.T) {
		b := &book.Book{
			ID:          1,
			OwnerID:     ownerID,
			ISBN:        "isbn-001",
			Title:       "Book Title",
			Description: "A book description",
			BriefReview: "A brief review",
			Author:      "Author Name",
			Year:        2024,
		}
		err := r.Save(b)
		assert.NoError(t, err)

		found, err := r.FindByID(1)
		assert.NoError(t, err)
		assertBookInfo(t, b, found)
	})

	t.Run("should update book", func(t *testing.T) {
		b := &book.Book{
			ID:          1,
			OwnerID:     ownerID,
			ISBN:        "isbn-001",
			Title:       "Updated Title",
			Description: "Updated description",
			BriefReview: "Updated review",
			Author:      "Updated Author",
			Year:        2025,
		}
		err := r.Update(b)
		assert.NoError(t, err)

		found, err := r.FindByID(1)
		assert.NoError(t, err)
		assertBookInfo(t, b, found)
	})

	t.Run("should list books", func(t *testing.T) {
		books, err := r.List()
		assert.NoError(t, err)
		assert.True(t, len(books) >= 1)

		// Check that the returned BookInfo matches the inserted book
		for _, info := range books {
			if info.ISBN == "isbn-001" {
				assert.Equal(t, "Updated Title", info.Title)
				assert.Equal(t, "Updated description", info.Description)
				assert.Equal(t, "Updated review", info.BriefReview)
				assert.Equal(t, "Updated Author", info.Author)
				assert.Equal(t, 2025, info.Year)
				assert.Equal(t, ownerID, info.Owner.ID)
			}
		}
	})

	t.Run("should delete book", func(t *testing.T) {
		err := r.Delete(1)
		assert.NoError(t, err)

		found, err := r.FindByID(1)
		assert.NoError(t, err)
		assert.Nil(t, found)
	})

	t.Run("should return nil if not found", func(t *testing.T) {
		found, err := r.FindByID(999)
		assert.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestUpsert(t *testing.T) {
	db := testutil.CreateConnection(t, "testdb", "testuser", "testpass")
	conn, _ := db.DB()
	testutil.MigrateTestDatabase(t, conn)
	r := repository.New(db)

	b1 := &book.Book{OwnerID: uuid.NewString(), ISBN: "isbn1", Title: "Book1", Author: "A", Year: 2020}
	b2 := &book.Book{OwnerID: uuid.NewString(), ISBN: "isbn2", Title: "Book2", Author: "B", Year: 2021}
	givenBooksInRepository(t, r, []*book.Book{b1, b2})

	books := []*book.Book{b1, b2}
	var updateBooks []*book.Book
	for _, info := range books {
		b := &book.Book{
			ID:      info.ID,
			OwnerID: info.OwnerID,
			ISBN:    info.ISBN,
			Title:   info.Title + " Updated",
			Author:  info.Author,
			Year:    info.Year + 1,
		}
		updateBooks = append(updateBooks, b)
	}

	err := r.Upsert(updateBooks)

	assert.NoError(t, err)
}

func TestGetBy(t *testing.T) {
	db := testutil.CreateConnection(t, "testdb", "testuser", "testpass")
	conn, _ := db.DB()
	testutil.MigrateTestDatabase(t, conn)
	r := repository.New(db)

	b1 := &book.Book{ID: 1, OwnerID: uuid.NewString(), ISBN: "isbn1", Title: "Book1", Author: "A", Year: 2020}
	b2 := &book.Book{ID: 2, OwnerID: uuid.NewString(), ISBN: "isbn2", Title: "Book2", Author: "B", Year: 2021}
	givenBooksInRepository(t, r, []*book.Book{b1, b2})

	result, err := r.GetBy([]int{b1.ID, b2.ID})
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	for _, book := range result {
		assertBook(t, book, book)
	}
}

func TestGetProposal(t *testing.T) {
	db := testutil.CreateConnection(t, "testdb", "testuser", "testpass")
	conn, _ := db.DB()
	testutil.MigrateTestDatabase(t, conn)
	r := repository.New(db)

	result := db.Exec(`INSERT INTO proposals (id,request_by,request_to,requested_id,for_exchange_id,status,requested_at) 
					  VALUES (1,?,?,1,2,'pending',NOW())`, uuid.NewString(), uuid.NewString())
	assert.NoError(t, result.Error)
	assert.Equal(t, int64(1), result.RowsAffected)

	t.Run("should return proposal details if found", func(t *testing.T) {
		got := r.GetProposal(1)
		assert.NotNil(t, got)
		assert.Equal(t, 1, got.ID)
		assert.Equal(t, 1, got.RequestedID)
		assert.Equal(t, 2, got.ForExchangeID)
	})

	t.Run("should return nil if proposal not found", func(t *testing.T) {
		got := r.GetProposal(999)
		assert.Nil(t, got)
	})
}

func assertBook(t *testing.T, expected *book.Book, actual *book.Book) {
	assert.NotNil(t, actual)
	assert.Equal(t, expected.ISBN, actual.ISBN)
	assert.Equal(t, expected.Title, actual.Title)
	assert.Equal(t, expected.Description, actual.Description)
	assert.Equal(t, expected.BriefReview, actual.BriefReview)
	assert.Equal(t, expected.Author, actual.Author)
	assert.Equal(t, expected.Year, actual.Year)
}

func assertBookInfo(t *testing.T, expected *book.Book, actual *book.BookInfo) {
	assert.Equal(t, expected.OwnerID, actual.Owner.ID)
	assertBook(t, expected, &book.Book{
		OwnerID:     actual.Owner.ID,
		ISBN:        actual.ISBN,
		Title:       actual.Title,
		Description: actual.Description,
		BriefReview: actual.BriefReview,
		Author:      actual.Author,
		Year:        actual.Year,
	})
}

func givenBooksInRepository(t *testing.T, r *repository.Repository, books []*book.Book) {
	t.Helper()
	for _, b := range books {
		assert.NoError(t, r.Save(b))
	}
}
