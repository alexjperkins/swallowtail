CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$ 
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'pager') THEN
		CREATE TYPE pager AS ENUM ('DISCORD', 'EMAIL', 'SMS', 'PHONE');
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'venue') THEN
		CREATE TYPE venue AS ENUM ('BINANCE', 'FTX', 'DERIBIT', 'BITFINEX');
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'dca_strategy') THEN
		CREATE TYPE dca_strategy AS ENUM ('CONSTANT', 'LINEAR', 'EXPONENTIAL');
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 's_account_venue_account_type') THEN
		CREATE TYPE s_account_venue_account_type AS ENUM ('TRADING', 'TESTING', 'TREASURY');
	END IF;
END
$$;

CREATE TABLE IF NOT EXISTS s_account_accounts (
	-- use the discord id associated to the user here
	user_id VARCHAR(20) NOT NULL UNIQUE,
	username VARCHAR(50) NOT NULL UNIQUE,

	email VARCHAR(50),
	phone_number VARCHAR(20),

	high_priority_pager pager NOT NULL DEFAULT 'DISCORD',
	low_priority_pager pager NOT NULL DEFAULT 'DISCORD',

	created TIME NOT NULL DEFAULT now(),
	updated TIME NOT NULL DEFAULT now(),
	last_payment_timestamp TIME NOT NULL DEFAULT now(),

	primary_venue venue NOT NULL DEFAULT 'BINANCE',

	is_admin BOOLEAN DEFAULT FALSE,
	is_futures_member BOOLEAN DEFAULT FALSE,

	default_dca_strategy dca_strategy NOT NULL DEFAULT 'LINEAR',

	PRIMARY KEY(user_id)
);

CREATE TABLE IF NOT EXISTS s_account_venue_accounts (
	venue_account_id uuid DEFAULT uuid_generate_v4(),
	venue_id venue,

	user_id VARCHAR(20) NOT NULL,
	
	api_key VARCHAR(200) NOT NULL,
	secret_key VARCHAR(200) NOT NULL,
	subaccount VARCHAR(256) NOT NULL DEFAULT 'UNKNOWN',

	url VARCHAR(512),
	ws_url VARCHAR(512),

	account_alias VARCHAR(256) NOT NULL DEFAULT 'PRIMARY',

	created TIME NOT NULL DEFAULT now(),
	updated TIME NOT NULL DEFAULT now(),

	is_active BOOLEAN DEFAULT FALSE,

	PRIMARY KEY(venue_account_id),
	CONSTRAINT fk_account
		FOREIGN KEY(user_id)
			REFERENCES s_account_accounts(user_id) ON DELETE SET NULL,
	
	UNIQUE(user_id, venue_id, subaccount),
	UNIQUE(user_id, account_alias)
);

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
