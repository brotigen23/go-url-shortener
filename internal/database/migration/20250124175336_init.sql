-- +goose Up
-- +goose StatementBegin
CREATE TABLE short_url(
    id SERIAL PRIMARY KEY,
    url VARCHAR UNIQUE UNIQUE, 
    short_url VARCHAR UNIQUE,
    username VARCHAR,
    is_deleted BOOLEAN DEFAULT FALSE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS short_url;
-- +goose StatementEnd
