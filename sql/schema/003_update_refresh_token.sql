-- +goose Up

ALTER TABLE refresh_tokens
ALTER COLUMN user_id TYPE uuid,
ALTER COLUMN user_id SET NOT NULL;

-- +goose Down
DROP TABLE users;
DROP TABLE chirps;
DROP TABLE refresh_tokens;