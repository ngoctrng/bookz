package usecases_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ngoctrng/bookz/internal/exchange"
	"github.com/ngoctrng/bookz/internal/exchange/mocks"
	"github.com/ngoctrng/bookz/internal/exchange/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateProposal(t *testing.T) {
	r := new(mocks.MockRepository)
	bus := new(mocks.MockMessageBus)
	svc := usecases.NewProposalService(r, bus)

	input := usecases.CreateProposalInput{
		RequestBy:     uuid.New(),
		Requested:     exchange.BookID(1),
		ForExchangeID: 42,
		Message:       "Let's trade!",
	}
	proposal := exchange.OpenProposal(input.RequestBy, input.Requested, input.ForExchangeID)
	proposal.AddMessage(input.Message)

	t.Run("should create proposal successfully", func(t *testing.T) {
		r.EXPECT().Save(mock.AnythingOfType("*exchange.Proposal")).Return(nil).Once()
		r.EXPECT().FetchRequestedBookOwner(int(input.Requested)).Return(uuid.New(), nil).Once()
		err := svc.CreateProposal(input)
		assert.NoError(t, err)
		r.AssertExpectations(t)
	})

	t.Run("should return error if repo fails", func(t *testing.T) {
		r.EXPECT().Save(mock.AnythingOfType("*exchange.Proposal")).Return(errors.New("db error")).Once()
		r.EXPECT().FetchRequestedBookOwner(int(input.Requested)).Return(uuid.New(), nil).Once()
		err := svc.CreateProposal(input)
		assert.Error(t, err)
		r.AssertExpectations(t)
	})
}

func TestGetProposalByID(t *testing.T) {
	r := new(mocks.MockRepository)
	bus := new(mocks.MockMessageBus)
	svc := usecases.NewProposalService(r, bus)

	proposal := &exchange.Proposal{ID: 1}
	t.Run("should get proposal by id", func(t *testing.T) {
		r.EXPECT().GetByID(1).Return(proposal, nil).Once()
		got, err := svc.GetProposalByID(1)
		assert.NoError(t, err)
		assert.Equal(t, proposal, got)
		r.AssertExpectations(t)
	})

	t.Run("should return error if repo fails", func(t *testing.T) {
		r.EXPECT().GetByID(1).Return(nil, errors.New("db error")).Once()
		got, err := svc.GetProposalByID(1)
		assert.Error(t, err)
		assert.Nil(t, got)
		r.AssertExpectations(t)
	})
}

func TestGetAllProposals(t *testing.T) {
	r := new(mocks.MockRepository)
	bus := new(mocks.MockMessageBus)
	svc := usecases.NewProposalService(r, bus)

	uid := uuid.New()
	proposals := []*exchange.Proposal{
		{ID: 1}, {ID: 2},
	}

	t.Run("should get all proposals for user", func(t *testing.T) {
		r.EXPECT().GetAll(uid).Return(proposals, nil).Once()
		got, err := svc.GetAllProposals(uid)
		assert.NoError(t, err)
		assert.Equal(t, proposals, got)
		r.AssertExpectations(t)
	})

	t.Run("should return error if repo fails", func(t *testing.T) {
		r.EXPECT().GetAll(uid).Return(nil, errors.New("db error")).Once()
		got, err := svc.GetAllProposals(uid)
		assert.Error(t, err)
		assert.Nil(t, got)
		r.AssertExpectations(t)
	})
}

func TestAcceptProposal(t *testing.T) {
	r := new(mocks.MockRepository)
	bus := new(mocks.MockMessageBus)
	svc := usecases.NewProposalService(r, bus)

	ownerID := uuid.New()
	requesterID := uuid.New()

	proposal := &exchange.Proposal{
		ID:        1,
		RequestBy: requesterID,
		RequestTo: ownerID,
		Status:    exchange.RequestStatusReviewing,
	}

	acceptedProposal := &exchange.Proposal{
		ID:        1,
		RequestBy: requesterID,
		RequestTo: ownerID,
		Status:    exchange.RequestStatusAccepted,
	}

	t.Run("should accept proposal successfully", func(t *testing.T) {
		r.EXPECT().GetByID(1).Return(proposal, nil).Once()
		r.EXPECT().Save(acceptedProposal).Return(nil).Once()
		bus.EXPECT().PublishProposalAccepted(acceptedProposal).Return(nil).Once()

		err := svc.AcceptProposal(1, ownerID)
		assert.NoError(t, err)
		r.AssertExpectations(t)
		bus.AssertExpectations(t)
	})

	t.Run("should return error if proposal not found", func(t *testing.T) {
		r.EXPECT().GetByID(99).Return(nil, nil).Once()

		err := svc.AcceptProposal(99, ownerID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
		r.AssertExpectations(t)
	})

	t.Run("should return error if repository fails", func(t *testing.T) {
		r.EXPECT().GetByID(1).Return(nil, errors.New("db error")).Once()

		err := svc.AcceptProposal(1, ownerID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
		r.AssertExpectations(t)
	})

	t.Run("should return error if unauthorized", func(t *testing.T) {
		wrongOwnerID := uuid.New()
		r.EXPECT().GetByID(1).Return(proposal, nil).Once()

		err := svc.AcceptProposal(1, wrongOwnerID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
		r.AssertExpectations(t)
	})

	t.Run("should return error if save fails", func(t *testing.T) {
		r.EXPECT().GetByID(1).Return(proposal, nil).Once()
		r.EXPECT().Save(acceptedProposal).Return(errors.New("save error")).Once()

		err := svc.AcceptProposal(1, ownerID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "save error")
		r.AssertExpectations(t)
	})

	t.Run("should return error if message bus fails", func(t *testing.T) {
		r.EXPECT().GetByID(1).Return(proposal, nil).Once()
		r.EXPECT().Save(acceptedProposal).Return(nil).Once()
		bus.EXPECT().PublishProposalAccepted(acceptedProposal).Return(errors.New("bus error")).Once()

		err := svc.AcceptProposal(1, ownerID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "bus error")
		r.AssertExpectations(t)
		bus.AssertExpectations(t)
	})
}
