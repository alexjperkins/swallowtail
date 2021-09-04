CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE IF EXISTS s_discord_touches;

CREATE TABLE IF NOT EXISTS s_discord_touches (
	touch_id uuid DEFAULT uuid_generate_v4(),
	idempotency_key VARCHAR(1024) NOT NULL UNIQUE,
	updated TIME NOT NULL DEFAULT now(),
	sender_id VARCHAR(255) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_s_discord_touches_idempotency_key
	ON s_discord_touches(idempotency_key);
