CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS s_discord_touches (
	touches_id uuid DEFAULT uuid_generate_v4(),
	idempotency_key VARCHAR(255) NOT NULL UNIQUE,
	updated TIME NOT NULL,
	sender_id VARCHAR(255) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_s_discord_touches_idempotency_key
ON s_discord_touches(idempotency_key);
