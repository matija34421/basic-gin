CREATE TABLE IF NOT EXISTS transactions (
  id              SERIAL PRIMARY KEY,
  from_account_id INT NOT NULL REFERENCES accounts(id),
  to_account_id   INT NOT NULL REFERENCES accounts(id),
  amount          NUMERIC(18,2) NOT NULL CHECK (amount > 0),
  created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_transactions_from ON transactions(from_account_id);
CREATE INDEX IF NOT EXISTS idx_transactions_to   ON transactions(to_account_id);
