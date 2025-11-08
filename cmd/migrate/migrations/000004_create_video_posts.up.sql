CREATE TABLE IF NOT EXISTS video_posts(
    id bigserial PRIMARY KEY NOT NULL,
    video_url text NOT NULL CHECK (video_url ~ '^https?://'),
    duration INT NOT NULL
);