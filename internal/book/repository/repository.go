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

func (r *Repository) List() ([]*book.BookInfo, error) {
	type result struct {
		ISBN        string
		Title       string
		Description string
		BriefReview string
		Author      string
		Year        int
		OwnerID     string
		Username    string
	}

	var rows []result
	err := r.db.Table(tblBooks).
		Select("books.isbn, books.title, books.description, books.brief_review, books.author, books.year, books.owner_id, users.username").
		Joins("LEFT JOIN users ON books.owner_id = users.id").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	infos := make([]*book.BookInfo, 0, len(rows))
	for _, row := range rows {
		infos = append(infos, &book.BookInfo{
			ISBN:        row.ISBN,
			Title:       row.Title,
			Description: row.Description,
			BriefReview: row.BriefReview,
			Author:      row.Author,
			Year:        row.Year,
			Owner: book.BookOwner{
				OwnerID:  row.OwnerID,
				Username: row.Username,
			},
		})
	}
	return infos, nil
}
