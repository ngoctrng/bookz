package usecases

import (
	"errors"
	"log/slog"

	"github.com/ngoctrng/bookz/internal/book"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Create(b *book.Book) error {
	return s.repo.Save(b)
}

func (s *Service) Get(id int) (*book.BookInfo, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Update(b *book.Book, ownerID string) error {
	existed, err := s.repo.FindByID(b.ID)
	if err != nil {
		return err
	}

	if existed.Owner.ID != ownerID {
		return errors.New("you do not have permission to update this book")
	}

	return s.repo.Update(b)
}

func (s *Service) Delete(id int, ownerID string) error {
	b, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if b.Owner.ID != ownerID {
		return errors.New("you do not have permission to delete this book")
	}

	return s.repo.Delete(id)
}

func (s *Service) List() ([]*book.BookInfo, error) {
	return s.repo.List()
}

func (s *Service) FulfillProposal(pid int) error {
	p := s.repo.GetProposal(pid)
	if p == nil {
		return errors.New("proposal not found")
	}

	books, err := s.repo.GetBy([]int{p.RequestedID, p.ForExchangeID})
	if err != nil {
		return err
	}
	if len(books) != 2 {
		return errors.New("expected 2 books in exchange")
	}

	var requestedBook, exchangeBook *book.Book
	for _, b := range books {
		switch b.ID {
		case p.RequestedID:
			requestedBook = b
		case p.ForExchangeID:
			exchangeBook = b
		}
	}

	slog.Info("books", "requested", requestedBook, "exchange", exchangeBook, "books", books)

	book.ChangeOwnerFor(exchangeBook, requestedBook)
	if err := s.repo.Upsert([]*book.Book{exchangeBook, requestedBook}); err != nil {
		return err
	}

	return nil
}
