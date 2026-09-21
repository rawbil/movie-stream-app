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
SELECT 
    movie_id,
    BIN_TO_UUID(public_id) AS public_id,
    imdb_id,
    title,
    poster_path,
    youtube_id,
    admin_review,
    ranking_value,
    ranking_name
 FROM movies 
ORDER BY ranking_value;

-- name: GetMovie :one
SELECT 
    movie_id, 
    BIN_TO_UUID(public_id) AS public_id, 
    imdb_id, title, 
    poster_path, youtube_id, 
    admin_review, 
    ranking_name, 
    ranking_value 
FROM movies
WHERE public_id = ?;

-- name: CreateMovie :execresult
INSERT INTO movies(public_id, imdb_id, title, poster_path, youtube_id)
VALUES(?, ?, ?, ?, ?);

-- name: CreateMovieGenre :execresult
INSERT INTO movie_genres(movie_id, genre_id)
VALUES (?, ?);

-- name: GetMovieGenre :one
SELECT * FROM movie_genres
WHERE movie_id = ? AND genre_id = ?;

-- name: GetUniqueMovie :one
SELECT * FROM movies
WHERE imdb_id = ?;