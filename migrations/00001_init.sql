-- +goose Up
CREATE TABLE users (
    id         BIGSERIAL   PRIMARY KEY,
    email      TEXT        NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE categories (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name       TEXT        NOT NULL,
    type       TEXT        NOT NULL CHECK (type IN ('income', 'expense')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    is_deleted BOOLEAN     NOT NULL GENERATED ALWAYS AS (deleted_at IS NOT NULL) STORED
);

CREATE UNIQUE INDEX uq_categories_active_name
    ON categories (user_id, name, type)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_categories_user_id ON categories (user_id);

CREATE TABLE transactions (
    id           BIGSERIAL   PRIMARY KEY,
    user_id      BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id  BIGINT      NOT NULL REFERENCES categories(id),
    amount_cents BIGINT      NOT NULL CHECK (amount_cents > 0),
    currency     VARCHAR(3)  NOT NULL DEFAULT 'RUB',
    comment      TEXT,
    occurred_on  DATE        NOT NULL DEFAULT CURRENT_DATE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_transactions_user_occurred ON transactions (user_id, occurred_on);

-- +goose Down
DROP TABLE transactions;
DROP TABLE categories;
DROP TABLE users;
