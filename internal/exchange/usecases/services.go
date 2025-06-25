package usecases

import (
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
