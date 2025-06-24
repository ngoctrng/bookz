package usecases

import "github.com/ngoctrng/bookz/internal/book"

type Repository interface {
	Save(b *book.Book) error
	FindByISBN(isbn string) (*book.Book, error)
	Update(b *book.Book) error
	Delete(isbn string) error
	List() ([]*book.Book, error)
}

type Usecase interface {
	Create(b *book.Book) error
	Get(isbn string) (*book.Book, error)
	Update(b *book.Book, ownerID string) error
	Delete(isbn, ownerID string) error
	List() ([]*book.Book, error)
}
