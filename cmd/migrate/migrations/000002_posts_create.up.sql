CREATE TABLE IF NOT EXISTS posts (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    title TEXT NOT NULL,
    user_id BIGINT NOT NULL,
    created_at timestamp(0) WITH time zone NOT NULL DEFAULT NOW()
)