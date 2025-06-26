package events

const ProposalAcceptedEvent = "proposal:accepted"

type ProposalAcceptedPayload struct {
	ProposalID int `json:"proposal_id"`
}
