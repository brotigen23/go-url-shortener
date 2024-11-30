-- +goose Up
-- +goose StatementBegin
CREATE TABLE Short_URLs(
    ID SERIAL PRIMARY KEY,
    URL VARCHAR, 
    Alias VARCHAR UNIQUE
);
CREATE TABLE Users(
    ID SERIAL PRIMARY KEY, 
    Name VARCHAR UNIQUE
);
CREATE TABLE Users_URLs(
    ID SERIAL PRIMARY KEY,
    User_ID INTEGER REFERENCES Users (ID),
    URL_ID INTEGER REFERENCES Short_URLs (ID)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Users_URLs;
DROP TABLE IF EXISTS Users;
DROP TABLE IF EXISTS Short_URLs;
-- +goose StatementEnd
