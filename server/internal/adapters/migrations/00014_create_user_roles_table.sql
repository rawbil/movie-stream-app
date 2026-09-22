-- +goose Up
CREATE TABLE IF NOT EXISTS user_roles(
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,

    PRIMARY KEY(user_id, role_id),

    CONSTRAINT roles_u_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT roles_r_id FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS user_roles;