CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$ 
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'pager') THEN
		CREATE TYPE pager AS ENUM ('DISCORD', 'EMAIL', 'SMS', 'PHONE');
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'exchange') THEN
		CREATE TYPE exchange AS ENUM ('BINANCE', 'FTX', 'DERIBIT', 'BITFINEX');
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'dca_strategy') THEN
		CREATE TYPE dca_strategy AS ENUM ('CONSTANT', 'LINEAR', 'EXPONENTIAL');
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'exchange_execution_strategy') THEN
		CREATE TYPE exchange_execution_strategy AS ENUM ('PRIMARY_ONLY', 'ATTEMPT_ALL_REGISTERED')
	END IF;
END
$$;

CREATE TABLE IF NOT EXISTS s_account_accounts (
	-- use the discord id associated to the user here
	user_id VARCHAR(20) NOT NULL UNIQUE,
	username VARCHAR(50) NOT NULL UNIQUE,
	password VARCHAR(64) NOT NULL,

	email VARCHAR(50),
	phone_number VARCHAR(20),

	high_priority_pager pager NOT NULL DEFAULT 'DISCORD',
	low_priority_pager pager NOT NULL DEFAULT 'DISCORD',

	created TIME NOT NULL DEFAULT now(),
	updated TIME NOT NULL DEFAULT now(),
	last_payment_timestamp TIME NOT NULL DEFAULT now(),

	primary_exchange exchange NOT NULL DEFAULT 'BINANCE',
	default_exchange_execution_strategy exchange_execution_strategy NOT NULL DEFAULT 'PRIMARY_ONLY',

	is_admin BOOLEAN DEFAULT FALSE,
	is_futures_member BOOLEAN DEFAULT FALSE,

	default_dca_strategy dca_strategy NOT NULL DEFAULT 'LINEAR',

	PRIMARY KEY(user_id)
);

CREATE TABLE IF NOT EXISTS s_account_exchanges (
	exchange_id uuid DEFAULT uuid_generate_v4(),
	exchange exchange,

	user_id VARCHAR(20) NOT NULL,
	
	api_key VARCHAR(200) NOT NULL,
	secret_key VARCHAR(200) NOT NULL,
	subaccount VARCHAR(256) NOT NULL DEFAULT 'UNKNOWN',

	created TIME NOT NULL DEFAULT now(),
	updated TIME NOT NULL DEFAULT now(),

	is_active BOOLEAN DEFAULT FALSE,

	PRIMARY KEY(exchange_id),
	CONSTRAINT fk_account
		FOREIGN KEY(user_id)
			REFERENCES s_account_accounts(user_id) ON DELETE SET NULL,
	
	UNIQUE(user_id, exchange)
);
