package exchange

import (
	"time"

	"github.com/google/uuid"
)

type BookID int
type RequestStatus string

const (
	RequestStatusReviewing RequestStatus = "REVIEWING"
	RequestStatusRejected  RequestStatus = "REJECTED"
	RequestStatusAccepted  RequestStatus = "ACCEPTED"
)

type Proposal struct {
	ID            int
	RequestedID   BookID
	ForExchangeID BookID
	Message       string
	Status        RequestStatus
	RequestedAt   time.Time
	RequestBy     uuid.UUID
	RequestTo     uuid.UUID // Owner of the requested book
}

func OpenProposal(by uuid.UUID, requested BookID, forExchangeID BookID) *Proposal {
	return &Proposal{
		RequestedID:   requested,
		ForExchangeID: forExchangeID,
		Status:        RequestStatusReviewing,
		RequestedAt:   time.Now(),
		RequestBy:     by,
	}
}

func (p *Proposal) AddMessage(message string) {
	p.Message = message
}

func (p *Proposal) SendRequestTo(owner uuid.UUID) {
	p.RequestTo = owner
}

func (p *Proposal) Accept() {
	p.Status = RequestStatusAccepted
}

func (p *Proposal) Reject() {
	p.Status = RequestStatusRejected
}

func (p *Proposal) IsAccepted() bool {
	return p.Status == RequestStatusAccepted
}
