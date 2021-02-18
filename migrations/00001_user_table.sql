-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id bigserial PRIMARY KEY NOT NULL UNIQUE,
    created_at timestamptz,
    updated_at timestamptz,
    deleted_at timestamptz,
    uuid uuid DEFAULT uuid_generate_v4() UNIQUE,
    email text NOT NULL UNIQUE,
    password text NOT NULL,
    name text NOT NULL,
    surname text NOT NULL,
    birthday timestamptz NOT NULL,
    photo text,
    phone text
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd