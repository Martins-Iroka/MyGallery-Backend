CREATE TABLE IF NOT EXISTS video_download_files(
    id bigserial PRIMARY KEY NOT NULL,
    video_post_id bigint REFERENCES video_posts,
    video_link text NOT NULL CHECK(video_link ~ '^https?://'),
    video_size bigint NOT NULL
);