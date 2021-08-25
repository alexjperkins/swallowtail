CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$ 
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'trade_type') THEN
		CREATE TYPE trade_type AS ENUM ('SPOT', 'PERPETUALS', 'FUTURES_QUARTERLY');
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'trade_side') THEN
		CREATE TYPE trade_side AS ENUM ('BUY', 'SELL');
	END IF;
END
$$;

CREATE TABLE IF NOT EXISTS s_binance_trades (
	trade_id uuid DEFAULT uuid_generate_v4(),
	idempotency_key VARCHAR(255) UNIQUE NOT NULL,
	user_discord_id VARCHAR(20) NOT NULL,

	
	side trade_side NOT NULL,
	type trade_type NOT NULL,
	asset_pair VARCHAR(15) NOT NULL,
	amount VARCHAR(15) NOT NULL,
	value VARCHAR(15) NOT NULL,

	received TIME NOT NULL,
	attemped TIME NOT NULL,
	attempted_retry_until TIME NOT NULL,

	PRIMARY KEY(trade_id)
);

CREATE INDEX IF NOT EXISTS idx_trades_idempotency_key
ON s_binance_trades(idempotency_key);
