package tasks_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/hibiken/asynq"
	"github.com/ngoctrng/bookz/internal/book/mocks"
	"github.com/ngoctrng/bookz/internal/book/tasks"
	"github.com/ngoctrng/bookz/internal/events"
	"github.com/stretchr/testify/assert"
)

func TestProposalAcceptedEventHandler(t *testing.T) {
	ctx := context.Background()
	uc := new(mocks.MockUsecase)
	h := tasks.NewHandler(uc)

	t.Run("should process valid event", func(t *testing.T) {
		payload := events.ProposalAcceptedPayload{ProposalID: 42}
		data, _ := json.Marshal(payload)
		task := asynq.NewTask(events.ProposalAcceptedEvent, data)

		uc.EXPECT().FulfillProposal(42).Return(nil).Once()
		err := h.ProposalAcceptedEventHandler(ctx, task)
		assert.NoError(t, err)
		uc.AssertExpectations(t)
	})

	t.Run("should return error for invalid task type", func(t *testing.T) {
		task := asynq.NewTask("other:event", []byte("{}"))
		err := h.ProposalAcceptedEventHandler(ctx, task)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid task type")
	})

	t.Run("should return error for invalid payload", func(t *testing.T) {
		task := asynq.NewTask(events.ProposalAcceptedEvent, []byte("not-json"))
		err := h.ProposalAcceptedEventHandler(ctx, task)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal task payload")
	})

	t.Run("should return error if FulfillProposal fails", func(t *testing.T) {
		payload := events.ProposalAcceptedPayload{ProposalID: 99}
		data, _ := json.Marshal(payload)
		task := asynq.NewTask(events.ProposalAcceptedEvent, data)

		uc.EXPECT().FulfillProposal(99).Return(errors.New("fail")).Once()
		err := h.ProposalAcceptedEventHandler(ctx, task)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fulfill proposal")
		uc.AssertExpectations(t)
	})
}
