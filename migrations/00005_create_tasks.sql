-- +goose Up
-- +goose StatementBegin
CREATE TABLE tasks (
    id BINARY(16) NOT NULL,
    team_id BINARY(16) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(32) NOT NULL,
    created_by BINARY(16) NOT NULL,
    assignee_id BINARY(16) NOT NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    closed_at DATETIME(6),
    version BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (team_id) REFERENCES teams (id),
    FOREIGN KEY (created_by) REFERENCES users (id),
    FOREIGN KEY (assignee_id) REFERENCES users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tasks;
-- +goose StatementEnd
