CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    account_number VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    nik VARCHAR(20) NOT NULL UNIQUE,
    phone_number VARCHAR(15) NOT NULL UNIQUE,
    balance DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for frequently queried fields
CREATE INDEX IF NOT EXISTS idx_account_number ON accounts(account_number);
CREATE INDEX IF NOT EXISTS idx_nik ON accounts(nik);
CREATE INDEX IF NOT EXISTS idx_phone_number ON accounts(phone_number);

-- Add trigger to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_accounts_updated_at
    BEFORE UPDATE ON accounts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();