-- +goose Up

CREATE TABLE IF NOT EXISTS movies(
    movie_id BIGINT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    public_id BINARY(16) NOT NULL UNIQUE,
    imdb_id VARCHAR(50) NOT NULL UNIQUE,
    title VARCHAR(50) NOT NULL,
    poster_path TEXT NOT NULL,
    youtube_id TEXT NOT NULL,
    admin_review TEXT NOT NULL,
    ranking_value INT NOT NULL,
    ranking_name VARCHAR(20) NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS movies;