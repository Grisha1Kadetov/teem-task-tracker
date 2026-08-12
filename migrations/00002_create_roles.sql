-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles (
    id VARCHAR(16) NOT NULL,
    PRIMARY KEY (id)
);

INSERT INTO roles (id)
VALUES ('owner'), ('admin'), ('member');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE roles;
-- +goose StatementEnd
