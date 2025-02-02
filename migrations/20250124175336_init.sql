-- +goose Up
-- +goose StatementBegin
CREATE TABLE short_url(
    id SERIAL PRIMARY KEY,
    url VARCHAR(2048) UNIQUE, 
    short_url VARCHAR(16) UNIQUE,
    username VARCHAR(16) UNIQUE,
    is_deleted BOOLEAN DEFAULT FALSE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE short_url;
-- +goose StatementEnd
