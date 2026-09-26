-- name: CreateGenre :execresult
INSERT IGNORE INTO genres(genre_name)
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
GROUP_CONCAT(g.genre_name ORDER BY g.genre_name) AS genres
from movies m
 JOIN movie_genres mg
ON m.movie_id = mg.movie_id
JOIN genres g
ON mg.genre_id = g.genre_id

GROUP BY
    m.movie_id,
    m.imdb_id, 
    m.title,
    m.poster_path, 
    m.youtube_id, 
    m.admin_review,
    m.ranking_value, 
    m.ranking_name
ORDER BY m.ranking_value
;

-- name: GetMovie :one
SELECT 
BIN_TO_UUID(m.public_id) AS public_id, 
m.imdb_id, 
m.title,
m.poster_path, 
m.youtube_id, 
m.admin_review,
m.ranking_value, 
m.ranking_name,  
GROUP_CONCAT(g.genre_name ORDER BY g.genre_name)  AS genres
from movies m
 JOIN movie_genres mg
ON m.movie_id = mg.movie_id
JOIN genres g
ON mg.genre_id = g.genre_id
WHERE m.public_id = ?
GROUP BY
  m.movie_id,
  m.public_id,
  m.imdb_id,
  m.title,
  m.poster_path,
  m.youtube_id,
  m.admin_review,
  m.ranking_value,
  m.ranking_name;

-- name: GetMovieInternal :one
SELECT * FROM movies
WHERE public_id = ?;

-- name: CreateMovie :execresult
INSERT INTO movies(public_id, imdb_id, title, poster_path, youtube_id)
VALUES(?, ?, ?, ?, ?);

-- name: CreateMovieGenre :execresult
INSERT IGNORE INTO movie_genres(movie_id, genre_id)
VALUES (?, ?);

-- name: GetMovieGenre :one
SELECT * FROM movie_genres
WHERE movie_id = ? AND genre_id = ?;

-- name: GetUniqueMovie :one
SELECT * FROM movies
WHERE imdb_id = ?;

-- name: CreateUser :execresult
INSERT INTO users(username, email, password)
VALUES (?, ?, ?);

-- name: CreateRefreshToken :execresult
INSERT IGNORE INTO refresh_tokens(user_id, hashed_token)
VALUES (?, ?);

-- name: UpdateRefreshToken :execresult
UPDATE refresh_tokens
SET hashed_token = ?,
    revoked = ?
WHERE user_id = ?;

-- name: RevokeRefreshToken :execresult
UPDATE refresh_tokens
SET revoked = ?
WHERE user_id = ?;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE user_id = ?;

-- name: CreateRole :execresult
INSERT INTO roles(role)
VALUES (?);

-- name: CreateUserPermission :execresult
INSERT INTO user_permissions(permission)
VALUES (?);

-- name: CreateRolePermission :execresult
INSERT IGNORE INTO role_permissions(role_id, permission_id)
VALUES (?, ?);

-- name: CreateUserRoles :execresult
INSERT IGNORE INTO user_roles(user_id, role_id)
VALUES (?, ?);

-- name: CreateUserGenre :execresult
INSERT IGNORE INTO user_fav_genres(user_id, genre_id)
VALUES (?, ?);

-- name: GetRole :one
SELECT * FROM roles
WHERE role = ?;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ?;

-- name: GetUserByID :one
SELECT * FROM users
WHERE user_id = ?;

-- name: UpdateMovieRankings :execresult
UPDATE movies
SET ranking_name = ?,
    ranking_value = ?
WHERE public_id = ?;

-- name: AddRanking :execresult
INSERT IGNORE INTO rankings(ranking_name, ranking_value)
VALUES (?, ?);

-- name: GetRanking :one
SELECT * FROM rankings
WHERE ranking_name = ?;

-- name: GetRankings :many
SELECT * FROM rankings
ORDER BY ranking_value;

-- name: AddMovieReview :execresult
INSERT INTO movie_reviews(movie_id, review)
values(?, ?);

-- name: GetMovieReviews :many
SELECT * FROM movie_reviews
WHERE movie_id = ?
ORDER BY updated_at;

-- name: GetRecommendedMovies :many
SELECT 
    BIN_TO_UUID(m.public_id) AS public_id, 
    m.imdb_id, 
    m.title, 
    m.poster_path, 
    m.youtube_id, 
    m.ranking_value, 
    m.ranking_name,
    GROUP_CONCAT(g.genre_name ORDER BY g.genre_name) AS genres
FROM user_fav_genres fav
JOIN movie_genres mg
    ON fav.genre_id = mg.genre_id
JOIN movies m
    ON m.movie_id = mg.movie_id
JOIN genres g
    ON g.genre_id = fav.genre_id
WHERE fav.user_id = ?
GROUP BY
    BIN_TO_UUID(m.public_id), 
    m.imdb_id, 
    m.title, 
    m.poster_path, 
    m.youtube_id, 
    m.ranking_value, 
    m.ranking_name
ORDER BY m.ranking_value
LIMIT 5;