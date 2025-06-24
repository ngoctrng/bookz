package book

type Book struct {
	OwnerID     string
	ISBN        string
	Title       string
	Description string
	BriefReview string
	Author      string
	Year        int
}

func New(ownerID, isbn, title, author string, year int) *Book {
	return &Book{
		OwnerID: ownerID,
		ISBN:    isbn,
		Title:   title,
		Author:  author,
		Year:    year,
	}
}

func (b *Book) AddDescription(description string) {
	b.Description = description
}

func (b *Book) AddBriefReview(review string) {
	b.BriefReview = review
}
