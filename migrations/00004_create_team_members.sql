-- +goose Up
-- +goose StatementBegin
CREATE TABLE team_members (
    team_id BINARY(16) NOT NULL,
    user_id BINARY(16) NOT NULL,
    role VARCHAR(16) NOT NULL,
    UNIQUE (team_id, user_id),
    FOREIGN KEY (team_id) REFERENCES teams (id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (role) REFERENCES roles (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE team_members;
-- +goose StatementEnd
