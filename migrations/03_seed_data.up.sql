WITH inserted AS (
  INSERT INTO clients (first_name, last_name, email, residence_address, birth_date)
  VALUES
    ('Ana',   'Ivić',      'ana@example.com',   'Beograd',   '1995-05-01'),
    ('Marko', 'Jovanović', 'marko@example.com', 'Novi Sad',  '1990-02-11')
  ON CONFLICT (email) DO NOTHING
  RETURNING id, email
)
INSERT INTO accounts (client_id, account_number, balance)
SELECT c.id, v.account_number, v.balance
FROM (
  VALUES
    ('ana@example.com',   '1234567890123456', 1000.00),
    ('marko@example.com', '2345678901234567',  250.50)
) AS v(email, account_number, balance)
JOIN (
  SELECT id, email FROM clients
) c ON c.email = v.email;