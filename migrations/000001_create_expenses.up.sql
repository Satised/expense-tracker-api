CREATE TABLE expenses (
    id          BIGSERIAL PRIMARY KEY,
    amount      NUMERIC(14, 2) NOT NULL CHECK (amount > 0),
    category    TEXT NOT NULL,
    note        TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_expenses_created_at ON expenses (created_at);