package exchange_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ngoctrng/bookz/internal/exchange"
	"github.com/stretchr/testify/assert"
)

func TestOpenAProposal(t *testing.T) {
	requested, forExchange := exchange.BookID(1), exchange.BookID(2)
	uid := uuid.New()

	p := exchange.OpenProposal(uid, requested, forExchange)

	assert.NotNil(t, p)
	assert.Equal(t, exchange.RequestStatusReviewing, p.Status)
	assert.Equal(t, uid, p.RequestBy)
	assert.Equal(t, requested, p.RequestedID)
	assert.Equal(t, forExchange, p.ForExchangeID)
}

func TestOpenAProposalWithMessage(t *testing.T) {
	requested, forExchange := exchange.BookID(1), exchange.BookID(2)
	uid := uuid.New()

	p := exchange.OpenProposal(uid, requested, forExchange)
	p.AddMessage("Hello, I would like to exchange this book.")

	assert.NotNil(t, p)
	assert.Equal(t, "Hello, I would like to exchange this book.", p.Message)

}
