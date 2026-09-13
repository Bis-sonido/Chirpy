-- name: GetChirps :many
SELECT *
FROM chirps
ORDER BY created_at;

-- name: GetChirp :one
SELECT *
FROM chirps
WHERE id = $1;

-- name: GetChirpsByUserId :many
SELECT *
FROM chirps
WHERE user_id = $1
ORDER BY created_at ASC;

-- name: GetChirpsSortedByCreatedAtASC :many
SELECT *
FROM chirps
ORDER BY created_at ASC;

-- name: GetChirpsSortedByCreatedAtDESC :many
SELECT *
FROM chirps
ORDER BY created_at DESC;