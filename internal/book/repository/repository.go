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

func (r *Repository) FindByID(id int) (*book.BookInfo, error) {
	var row BookInfoResult
	err := r.db.Table(tblBooks).
		Select("books.isbn, books.title, books.description, books.brief_review, books.author, books.year, books.owner_id, users.username").
		Joins("LEFT JOIN users ON books.owner_id = users.id").
		Where("books.id = ?", id).
		First(&row).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil // No book found with the given ID
	}
	if err != nil {
		return nil, err
	}

	return &book.BookInfo{
		ISBN:        row.ISBN,
		Title:       row.Title,
		Description: row.Description,
		BriefReview: row.BriefReview,
		Author:      row.Author,
		Year:        row.Year,
		Owner: book.BookOwner{
			ID:       row.OwnerID,
			Username: row.Username,
		},
	}, nil
}

func (r *Repository) Update(b *book.Book) error {
	schema := DomainToBookSchema(b)
	return r.db.Table(tblBooks).Where("id = ? AND owner_id = ?", b.ID, b.OwnerID).Updates(schema).Error
}

func (r *Repository) Delete(id int) error {
	return r.db.Table(tblBooks).Where("id = ?", id).Delete(nil).Error
}

func (r *Repository) List() ([]*book.BookInfo, error) {
	var rows []BookInfoResult
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
				ID:       row.OwnerID,
				Username: row.Username,
			},
		})
	}
	return infos, nil
}
