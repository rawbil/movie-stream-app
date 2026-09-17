-- +goose Up
CREATE TABLE IF NOT EXISTS fav_genres(
    id BIGINT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    genre_id BIGINT NOT NULL,

    CONSTRAINT fk_genre_id FOREIGN KEY (genre_id) REFERENCES genres(genre_id)
);

-- +goose Down
DROP TABLE IF EXISTS fav_genres;
