CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGSERIAL NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id BIGSERIAL NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at timestamp(0) WITH time zone NOT NULL DEFAULT NOW()
)