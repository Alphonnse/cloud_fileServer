-- name: SignUp :one
INSERT INTO users(id, created_at, updated_at, name, email, password) 
values ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: LogIn :one
SELECT * FROM users WHERE email = $1;

-- name: FindWithID :one
SELECT * FROM users WHERE id = $1;
