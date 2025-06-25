package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/ngoctrng/bookz/internal/exchange"
)

type CreateProposalRequest struct {
	Requested     exchange.BookID `json:"requested" validate:"required"`
	ForExchangeID exchange.BookID `json:"for_exchange_id" validate:"required"`
	Message       string          `json:"message"`
}

func (r CreateProposalRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
