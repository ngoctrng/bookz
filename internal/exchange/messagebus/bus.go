package messagebus

import (
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/ngoctrng/bookz/internal/events"
	"github.com/ngoctrng/bookz/internal/exchange"
)

type MessageBus struct {
	client *asynq.Client
}

func New(client *asynq.Client) *MessageBus {
	return &MessageBus{client: client}
}

func (b *MessageBus) PublishProposalAccepted(p *exchange.Proposal) error {
	payload := events.ProposalAcceptedPayload{
		ProposalID: p.ID,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(events.ProposalAcceptedEvent, data)
	if _, err := b.client.Enqueue(task); err != nil {
		return err
	}

	return nil
}
