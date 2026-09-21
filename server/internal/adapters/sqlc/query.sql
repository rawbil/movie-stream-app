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
BIN_TO_UUID(m.public_id) AS public_id,
m.imdb_id, 
m.title,
m.poster_path, 
m.youtube_id, 
m.admin_review,
m.ranking_value, 
m.ranking_name,  
g.genre_name 
from movies m
 JOIN movie_genres mg
ON m.movie_id = mg.movie_id
JOIN genres g
ON mg.genre_id = g.genre_id
ORDER BY m.ranking_value
;

-- name: GetMovie :one
SELECT 
BIN_TO_UUID(m.public_id), 
m.imdb_id, 
m.title,
m.poster_path, 
m.youtube_id, 
m.admin_review,
m.ranking_value, 
m.ranking_name,  
g.genre_name 
from movies m
 JOIN movie_genres mg
ON m.movie_id = mg.movie_id
JOIN genres g
ON mg.genre_id = g.genre_id
WHERE m.public_id = ?
;

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