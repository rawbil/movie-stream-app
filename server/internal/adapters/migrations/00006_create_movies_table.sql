-- +goose Up

CREATE TABLE IF NOT EXISTS movies(
    movie_id BIGINT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    imdb_id VARCHAR(50) NOT NULL,
    title VARCHAR(50) NOT NULL,
    poster_path TEXT NOT NULL,
    youtube_id TEXT NOT NULL,
    admin_review TEXT NOT NULL,
    genre_id BIGINT NOT NULL,
    ranking_value INT NOT NULL,
    ranking_name VARCHAR(20) NOT NULL,

    FOREIGN KEY (genre_id) REFERENCES genres(genre_id)
);

-- +goose Down
DROP TABLE IF EXISTS movies;