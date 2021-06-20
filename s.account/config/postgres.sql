CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$ 
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'pager') THEN
		CREATE TYPE pager AS ENUM ('discord', 'email', 'sms', 'phone', 'unknown');
	END IF;
END
$$;

DO $$ 
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'exchange') THEN
		CREATE TYPE exchange AS ENUM ('binance', 'ftx');
	END IF;
END
$$;

CREATE TABLE IF NOT EXISTS googlesheets (
	googlesheets_id uuid DEFAULT uuid_generate_v4(),
	spreadsheet_id VARCHAR(200) NOT NULL UNIQUE,
	sheet_id VARCHAR(200) NOT NULL UNIQUE,
	account_id uuid,

	PRIMARY KEY(googlesheets_id),
	CONSTRAINT fk_account
		FOREIGN KEY(account_id)
			REFERENCES accounts(account_id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS exchanges (
	exchange_id uuid DEFAULT uuid_generate_v4(),
	exchange exchange,
	api_key VARCHAR(200),
	secret_key VARCHAR(200),
	account_id uuid,

	PRIMARY KEY(exchange_id),
	CONSTRAINT fk_account
		FOREIGN KEY(account_id)
			REFERENCES accounts(account_id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS accounts (
	account_id uuid DEFAULT uuid_generate_v4(),
	username VARCHAR(50) NOT NULL UNIQUE,
	password VARCHAR(50) NOT NULL,
	email VARCHAR(50),
	discord_id VARCHAR(20),
	phone_number VARCHAR(20),

	high_priority_pager pager NOT NULL DEFAULT 'unknown',
	low_priority_pager pager NOT NULL DEFAULT 'unknown',

	PRIMARY KEY(account_id)
);
