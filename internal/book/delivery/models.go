package delivery

import "github.com/go-playground/validator/v10"

type BookRequest struct {
	ISBN        string `json:"isbn" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	BriefReview string `json:"brief_review"`
	Author      string `json:"author" validate:"required"`
	Year        int    `json:"year"`
}

func (r BookRequest) Validate() error {
	validator := validator.New()
	return validator.Struct(r)
}

type BookResponse struct {
	ISBN        string `json:"isbn"`
	Title       string `json:"title"`
	Description string `json:"description"`
	BriefReview string `json:"brief_review"`
	Author      string `json:"author"`
	Year        int    `json:"year"`
}
