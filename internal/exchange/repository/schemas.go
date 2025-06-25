package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/ngoctrng/bookz/internal/exchange"
)

// ProposalSchema maps to the proposals table in the database.
type ProposalSchema struct {
	ID            int
	RequestBy     uuid.UUID
	RequestedID   int
	ForExchangeID int
	Message       string
	Status        string
	RequestedAt   time.Time
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

func DomainToProposalSchema(p *exchange.Proposal) *ProposalSchema {
	return &ProposalSchema{
		ID:            p.ID,
		RequestBy:     p.RequestBy,
		RequestedID:   int(p.RequestedID),
		ForExchangeID: int(p.ForExchangeID),
		Message:       p.Message,
		Status:        string(p.Status),
		RequestedAt:   p.RequestedAt,
	}
}

func ProposalSchemaToDomain(s *ProposalSchema) *exchange.Proposal {
	return &exchange.Proposal{
		ID:            s.ID,
		RequestBy:     s.RequestBy,
		RequestedID:   exchange.BookID(s.RequestedID),
		ForExchangeID: exchange.BookID(s.ForExchangeID),
		Message:       s.Message,
		Status:        exchange.RequestStatus(s.Status),
		RequestedAt:   s.RequestedAt,
	}
}
