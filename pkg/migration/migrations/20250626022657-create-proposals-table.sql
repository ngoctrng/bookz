
-- +migrate Up
CREATE TABLE proposals (
    id              SERIAL PRIMARY KEY,
    request_by      UUID        NOT NULL,
    request_to      UUID        NOT NULL,
    requested_id    INTEGER     NOT NULL, -- ID of the book wanting to be exchanged
    for_exchange_id INTEGER     NOT NULL, -- ID of the book being offered in exchange
    message         TEXT,
    status          VARCHAR(32) NOT NULL,
    requested_at    TIMESTAMPTZ NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_proposals_request_by ON proposals (request_by);

-- +migrate Down
DROP TABLE IF EXISTS proposals;