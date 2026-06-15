ALTER TABLE accounts
DROP CONSTRAINT IF EXISTS accounts_account_number_key;

DROP SEQUENCE IF EXISTS account_number_seq;