-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS metrics (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    type VARCHAR NOT NULL,
    delta INT,
    value DOUBLE PRECISION
);
CREATE UNIQUE INDEX IF NOT EXISTS name_type_idx ON metrics (name,type);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP INDEX name_type_idx;
DROP TABLE IF EXISTS metrics;
-- +goose StatementEnd
