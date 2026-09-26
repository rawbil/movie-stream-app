-- +goose Up
CREATE TABLE IF NOT EXISTS rankings(
    id int PRIMARY KEY NOT NULL AUTO_INCREMENT,
    ranking_name VARCHAR(20) UNIQUE NOT NULL,
    ranking_value int NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE IF EXISTS rankings;
