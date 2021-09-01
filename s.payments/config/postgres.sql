DO $$ 
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'payment_tyep') THEN
		CREATE TYPE pager AS ENUM ('FUTURES');
	END IF;
END
$$;

CREATE TABLE IF NOT EXISTS s_payments_payments (
	user_id VARCHAR(20) NOT NULL,
	transaction_id VARCHAR(256) NOT NULL UNIQUE,
	timestamp TIME NOT NULL,
	amount_in_usdt DECIMAL NOT NULL,
	audit_note VARCHAR(256),

	PRIMARY KEY(transaction_id)
)

CREATE INDEX IF NOT EXISTS idx_s_payments_txis 
	ON s_payments(transaction_id);
