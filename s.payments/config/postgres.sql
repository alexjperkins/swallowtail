CREATE TABLE IF NOT EXISTS s_payments_payments (
	user_id VARCHAR(20) NOT NULL,
	transaction_id VARCHAR(256) NOT NULL UNIQUE,
	timestamp TIME NOT NULL,
	amount_in_usdt DECIMAL NOT NULL,
)

CREATE INDEX IF NOT EXISTS idx_s_payments_txis 
	ON s_payments(transaction_id);
