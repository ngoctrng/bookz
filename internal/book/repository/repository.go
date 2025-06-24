package repository

import (
	"github.com/ngoctrng/bookz/internal/book"
	"gorm.io/gorm"
)

const tblBooks = "books"

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(b *book.Book) error {
	schema := DomainToBookSchema(b)
	return r.db.Table(tblBooks).Create(schema).Error
}

func (r *Repository) FindByISBN(isbn string) (*book.Book, error) {
	var schema BookSchema
	err := r.db.Table(tblBooks).Where("isbn = ?", isbn).First(&schema).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return BookSchemaToDomain(&schema), nil
}

func (r *Repository) Update(b *book.Book) error {
	schema := DomainToBookSchema(b)
	return r.db.Table(tblBooks).Where("isbn = ? AND owner_id = ?", b.ISBN, b.OwnerID).Updates(schema).Error
}

func (r *Repository) Delete(isbn string) error {
	return r.db.Table(tblBooks).Where("isbn = ?", isbn).Delete(nil).Error
}

func (r *Repository) List() ([]*book.Book, error) {
	var schemas []BookSchema
	if err := r.db.Table(tblBooks).Find(&schemas).Error; err != nil {
		return nil, err
	}
	var books []*book.Book
	for _, s := range schemas {
		books = append(books, BookSchemaToDomain(&s))
	}
	return books, nil
}
