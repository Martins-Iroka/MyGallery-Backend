CREATE TABLE IF NOT EXISTS photo_posts(
    id bigserial PRIMARY KEY,
    photographer varchar(255) NOT NULL,
    original text NOT NULL CHECK (original ~ '^https?://'),
    large2x text NOT NULL CHECK (large2x ~ '^https?://'),
    large text NOT NULL CHECK (large ~ '^https?://'),
    medium text NOT NULL CHECK (medium ~ '^https?://'),
    small text NOT NULL CHECK (small ~ '^https?://'),
    portrait text NOT NULL CHECK (portrait ~ '^https?://'),
    landscape text NOT NULL CHECK (landscape ~ '^https?://'),
    tiny text NOT NULL CHECK (tiny ~ '^https?://')
);