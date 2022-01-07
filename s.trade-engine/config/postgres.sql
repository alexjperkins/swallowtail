CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 's_tradeengine_venue') THEN
		CREATE TYPE s_tradeengine_venue AS ENUM ('BINANCE', 'FTX', 'DERIBIT', 'BITFINEX');
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 's_tradeengine_actor_type') THEN
		CREATE TYPE s_tradeengine_actor_type AS ENUM ('AUTOMATED', 'MANUAL', 'INTERNAL', 'EXTERNAL');
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 's_tradeengine_trade_side') THEN
		CREATE TYPE s_tradeengine_trade_side AS ENUM ('BUY', 'SELL', 'LONG', 'SHORT');
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 's_tradeengine_order_status') THEN
		CREATE TYPE s_tradeengine_order_status AS ENUM (
			'PENDING_NEW_ORDER',
			'NEW_ORDER',
			'PENDING_CANCEL_ORDER',
			'CANCELLED_ORDER',
			'PARTIALLY_FILLED_ORDER',
			'FILLED_ORDER',
			'REJECTED',
			'EXPIRED'
		);
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 's_tradeengine_trade_strategy_status') THEN
		CREATE TYPE s_tradeengine_trade_strategy_status AS ENUM ('POLLING', 'NEW', 'ACTIVE', 'COMPLETE', 'CANCELLED');
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 's_tradeengine_instrument_type') THEN
		CREATE TYPE s_tradeengine_instrument_type AS ENUM ('SPOT', 'FUTURE_PERPETUAL', 'FUTURE', 'OPTION', 'FORWARD', 'MOVE');
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 's_tradeengine_asset_pair') THEN
		CREATE TYPE s_tradeengine_asset_pair AS ENUM ('USDT', 'BTC', 'USD', 'USDC', 'ETH');
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 's_tradeengine_order_type') THEN
		CREATE TYPE s_tradeengine_order_type AS ENUM ('LIMIT', 'MARKET', 'STOP_MARKET', 'STOP_LIMIT', 'TAKE_PROFIT_LIMIT', 'TAKE_PROFIT_MARKET', 'TRAILING_STOP_MARKET');
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 's_tradeengine_time_in_force') THEN
		CREATE TYPE s_tradeengine_time_in_force AS ENUM (
			'TIME_IN_FORCE_UNREQUIRED',
			'GOOD_TILL_CANCELLED',
			'IMMEDIATE_OR_CANCEL',
			'FILL_OR_KILL',
			'GOOD_TILL_CROSSING'
		);
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 's_tradeengine_working_type') THEN
		CREATE TYPE s_tradeengine_working_type AS ENUM (
			'WORKING_TYPE_UNREQUIRED',
			'MARK_PRICE',
			'CONTRACT_PRICE'
		);
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 's_tradeengine_execution_strategy') THEN
		CREATE TYPE s_tradeengine_execution_strategy AS ENUM (
			'DMA_LIMIT',
			'DMA_MARKET',
			'DCA_FIRST_MARKET_REST_LIMIT',
			'DCA_ALL_LIMIT',
			'TWAP',
			'VWAP'
		);
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 's_tradeengine_dca_execution_strategy') THEN
		CREATE TYPE s_tradeengine_dca_execution_strategy AS ENUM (
			'CONSTANT',
			'LINEAR',
			'EXPONENTIAL'
		);
	END IF;
END
$$;

CREATE TABLE IF NOT EXISTS s_tradeengine_trade_strategies (
	trade_strategy_id uuid DEFAULT uuid_generate_v4(),

	actor_id VARCHAR(32) NOT NULL,
	humanized_actor_name VARCHAR(256) NOT NULL,
	actor_type s_tradeengine_actor_type NOT NULL,

	idempotency_key VARCHAR(256) UNIQUE,

	execution_strategy s_tradeengine_execution_strategy NOT NULL,
	instrument_type s_tradeengine_instrument_type NOT NULL,
	trade_side s_tradeengine_trade_side NOT NULL,

	asset VARCHAR(8) NOT NULL,
	pair VARCHAR(4) NOT NULL,

	entries DECIMAL[] NOT NULL,
	stop_loss DECIMAL NOT NULL,
	take_profits DECIMAL[] NOT NULL,

	current_price DECIMAL NOT NULL,

	status s_tradeengine_trade_strategy_status NOT NULL DEFAULT 'POLLING',

	tradeable_venues VARCHAR(64)[] NOT NULL,

	created TIME NOT NULL,
	last_updated TIME NOT NULL,

	PRIMARY KEY(trade_strategy_id)
);

CREATE TABLE IF NOT EXISTS s_tradeengine_trade_strategy_participants (
	trade_participant_id uuid DEFAULT uuid_generate_v4(),

	trade_strategy_id uuid,
	user_id VARCHAR(20) NOT NULL,
	
	is_bot BOOLEAN NOT NULL DEFAULT FALSE,

	size DECIMAL NOT NULL,
	risk DECIMAL NOT NULL,
 
	venue s_tradeengine_venue NOT NULL,
	exchange_order_ids VARCHAR(256)[] NOT NULL,

	executed TIME NOT NULL,

	PRIMARY KEY(trade_participant_id),
	CONSTRAINT fk_tradeengine
		FOREIGN KEY(trade_strategy_id)
			REFERENCES s_tradeengine_trade_strategies(trade_strategy_id) ON DELETE SET NULL,
	
	UNIQUE(trade_strategy_id, user_id)
);
