package repository

import (
	"time"

	"github.com/ngoctrng/bookz/internal/book"
)

type BookSchema struct {
	ISBN        string `gorm:"primaryKey"`
	OwnerID     string
	Title       string
	Description string
	BriefReview string
	Author      string
	Year        int
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func DomainToBookSchema(b *book.Book) *BookSchema {
	return &BookSchema{
		ISBN:        b.ISBN,
		OwnerID:     b.OwnerID,
		Title:       b.Title,
		Description: b.Description,
		BriefReview: b.BriefReview,
		Author:      b.Author,
		Year:        b.Year,
	}
}

func BookSchemaToDomain(s *BookSchema) *book.Book {
	return &book.Book{
		ISBN:        s.ISBN,
		OwnerID:     s.OwnerID,
		Title:       s.Title,
		Description: s.Description,
		BriefReview: s.BriefReview,
		Author:      s.Author,
		Year:        s.Year,
	}
}

type BookInfoResult struct {
	ISBN        string
	Title       string
	Description string
	BriefReview string
	Author      string
	Year        int
	OwnerID     string
	Username    string
}
