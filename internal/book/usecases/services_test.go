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
	b := &book.BookInfo{ID: 1, ISBN: "isbn-001", Title: "Book Title"}
	r := new(mocks.MockRepository)
	svc := usecases.NewService(r)

	t.Run("should get book by id", func(t *testing.T) {
		r.EXPECT().FindByID(1).Return(b, nil).Once()
		got, err := svc.Get(1)
		assert.NoError(t, err)
		assert.Equal(t, b, got)
		r.AssertExpectations(t)
	})

	t.Run("should return error if repo fails", func(t *testing.T) {
		r.EXPECT().FindByID(1).Return(nil, errors.New("db error")).Once()
		got, err := svc.Get(1)
		assert.Error(t, err)
		assert.Nil(t, got)
		r.AssertExpectations(t)
	})
}

func TestUpdateBook(t *testing.T) {
	origin := &book.BookInfo{
		ID:    1,
		Owner: book.BookOwner{ID: "owner1"},
		ISBN:  "isbn-001",
		Title: "Old Title",
	}
	updated := &book.Book{
		ID:      1,
		OwnerID: "owner1",
		ISBN:    "isbn-001",
		Title:   "New Title",
	}
	r := new(mocks.MockRepository)
	svc := usecases.NewService(r)

	t.Run("should update book if owner matches", func(t *testing.T) {
		r.EXPECT().FindByID(1).Return(origin, nil).Once()
		r.EXPECT().Update(updated).Return(nil).Once()
		err := svc.Update(updated, "owner1")
		assert.NoError(t, err)
		r.AssertExpectations(t)
	})

	t.Run("should not update if owner does not match", func(t *testing.T) {
		r.EXPECT().FindByID(1).Return(origin, nil).Once()
		err := svc.Update(updated, "other-owner")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "permission")
		r.AssertExpectations(t)
	})

	t.Run("should return error if repo fails", func(t *testing.T) {
		r.EXPECT().FindByID(1).Return(nil, errors.New("db error")).Once()
		err := svc.Update(updated, "owner1")
		assert.Error(t, err)
		r.AssertExpectations(t)
	})
}

func TestDeleteBook(t *testing.T) {
	b := &book.BookInfo{
		ID:    1,
		Owner: book.BookOwner{ID: "owner1"},
		ISBN:  "isbn-001",
		Title: "Book Title",
	}
	r := new(mocks.MockRepository)
	svc := usecases.NewService(r)

	t.Run("should delete book if owner matches", func(t *testing.T) {
		r.EXPECT().FindByID(1).Return(b, nil).Once()
		r.EXPECT().Delete(1).Return(nil).Once()
		err := svc.Delete(1, "owner1")
		assert.NoError(t, err)
		r.AssertExpectations(t)
	})

	t.Run("should not delete if owner does not match", func(t *testing.T) {
		r.EXPECT().FindByID(1).Return(b, nil).Once()
		err := svc.Delete(1, "other-owner")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "permission")
		r.AssertExpectations(t)
	})

	t.Run("should return error if repo fails", func(t *testing.T) {
		r.EXPECT().FindByID(1).Return(nil, errors.New("db error")).Once()
		err := svc.Delete(1, "owner1")
		assert.Error(t, err)
		r.AssertExpectations(t)
	})
}

func TestListBooks(t *testing.T) {
	books := []*book.BookInfo{
		{ID: 1, ISBN: "1", Title: "Book 1", Author: "Author 1", Year: 2021, Owner: book.BookOwner{ID: "owner1"}},
		{ID: 2, ISBN: "2", Title: "Book 2", Author: "Author 2", Year: 2022, Owner: book.BookOwner{ID: "owner2"}},
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

func TestFulfillProposal(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	svc := usecases.NewService(mockRepo)

	proposal := &book.ProposalDetails{
		RequestedID:   1,
		ForExchangeID: 2,
	}
	book1 := &book.Book{ID: 1, OwnerID: "owner1"}
	book2 := &book.Book{ID: 2, OwnerID: "owner2"}

	t.Run("should fulfill proposal successfully", func(t *testing.T) {
		mockRepo.EXPECT().GetProposal(10).Return(proposal).Once()
		mockRepo.EXPECT().GetBy([]int{1, 2}).Return([]*book.Book{book1, book2}, nil).Once()
		mockRepo.EXPECT().Upsert([]*book.Book{book2, book1}).Return(nil).Once()

		err := svc.FulfillProposal(10)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error if proposal not found", func(t *testing.T) {
		mockRepo.EXPECT().GetProposal(11).Return(nil).Once()

		err := svc.FulfillProposal(11)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "proposal not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error if repo.GetBy fails", func(t *testing.T) {
		mockRepo.EXPECT().GetProposal(12).Return(proposal).Once()
		mockRepo.EXPECT().GetBy([]int{1, 2}).Return(nil, errors.New("db error")).Once()

		err := svc.FulfillProposal(12)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error if not exactly 2 books", func(t *testing.T) {
		mockRepo.EXPECT().GetProposal(13).Return(proposal).Once()
		mockRepo.EXPECT().GetBy([]int{1, 2}).Return([]*book.Book{book1}, nil).Once()

		err := svc.FulfillProposal(13)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected 2 books in exchange")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error if upsert fails", func(t *testing.T) {
		mockRepo.EXPECT().GetProposal(14).Return(proposal).Once()
		mockRepo.EXPECT().GetBy([]int{1, 2}).Return([]*book.Book{book1, book2}, nil).Once()
		mockRepo.EXPECT().Upsert([]*book.Book{book2, book1}).Return(errors.New("upsert error")).Once()

		err := svc.FulfillProposal(14)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "upsert error")
		mockRepo.AssertExpectations(t)
	})
}
