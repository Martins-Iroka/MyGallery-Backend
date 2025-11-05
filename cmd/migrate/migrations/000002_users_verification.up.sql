CREATE TABLE IF NOT EXISTS users_verification_tracking(
    token bytea PRIMARY KEY,
    user_id bigserial NOT NULL
);