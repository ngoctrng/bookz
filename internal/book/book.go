package book

type Book struct {
	ID          int
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

func (b *Book) ChangeOwner(newOwnerID string) {
	b.OwnerID = newOwnerID
}

func ChangeOwnerFor(b1 *Book, b2 *Book) (string, string) {
	b1Owner := b1.OwnerID
	b2Owner := b2.OwnerID
	b1.ChangeOwner(b2Owner)
	b2.ChangeOwner(b1Owner)
	return b1Owner, b2Owner
}

type BookInfo struct {
	ID          int       `json:"id"`
	ISBN        string    `json:"isbn"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	BriefReview string    `json:"brief_review"`
	Author      string    `json:"author"`
	Year        int       `json:"year"`
	Owner       BookOwner `json:"owner"`
}

type BookOwner struct {
	ID       string `json:"owner_id"`
	Username string `json:"username"`
}

type ProposalDetails struct {
	ID            int `json:"id"`
	RequestedID   int `json:"requested_id"`    // ID of the book wanting
	ForExchangeID int `json:"for_exchange_id"` // ID of the book being offered in exchange
}
