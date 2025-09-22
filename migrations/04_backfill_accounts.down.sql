DELETE FROM accounts a
USING clients c
WHERE a.client_id = c.id
  AND a.balance = 0.00;