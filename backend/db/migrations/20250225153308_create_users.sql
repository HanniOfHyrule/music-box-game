-- +goose Up
-- +goose StatementBegin
CREATE TABLE "users" (
  id SERIAL PRIMARY KEY,
  api_token VARCHAR(64) UNIQUE NOT NULL,
  spotify_token VARCHAR(255) NOT NULL,
  spotify_refresh_token VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "users";
-- +goose StatementEnd
