-- Drop constraints
ALTER TABLE IF EXISTS s_account_venue_accounts DROP CONSTRAINT s_account_venue_accounts_user_id_exchange_key;
ALTER TABLE IF EXISTS s_account_venue_accounts DROP CONSTRAINT s_account_venue_accounts_pkey;

-- Add primary key
ALTER TABLE IF EXISTS s_account_venue_accounts RENAME COLUMN exchange_id TO venue_account_id;

-- Rename columns
ALTER TABLE IF EXISTS s_account_venue_accounts RENAME COLUMN exchange_id TO venue_id;

-- Add constraints
ALTER TABLE IF EXISTS s_account_venue_accounts ADD PRIMARY KEY (venue_account_id);
ALTER TABLE IF EXISTS s_account_venue_accounts ADD UNIQUE (user_id, venue_id, subaccount);

-- Add new account_alias column & constraint
ALTER TABLE IF EXISTS ADD COLUMN account_alias VARCHAR(256);
ALTER TABLE IF EXISTS s_account_venue_accounts ADD UNIQUE (user_id, account_alias);

-- Rename table
ALTER TABLE IF EXISTS s_account_exchanges RENAME s_account_venue_accounts;

-- TODO drop primary exchange, enum & move to UUID on the s_account_venue_accounts table.
