ALTER TABLE s_payments_payments ADD COLUMN payment_timestamp TIMESTAMP NOT NULL DEFAULT now();
ALTER TABLE s_payments_payments DROP COLUMN timestamp;
