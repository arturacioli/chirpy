-- name: CreateUser :one
INSERT INTO users (id, createdAt, UpdatedAt, email)
VALUES(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1
)
RETURNING *;
