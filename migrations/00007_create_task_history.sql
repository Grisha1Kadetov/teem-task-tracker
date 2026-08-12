-- +goose Up
-- +goose StatementBegin
CREATE TABLE task_history (
    id BINARY(16) NOT NULL,
    task_id BINARY(16) NOT NULL,
    changed_by BINARY(16) NOT NULL,
    changes JSON NOT NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    PRIMARY KEY (id),
    FOREIGN KEY (task_id) REFERENCES tasks (id),
    FOREIGN KEY (changed_by) REFERENCES users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE task_history;
-- +goose StatementEnd
