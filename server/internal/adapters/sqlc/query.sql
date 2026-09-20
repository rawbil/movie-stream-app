-- name: CreateGenre :execresult
INSERT INTO genres(genre_name)
VALUES (?);

-- name: UpdateGenre :execresult
UPDATE genres
SET genre_name = ?
WHERE genre_id = ?;

-- name: ListGenres :many
SELECT * FROM genres
ORDER BY genre_name;

-- name: GetGenre :one
SELECT * FROM genres
WHERE genre_name = ?;

-- name: GetGenreByID :one
SELECT * FROM genres
WHERE genre_id = ?;

-- name: DeleteGenre :exec
DELETE FROM genres
WHERE genre_name = ?;

-- name: ListMovies :many
SELECT * FROM movies 
ORDER BY ranking_value;

-- name: GetMovie :one
SELECT * FROM movies
WHERE movie_id = ?;

-- name: CreateMovie :execresult
INSERT INTO movies(imdb_id, title, poster_path)
VALUES(?, ?, ?);