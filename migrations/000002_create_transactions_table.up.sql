CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT REFERENCES accounts(id),
    transaction_ref VARCHAR(50) UNIQUE NOT NULL, -- ฟอร์แมต TXN<YYYYMMDD><RUNNING_NUMBER>
    transaction_type VARCHAR(20) NOT NULL,       -- ALLOWED: 'DEPOSIT', 'WITHDRAW'
    amount DECIMAL(15,2) NOT NULL,
    balance_before DECIMAL(15,2) NOT NULL,
    balance_after DECIMAL(15,2) NOT NULL,
    description VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);



