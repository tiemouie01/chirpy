-- +goose Up

ALTER TABLE users
ADD is_chirpy_red BOOLEAN DEFAULT false;

-- +goose Down
DROP TABLE users;
DROP TABLE chirps;
DROP TABLE refresh_tokens;