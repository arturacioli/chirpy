-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users where email = $1;

-- name: GetUserFromRefreshToken :one
SELECT u.* FROM users u
JOIN refresh_tokens r ON u.id = r.user_id
WHERE r.token = $1 AND
r.revoked_at IS null AND
r.expires_at > NOW();

-- name: UpdateUserEmailAndPassword :one
UPDATE users SET email = $1, hashed_password = $2, updated_at = NOW()
    WHERE id = $3 RETURNING *;

-- name: UpdateUserToRed :one
UPDATE users SET is_chirpy_red = true where id = $1 RETURNING *;
