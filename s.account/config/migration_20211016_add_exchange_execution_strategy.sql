DO $$ 
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'exchange_execution_strategy') THEN
		CREATE TYPE exchange_execution_strategy AS ENUM ('PRIMARY_ONLY', 'ATTEMPT_ALL_REGISTERED')
	END IF;
END
$$;

ALTER TABLE s_account_accounts ADD COLUMN default_exchange_execution_strategy NOT NULL DEFAULT 'PRIMARY_ONLY';
