-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS aliases(
    URL VARCHAR UNIQUE, 
    Alias VARCHAR
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS aliases;
-- +goose StatementEnd
