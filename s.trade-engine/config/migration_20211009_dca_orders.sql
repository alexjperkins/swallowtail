ALTER TYPE s_tradeengine_order_type ADD VALUE 'STOP_MARKET';
ALTER TYPE s_tradeengine_order_type ADD VALUE 'DCA_ALL_LIMIT';
ALTER TYPE s_tradeengine_order_type ADD VALUE 'DCA_FIRST_MARKET_REST_LIMIT';
   
ALTER TABLE s_tradeengine_trades 
	ALTER COLUMN entry TYPE DECIMAL[]
	USING array[entry]::DECIMAL[];

ALTER TABLE s_tradeengine_trades RENAME entry TO entries;
