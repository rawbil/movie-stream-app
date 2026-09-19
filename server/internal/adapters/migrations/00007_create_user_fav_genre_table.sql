-- +goose Up
CREATE TABLE IF NOT EXISTS user_fav_genres(
    user_id BIGINT NOT NULL,
    genre_id BIGINT NOT NULL,

    PRIMARY KEY(user_id, genre_id),

    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (genre_id) REFERENCES genres(genre_id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS user_fav_genres;
