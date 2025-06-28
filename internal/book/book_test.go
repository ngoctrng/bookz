package book_test

import (
	"testing"

	"github.com/ngoctrng/bookz/internal/book"
	"github.com/stretchr/testify/assert"
)

func TestNewBook(t *testing.T) {
	b := book.New("owner1", "1234567890", "Go Programming", "Alan A. A. Donovan", 2024)
	assert.Equal(t, "owner1", b.OwnerID)
	assert.Equal(t, "1234567890", b.ISBN)
	assert.Equal(t, "Go Programming", b.Title)
	assert.Equal(t, "Alan A. A. Donovan", b.Author)
	assert.Equal(t, 2024, b.Year)
	assert.Empty(t, b.Description)
	assert.Empty(t, b.BriefReview)
}

func TestAddDescription(t *testing.T) {
	b := book.New("owner1", "1234567890", "Go Programming", "Alan", 2024)
	b.AddDescription("A book about Go.")
	assert.Equal(t, "A book about Go.", b.Description)
}

func TestAddBriefReview(t *testing.T) {
	b := book.New("owner1", "1234567890", "Go Programming", "Alan", 2024)
	b.AddBriefReview("Excellent resource.")
	assert.Equal(t, "Excellent resource.", b.BriefReview)
}

func TestChangeOwner(t *testing.T) {
	b := book.New("owner1", "1234567890", "Go Programming", "Alan", 2024)
	b.ChangeOwner("owner2")
	assert.Equal(t, "owner2", b.OwnerID)
}

func TestChangeOwnerFor(t *testing.T) {
	b1 := book.New("owner1", "111", "Book One", "Author1", 2020)
	b2 := book.New("owner2", "222", "Book Two", "Author2", 2021)

	origOwner1, origOwner2 := b1.OwnerID, b2.OwnerID

	returnedOwner1, returnedOwner2 := book.ChangeOwnerFor(b1, b2)

	assert.Equal(t, origOwner1, returnedOwner1)
	assert.Equal(t, origOwner2, returnedOwner2)
	assert.Equal(t, origOwner2, b1.OwnerID)
	assert.Equal(t, origOwner1, b2.OwnerID)
}
