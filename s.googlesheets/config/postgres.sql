CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS s_googlehseets_sheet (
	googlesheets_id uuid DEFAULT uuid_generate_v4();

	spreadsheet_id VARCHAR(32) NOT NULL,
	sheets_id VARCHAR(32) NOT NULL UNIQUE,

	user_id VARCHAR(20) NOT NULL,
	with_pager_on_error BOOLEAN DEFAULT FALSE,
	with_pager_on_target BOOLEAN DEFAULT FALSE,

	created TIME NOT NULL DEFAULT now(),
	updated TIME NOT NULL DEFAULT now(),

	PRIMARY KEY(internal_id)
);
