INSERT INTO accounts (client_id, account_number, balance, created_at)
SELECT c.id,
       LPAD(CAST(FLOOR(RANDOM()*10000000000000000)::bigint AS TEXT), 16, '0') AS account_number,
       0.00,
       NOW()
FROM clients c
WHERE NOT EXISTS (
  SELECT 1 FROM accounts a WHERE a.client_id = c.id
);