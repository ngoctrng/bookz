package usecases

import "github.com/ngoctrng/bookz/internal/book"

type Repository interface {
	Save(b *book.Book) error
	FindByID(id int) (*book.BookInfo, error)
	Update(b *book.Book) error
	Delete(id int) error
	List() ([]*book.BookInfo, error)
	GetProposal(pid int) *book.ProposalDetails
	GetBy(ids []int) ([]*book.Book, error)
	Upsert(books []*book.Book) error
}

type Usecase interface {
	Create(b *book.Book) error
	Get(id int) (*book.BookInfo, error)
	Update(b *book.Book, ownerID string) error
	Delete(id int, ownerID string) error
	List() ([]*book.BookInfo, error)
	FulfillProposal(pid int) error
}
