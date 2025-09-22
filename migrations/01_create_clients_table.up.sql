CREATE TABLE clients (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    residence_address TEXT,
    birth_date DATE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);