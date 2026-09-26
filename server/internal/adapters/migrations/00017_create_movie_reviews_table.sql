-- +goose Up
CREATE TABLE IF NOT EXISTS movie_reviews(
    id BIGINT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    movie_id BIGINT NOT NULL,
    review TEXT NOT NULL,
    created_at DATETIME DEFAULT NOW() NOT NULL,
    updated_at DATETIME DEFAULT NOW() NOT NULL,
    FOREIGN KEY (movie_id) REFERENCES movies(movie_id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS movie_reviews;
