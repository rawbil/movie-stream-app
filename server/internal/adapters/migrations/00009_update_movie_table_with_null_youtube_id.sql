-- +goose Up
ALTER TABLE movies
MODIFY COLUMN youtube_id TEXT,
MODIFY COLUMN admin_review TEXT,
MODIFY COLUMN ranking_value INT,
MODIFY COLUMN ranking_name VARCHAR(20);

-- youtube_id TEXT NOT NULL,
--     admin_review TEXT NOT NULL,
--     ranking_value INT NOT NULL,
--     ranking_name VARCHAR(20) NOT NULL

-- +goose Down
SELECT 'down SQL query';
