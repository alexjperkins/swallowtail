-- Rename table.
ALTER TABLE IF EXISTS s_account_exchanges RENAME TO s_account_venue_accounts;

-- Drop constraints.
ALTER TABLE IF EXISTS s_account_venue_accounts DROP CONSTRAINT s_account_exchanges_user_id_exchange_key;
ALTER TABLE IF EXISTS s_account_venue_accounts DROP CONSTRAINT s_account_exchanges_pkey;

-- Rename columns.
ALTER TABLE IF EXISTS s_account_venue_accounts RENAME COLUMN exchange_id TO venue_account_id;
ALTER TABLE IF EXISTS s_account_venue_accounts RENAME COLUMN exchange TO venue_id;

-- Add constraints.
ALTER TABLE IF EXISTS s_account_venue_accounts ADD PRIMARY KEY (venue_account_id);
ALTER TABLE IF EXISTS s_account_venue_accounts ADD UNIQUE (user_id, venue_id, subaccount);

-- Add new account_alias column & constraint.
ALTER TABLE IF EXISTS s_account_venue_accounts ADD COLUMN account_alias VARCHAR(256) DEFAULT 'PRIMARY';
ALTER TABLE IF EXISTS s_account_venue_accounts ADD UNIQUE (user_id, account_alias);

-- Drop columns.
ALTER TABLE IF EXISTS s_account_accounts DROP COLUMN primary_exchange;

-- Rename enum type.
ALTER TYPE exchange RENAME TO venue;

ALTER TABLE IF EXISTS s_account_accounts ADD COLUMN primary_venue venue DEFAULT 'BINANCE';

-- Drop unused field.
DROP TYPE IF EXISTS exchange_execution_strategy;

-- Add venue urls.
ALTER TABLE IF EXISTS s_account_venue_accounts ADD COLUMN url VARCHAR(512);
ALTER TABLE IF EXISTS s_account_venue_accounts ADD COLUMN ws_url VARCHAR(512);

-- Add venue account type.
DO $$ 
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 's_account_venue_account_type') THEN
		CREATE TYPE s_account_venue_account_type AS ENUM ('TRADING', 'TESTING', 'TREASURY');
	END IF;
END
$$;

-- Add internal account table.
CREATE TABLE IF NOT EXISTS s_account_internal_venue_accounts (
	internal_account_id uuid DEFAULT uuid_generate_v4(),
	venue_id venue,

	api_key VARCHAR(200) NOT NULL,
	secret_key VARCHAR(200) NOT NULL,
	subaccount VARCHAR(256) NOT NULL DEFAULT 'UNKNOWN',

	url VARCHAR(512),
	ws_url VARCHAR(512),

	venue_account_type s_account_venue_account_type NOT NULL DEFAULT 'TESTING',

	created TIME NOT NULL DEFAULT now(),
	updated TIME NOT NULL DEFAULT now(),

	PRIMARY KEY(internal_account_id),

	UNIQUE(venue_id, subaccount, venue_account_type)
);
