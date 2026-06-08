-- name: GetUser :one
SELECT * FROM app.users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM app.users
WHERE email = $1;

-- name: ListUsers :many
SELECT * FROM app.users
ORDER BY email;

-- name: SearchUsers :many
SELECT u.* FROM app.users u
LEFT JOIN app.organization_members om
  ON om.user_id = u.id
 AND om.organization_id = sqlc.arg('organization_id')
WHERE u.email ILIKE ('%' || sqlc.arg('query') || '%')
  AND (sqlc.arg('organization_id') = '' OR om.user_id IS NULL)
ORDER BY u.email;

-- name: CreateUser :one
INSERT INTO app.users (
  id, name, email, email_verified, password, image
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateUser :one
UPDATE app.users
SET name = COALESCE(sqlc.narg('name'), name),
    email = COALESCE(sqlc.narg('email'), email),
    email_verified = COALESCE(sqlc.narg('email_verified'), email_verified),
    password = COALESCE(sqlc.narg('password'), password),
    image = COALESCE(sqlc.narg('image'), image)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM app.users
WHERE id = $1;

-- name: VerifyUserEmail :one
UPDATE app.users
SET email_verified = NOW()
WHERE id = $1
RETURNING *;

-- name: GetUserWithOrganizations :many
SELECT 
  u.*,
  om.organization_id,
  om.role as organization_role,
  o.name as organization_name
FROM app.users u
LEFT JOIN app.organization_members om ON u.id = om.user_id
LEFT JOIN app.organizations o ON om.organization_id = o.id
WHERE u.id = $1;
