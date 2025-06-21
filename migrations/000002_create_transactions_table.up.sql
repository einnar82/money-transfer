CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    source_account_id INT NOT NULL REFERENCES accounts(id),
    destination_account_id INT NOT NULL REFERENCES accounts(id),
    amount DECIMAL NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);