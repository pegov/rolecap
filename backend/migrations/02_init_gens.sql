CREATE TABLE IF NOT EXISTS gen (
	id SERIAL PRIMARY KEY,
	user_id INTEGER REFERENCES auth_user(id) ON DELETE CASCADE,

	title TEXT,
	description TEXT,
	access INTEGER,

	body JSONB,
	body_meta JSONB,

	views BIGINT DEFAULT 0,

	date_added TIMESTAMP WITH TIME ZONE,
	date_updated TIMESTAMP WITH TIME ZONE,

	active BOOLEAN DEFAULT TRUE
);