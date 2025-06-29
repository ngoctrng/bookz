package usecases

import (
	"github.com/google/uuid"
	"github.com/ngoctrng/bookz/internal/exchange"
)

type Repository interface {
	Save(p *exchange.Proposal) error
	GetByID(id int) (*exchange.Proposal, error)
	GetAll(uid uuid.UUID) ([]*exchange.Proposal, error)
	FetchRequestedBookOwner(id int) (uuid.UUID, error)
}

type MessageBus interface {
	PublishProposalAccepted(p *exchange.Proposal) error
}

type Usecase interface {
	CreateProposal(in CreateProposalInput) error
	GetProposalByID(id int) (*exchange.Proposal, error)
	GetAllProposals(uid uuid.UUID) ([]*exchange.Proposal, error)
	AcceptProposal(id int, uid uuid.UUID) error
}

type CreateProposalInput struct {
	Requested     exchange.BookID
	ForExchangeID exchange.BookID
	Message       string
	RequestBy     uuid.UUID
}
