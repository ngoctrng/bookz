
-- +migrate Up
CREATE TABLE IF NOT EXISTS books (
    id            SERIAL PRIMARY KEY,
    isbn          VARCHAR(64) NOT NULL,
    owner_id      uuid NOT NULL,
    title         VARCHAR(255) NOT NULL,
    description   TEXT,
    brief_review  TEXT,
    author        VARCHAR(128) NOT NULL,
    year          INT,
    created_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE IF EXISTS books;