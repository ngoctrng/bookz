package tasks

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/hibiken/asynq"
	"github.com/ngoctrng/bookz/internal/book/usecases"
	"github.com/ngoctrng/bookz/internal/events"
)

type Handler struct {
	uc usecases.Usecase
}

func NewHandler(uc usecases.Usecase) *Handler {
	return &Handler{
		uc: uc,
	}
}

func (h *Handler) ProposalAcceptedEventHandler(ctx context.Context, t *asynq.Task) error {
	if t.Type() != events.ProposalAcceptedEvent {
		return errors.New("invalid task type")
	}

	var p events.ProposalAcceptedPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return errors.Join(err, errors.New("failed to unmarshal task payload"))
	}

	if err := h.uc.FulfillProposal(p.ProposalID); err != nil {
		return errors.Join(err, errors.New("failed to fulfill proposal"))
	}

	return nil
}
