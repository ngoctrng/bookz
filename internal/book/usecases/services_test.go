package usecases_test

import (
	"errors"
	"testing"

	"github.com/ngoctrng/bookz/internal/book"
	"github.com/ngoctrng/bookz/internal/book/mocks"
	"github.com/ngoctrng/bookz/internal/book/usecases"
	"github.com/stretchr/testify/assert"
)

func TestCreateBook(t *testing.T) {
	b := &book.Book{
		OwnerID: "owner1",
		ISBN:    "isbn-001",
		Title:   "Book Title",
		Author:  "Author",
		Year:    2024,
	}
	r := new(mocks.MockRepository)
	svc := usecases.NewService(r)

	t.Run("should create book successfully", func(t *testing.T) {
		r.EXPECT().Save(b).Return(nil).Once()
		err := svc.Create(b)
		assert.NoError(t, err)
		r.AssertExpectations(t)
	})

	t.Run("should return error if repo fails", func(t *testing.T) {
		r.EXPECT().Save(b).Return(errors.New("db error")).Once()
		err := svc.Create(b)
		assert.Error(t, err)
		r.AssertExpectations(t)
	})
}

func TestGetBook(t *testing.T) {
	b := &book.Book{ISBN: "isbn-001", Title: "Book Title"}
	r := new(mocks.MockRepository)
	svc := usecases.NewService(r)

	t.Run("should get book by isbn", func(t *testing.T) {
		r.EXPECT().FindByISBN("isbn-001").Return(b, nil).Once()
		got, err := svc.Get("isbn-001")
		assert.NoError(t, err)
		assert.Equal(t, b, got)
		r.AssertExpectations(t)
	})

	t.Run("should return error if repo fails", func(t *testing.T) {
		r.EXPECT().FindByISBN("isbn-001").Return(nil, errors.New("db error")).Once()
		got, err := svc.Get("isbn-001")
		assert.Error(t, err)
		assert.Nil(t, got)
		r.AssertExpectations(t)
	})
}

func TestUpdateBook(t *testing.T) {
	origin := &book.Book{
		OwnerID: "owner1",
		ISBN:    "isbn-001",
		Title:   "Old Title",
	}
	updated := &book.Book{
		OwnerID: "owner1",
		ISBN:    "isbn-001",
		Title:   "New Title",
	}
	r := new(mocks.MockRepository)
	svc := usecases.NewService(r)

	t.Run("should update book if owner matches", func(t *testing.T) {
		r.EXPECT().FindByISBN("isbn-001").Return(origin, nil).Once()
		r.EXPECT().Update(updated).Return(nil).Once()
		err := svc.Update(updated, "owner1")
		assert.NoError(t, err)
		r.AssertExpectations(t)
	})

	t.Run("should not update if owner does not match", func(t *testing.T) {
		r.EXPECT().FindByISBN("isbn-001").Return(origin, nil).Once()
		err := svc.Update(updated, "other-owner")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "permission")
		r.AssertExpectations(t)
	})

	t.Run("should return error if repo fails", func(t *testing.T) {
		r.EXPECT().FindByISBN("isbn-001").Return(nil, errors.New("db error")).Once()
		err := svc.Update(updated, "owner1")
		assert.Error(t, err)
		r.AssertExpectations(t)
	})
}

func TestDeleteBook(t *testing.T) {
	b := &book.Book{
		OwnerID: "owner1",
		ISBN:    "isbn-001",
		Title:   "Book Title",
	}
	r := new(mocks.MockRepository)
	svc := usecases.NewService(r)

	t.Run("should delete book if owner matches", func(t *testing.T) {
		r.EXPECT().FindByISBN("isbn-001").Return(b, nil).Once()
		r.EXPECT().Delete("isbn-001").Return(nil).Once()
		err := svc.Delete("isbn-001", "owner1")
		assert.NoError(t, err)
		r.AssertExpectations(t)
	})

	t.Run("should not delete if owner does not match", func(t *testing.T) {
		r.EXPECT().FindByISBN("isbn-001").Return(b, nil).Once()
		err := svc.Delete("isbn-001", "other-owner")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "permission")
		r.AssertExpectations(t)
	})

	t.Run("should return error if repo fails", func(t *testing.T) {
		r.EXPECT().FindByISBN("isbn-001").Return(nil, errors.New("db error")).Once()
		err := svc.Delete("isbn-001", "owner1")
		assert.Error(t, err)
		r.AssertExpectations(t)
	})
}

func TestListBooks(t *testing.T) {
	books := []*book.BookInfo{
		{ISBN: "1", Title: "Book 1", Author: "Author 1", Year: 2021, Owner: book.BookOwner{OwnerID: "owner1"}},
		{ISBN: "2", Title: "Book 2", Author: "Author 2", Year: 2022, Owner: book.BookOwner{OwnerID: "owner2"}},
	}
	r := new(mocks.MockRepository)
	svc := usecases.NewService(r)

	t.Run("should list books", func(t *testing.T) {
		r.EXPECT().List().Return(books, nil).Once()
		got, err := svc.List()
		assert.NoError(t, err)
		assert.Equal(t, books, got)
		r.AssertExpectations(t)
	})

	t.Run("should return error if repo fails", func(t *testing.T) {
		r.EXPECT().List().Return(nil, errors.New("db error")).Once()
		got, err := svc.List()
		assert.Error(t, err)
		assert.Nil(t, got)
		r.AssertExpectations(t)
	})
}
