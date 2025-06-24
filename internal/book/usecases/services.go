package usecases

import (
	"errors"

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

func (s *Service) Get(isbn string) (*book.Book, error) {
	return s.repo.FindByISBN(isbn)
}

func (s *Service) Update(b *book.Book, ownerID string) error {
	existed, err := s.repo.FindByISBN(b.ISBN)
	if err != nil {
		return err
	}

	if existed.OwnerID != ownerID {
		return errors.New("you do not have permission to update this book")
	}

	return s.repo.Update(b)
}

func (s *Service) Delete(isbn, ownerID string) error {
	b, err := s.repo.FindByISBN(isbn)
	if err != nil {
		return err
	}

	if b.OwnerID != ownerID {
		return errors.New("you do not have permission to delete this book")
	}

	return s.repo.Delete(isbn)
}

func (s *Service) List() ([]*book.Book, error) {
	return s.repo.List()
}
