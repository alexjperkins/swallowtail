ALTER TABLE s_tradeengine_trades 
	ADD COLUMN tradeable_exchanges VARCHAR(64)[] NOT NULL;
