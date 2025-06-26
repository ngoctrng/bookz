package messagebus_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/hibiken/asynq"
	"github.com/ngoctrng/bookz/internal/events"
	"github.com/ngoctrng/bookz/internal/exchange"
	"github.com/ngoctrng/bookz/internal/exchange/messagebus"
	"github.com/ngoctrng/bookz/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublishProposalAccepted(t *testing.T) {
	client, inspector := testutil.CreateAsynqClient(t)

	bus := messagebus.New(client)

	t.Run("should publish proposal accepted event successfully", func(t *testing.T) {
		proposal := exchange.Proposal{ID: 123}

		err := bus.PublishProposalAccepted(&proposal)
		assert.NoError(t, err)

		time.Sleep(100 * time.Millisecond)

		pendingTasks, err := inspector.ListPendingTasks("default")
		require.NoError(t, err)
		assertSingleEvents(t, pendingTasks, proposal.ID)
	})

	t.Run("should handle multiple events", func(t *testing.T) {
		proposals := []exchange.Proposal{{ID: 1}, {ID: 2}, {ID: 3}}

		for _, p := range proposals {
			err := bus.PublishProposalAccepted(&p)
			assert.NoError(t, err)
		}

		time.Sleep(200 * time.Millisecond)

		pendingTasks, err := inspector.ListPendingTasks("default")
		require.NoError(t, err)
		assertMultipleEvents(t, pendingTasks, proposals)
	})
}

func assertMultipleEvents(t *testing.T, pendingTasks []*asynq.TaskInfo, proposals []exchange.Proposal) {
	eventCount := 0
	for _, task := range pendingTasks {
		if task.Type == events.ProposalAcceptedEvent {
			eventCount++
		}
	}

	assert.True(t, eventCount >= len(proposals),
		"Expected at least %d ProposalAcceptedEvent tasks, got %d", len(proposals), eventCount)
}

func assertSingleEvents(t *testing.T, pendingTasks []*asynq.TaskInfo, proposalID int) {
	assert.True(t, len(pendingTasks) > 0, "expected at least one pending task")

	var foundTask *asynq.TaskInfo
	for _, task := range pendingTasks {
		if task.Type == events.ProposalAcceptedEvent {
			foundTask = task
			break
		}
	}

	assert.NotNil(t, foundTask, "expected to find ProposalAcceptedEvent task")

	var payload events.ProposalAcceptedPayload
	err := json.Unmarshal(foundTask.Payload, &payload)
	assert.NoError(t, err)
	assert.Equal(t, proposalID, payload.ProposalID)
}
