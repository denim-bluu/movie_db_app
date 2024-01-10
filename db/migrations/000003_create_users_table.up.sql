CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    name TEXT NOT NULL,
    email CITEXT UNIQUE NOT NULL,
    password_hash BYTEA NOT NULL,
    activated BOOLEAN NOT NULL,
    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS users_email_idx ON users (email);

-- Insert statements for each user
INSERT INTO users (name, email, password_hash, activated)
VALUES
    ('Alice Smith', 'alice@example.com', '\x243261243132246c717170677644692e6161597441457a456958514f2e624731302e747a50384b767a2e34364c4742716f4d76494b4158464a426471', true),
    ('Faith Smith', 'faith@example.com', '\x243261243132246c717170677644692e6161597441457a456958514f2e624731302e747a50384b767a2e34364c4742716f4d76494b4158464a426471', true);
