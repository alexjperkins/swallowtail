-- Rename table
ALTER TABLE IF EXISTS s_account_exchanges RENAME s_account_venue_accounts;

-- Drop constraints
ALTER TABLE IF EXISTS s_account_venue_accounts DROP CONSTRAINT s_account_venue_accounts_user_id_exchange_key;
ALTER TABLE IF EXISTS s_account_venue_accounts DROP CONSTRAINT s_account_venue_accounts_pkey;

-- Rename columns
ALTER TABLE IF EXISTS s_account_venue_accounts RENAME COLUMN exchange_id TO venue_account_id;
ALTER TABLE IF EXISTS s_account_venue_accounts RENAME COLUMN exchange TO venue_id;

-- Add constraints
ALTER TABLE IF EXISTS s_account_venue_accounts ADD PRIMARY KEY (venue_account_id);
ALTER TABLE IF EXISTS s_account_venue_accounts ADD UNIQUE (user_id, venue_id, subaccount);

-- Add new account_alias column & constraint
ALTER TABLE IF EXISTS ADD COLUMN account_alias VARCHAR(256);
ALTER TABLE IF EXISTS s_account_venue_accounts ADD UNIQUE (user_id, account_alias);

-- TODO drop primary exchange, enum & move to UUID on the s_account_venue_accounts table.

-- Drop columns
ALTER TABLE IF EXISTS s_account_accounts DROP COLUMN primary_exchange;
ALTER TABLE IF EXISTS s_account_accounts DROP COLUMN default_exchange_execution_strategy;

-- Rename enum type
ALTER TYPE IF EXISTS exchange RENAME to venue CASCADE;

ALTER TABLE IF EXISTS s_account_venue_accounts ADD COLUMN primary_venue venue DEFAULT = 'BINANCE';

DROP TYPE IF EXISTS exchange_execution_strategy;
