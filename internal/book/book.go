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

type BookInfo struct {
	ISBN        string    `json:"isbn"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	BriefReview string    `json:"brief_review"`
	Author      string    `json:"author"`
	Year        int       `json:"year"`
	Owner       BookOwner `json:"owner"`
}

type BookOwner struct {
	OwnerID  string `json:"owner_id"`
	Username string `json:"username"`
}
