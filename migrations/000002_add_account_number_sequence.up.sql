-- Create sequence for account numbers
CREATE SEQUENCE IF NOT EXISTS account_number_seq
START WITH 0000000001
INCREMENT BY 1;

-- Ensure account_number is unique