CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email citext UNIQUE NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password bytea NOT NULL,
    created_at timestamp(0) WITH time zone NOT NULL DEFAULT NOW()
)