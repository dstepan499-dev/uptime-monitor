CREATE TABLE websites (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	url VARCHAR(255) NOT NULL,
	interval_seconds INTEGER DEFAULT 60,
	timeout_seconds INTEGER DEFAULT 5,
	is_active BOOLEAN DEFAULT TRUE,
	last_status_code INTEGER,
	last_checked_at TIMESTAMPTZ,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE check_logs (
	id SERIAL PRIMARY KEY,
	website_id INTEGER NOT NULL REFERENCES websites(id) ON DELETE CASCADE,
	status_code INTEGER,
	response_time_ms INTEGER,
	is_successful BOOLEAN NOT NULL,
	error_message TEXT,
	checked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_check_logs_website_id_checked_at ON check_logs(website_id, checked_at DESC);