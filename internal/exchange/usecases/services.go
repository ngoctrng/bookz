package usecases

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ngoctrng/bookz/internal/exchange"
)

type Service struct {
	repo Repository
}

func NewProposalService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateProposal(in CreateProposalInput) error {
	p := exchange.OpenProposal(in.RequestBy, in.Requested, in.ForExchangeID)
	if in.Message != "" {
		p.AddMessage(in.Message)
	}

	ownerID, err := s.repo.FetchRequestedBookOwner(int(in.Requested))
	if err != nil {
		return err
	}
	p.SendRequestTo(ownerID)

	if err := s.repo.Save(p); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetProposalByID(id int) (*exchange.Proposal, error) {
	proposal, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return proposal, nil
}

func (s *Service) GetAllProposals(uid uuid.UUID) ([]*exchange.Proposal, error) {
	proposals, err := s.repo.GetAll(uid)
	if err != nil {
		return nil, err
	}

	return proposals, nil
}

func (s *Service) AcceptProposal(id int, uid uuid.UUID) error {
	proposal, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if proposal == nil {
		return errors.New("proposal not found")
	}

	if proposal.RequestTo != uid {
		return errors.New("unauthorized: only the book owner can accept this proposal")
	}

	proposal.Accept()

	return s.repo.Save(proposal)
}
