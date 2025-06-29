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
		Select("books.id, books.isbn, books.title, books.description, books.brief_review, books.author, books.year, books.owner_id, users.username").
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
		ID:          row.ID,
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
		Select("books.id, books.isbn, books.title, books.description, books.brief_review, books.author, books.year, books.owner_id, users.username").
		Joins("LEFT JOIN users ON books.owner_id = users.id").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	infos := make([]*book.BookInfo, 0, len(rows))
	for _, row := range rows {
		infos = append(infos, &book.BookInfo{
			ID:          row.ID,
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

func (r *Repository) GetProposal(pid int) *book.ProposalDetails {
	var proposal book.ProposalDetails
	err := r.db.Table("proposals").
		Select("id, requested_id, for_exchange_id").
		Where("id = ?", pid).
		First(&proposal).Error
	if err != nil {
		return nil
	}

	return &proposal
}

func (r *Repository) GetBy(ids []int) ([]*book.Book, error) {
	var schemas []BookSchema
	err := r.db.Table(tblBooks).
		Where("id IN ?", ids).
		Find(&schemas).Error
	if err != nil {
		return nil, err
	}

	books := make([]*book.Book, 0, len(schemas))
	for _, s := range schemas {
		books = append(books, BookSchemaToDomain(&s))
	}

	return books, nil
}

func (r *Repository) Upsert(books []*book.Book) error {
	r.db.Begin()
	defer r.db.Commit()

	for _, b := range books {
		schema := DomainToBookSchema(b)
		err := r.db.Table(tblBooks).
			Where("id = ?", b.ID).
			Updates(schema).Error
		if err != nil {
			r.db.Rollback()
			return err
		}
	}

	return nil
}
