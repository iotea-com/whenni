-- name: CreateSession :one
INSERT INTO app.sessions (
  session_token, user_id, expires
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetSession :one
SELECT * FROM app.sessions
WHERE session_token = $1;

-- name: GetSessionWithUser :one
SELECT 
  s.*,
  u.id as user_id,
  u.name as user_name,
  u.email as user_email,
  u.image as user_image
FROM app.sessions s
JOIN app.users u ON s.user_id = u.id
WHERE s.session_token = $1 AND s.expires > NOW();

-- name: UpdateSession :one
UPDATE app.sessions
SET expires = $2
WHERE session_token = $1
RETURNING *;

-- name: DeleteSession :exec
DELETE FROM app.sessions
WHERE session_token = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM app.sessions
WHERE expires < NOW();

-- name: CreateAccount :one
INSERT INTO app.accounts (
  user_id, type, provider, provider_account_id,
  refresh_token, access_token, expires_at, token_type,
  scope, id_token, session_state
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM app.accounts
WHERE provider = $1 AND provider_account_id = $2;

-- name: GetAccountsByUser :many
SELECT * FROM app.accounts
WHERE user_id = $1;

-- name: DeleteAccount :exec
DELETE FROM app.accounts
WHERE id = $1;

-- name: CreateAuthToken :one
INSERT INTO app.auth_tokens (
  user_id, token, type, expires_at
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetAuthToken :one
SELECT * FROM app.auth_tokens
WHERE token = $1 AND expires_at > NOW();

-- name: DeleteAuthToken :exec
DELETE FROM app.auth_tokens
WHERE token = $1;

-- name: DeleteExpiredAuthTokens :exec
DELETE FROM app.auth_tokens
WHERE expires_at < NOW();

-- name: DeleteAuthTokensByUser :exec
DELETE FROM app.auth_tokens
WHERE user_id = $1 AND type = $2;
