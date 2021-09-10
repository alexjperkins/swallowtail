CREATE EXTENSION IF NOT EXISTS = "uuid-ossp"

DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 's_tradeengine_actor_type') THEN
		CREATE TYPE s_tradeengine_actor_type AS ENUM ('AUTOMATED', 'MANUAL', 'INTERNAL', 'EXTERNAL');
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 's_tradeengine_trade_type') THEN
		CREATE TYPE s_tradeengine_trade_type AS ENUM ('SPOT', 'FUTURESPERP', 'FUTURESQUARTERLY');
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'exchange') THEN
		CREATE TYPE exchange AS ENUM ('BINANCE', 'FTX');
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 's_tradeengine_asset_pair') THEN
		CREATE TYPE s_tradeengine_asset_pair AS ENUM ('USDT', 'BTC', 'USD');
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 's_tradeengine_trade_status') THEN
		CREATE TYPE s_tradeengine_trade_status AS ENUM ('PENDING', 'ACTIVE', 'COMPLETE', "CANCELLED");
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 's_tradeengine_trade_side') THEN
		CREATE TYPE s_tradeengine_trade_side AS ENUM ('BUY', 'SELL' );
	END IF;
END
$$;

CREATE TABLE IF NOT EXISTS s_tradeengine_trades (
	trade_id uuid DEFAULT uuid_generate_v4(),

	actor_id VARCHAR(32) NOT NULL,
	actor_type s_tradeengine_actor_type NOT NULL,

	idempotency_key VARCHAR(256) UNIQUE,

	trade_type s_tradeengine_trade_type NOT NULL,

	asset VARCHAR(8) NOT NULL,
	pair VARCHAR(4) NOT NULL,

	entry DECIMAL NOT NULL,
	stop_loss DECIMAL NOT NULL,
	take_profits []DECIMAL NOT NULL,

	status s_tradeengine_trade_status NOT NULL DEFAULT 'PENDING',
	risk_return DECIMAL,

	created TIME NOT NULL,
	last_updated TIME NOT NULL,
	PRIMARY KEY(trade_id)
)

CREATE TABLE IF NOT EXISTS s_tradeengine_trade_participants (
	trade_id VARCHAR(32),
	user_id VARCHAR(20),
	
	is_bot BOOLEAN NOT NULL DEFAULT FALSE,

	size DECIMAL NOT NULL,

	exchange exchange NOT NULL,

	executed TIME NOT NULL,

	CONSTRAINT fk_tradeengine_
		FOREIGN KEY(trade_id)
			REFERENCES s_tradeengine_trades(trade_id) ON DELETE SET NULL,
	
	UNIQUE(trade_id, user_id)
)
