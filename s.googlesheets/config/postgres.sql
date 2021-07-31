CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$ 
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'googlesheet_type') THEN
		CREATE TYPE pager AS ENUM ('PORTFOLIO', 'PLAIN');
	END IF;
END
$$;

CREATE TABLE IF NOT EXISTS s_googlehseets_sheet (
	googlesheet_id uuid DEFAULT uuid_generate_v4();

	spreadsheet_id VARCHAR(32) NOT NULL,
	sheet_id VARCHAR(32) NOT NULL UNIQUE,

	sheet_type googlesheet_type NOT NULL DEFAULT 'PLAIN',

	user_id VARCHAR(20) NOT NULL,
	with_pager_on_error BOOLEAN DEFAULT FALSE,
	with_pager_on_target BOOLEAN DEFAULT FALSE,

	created TIME NOT NULL DEFAULT now(),
	updated TIME NOT NULL DEFAULT now(),

	active BOOLEAN DEFAULT TRUE,

	PRIMARY KEY(googlesheet_id)
);
