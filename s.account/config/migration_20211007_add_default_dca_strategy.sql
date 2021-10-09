DO $$ 
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'dca_strategy') THEN
		CREATE TYPE dca_strategy AS ENUM ('CONSTANT', 'LINEAR', 'EXPONENTIAL');
	END IF;
END
$$;

ALTER TABLE s_account_accounts ADD COLUMN default_dca_strategy dca_strategy NOT NULL DEFAULT 'LINEAR';
