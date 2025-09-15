-- +goose Up
CREATE TABLE todos (
    id SERIAL PRIMARY KEY,
    task TEXT NOT NULL
);

-- +goose Down
DROP TABLE todos;