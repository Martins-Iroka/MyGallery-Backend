CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);

CREATE INDEX IF NOT EXISTS idx_users_verification ON users_verification_tracking (token);

CREATE INDEX IF NOT EXISTS idx_photos_comment ON photos_comment (post_id);

CREATE INDEX IF NOT EXISTS idx_videos_comment ON video_comments (post_id);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);