-- name: GetPermissionSet :one
SELECT * FROM app.permissions
WHERE id = $1;

-- name: ListOrganizationPermissionSets :many
SELECT * FROM app.permissions
WHERE organization_id = $1 AND space_id IS NULL
ORDER BY name;

-- name: CountOrganizationPermissionSets :one
SELECT COUNT(*) FROM app.permissions
WHERE organization_id = $1
  AND space_id IS NULL;

-- name: ListOrganizationPermissionSetsPaginated :many
SELECT * FROM app.permissions
WHERE organization_id = $1
  AND space_id IS NULL
ORDER BY updated_at DESC
LIMIT $2 OFFSET $3;

-- name: ListSpacePermissionSets :many
SELECT * FROM app.permissions
WHERE space_id = $1
ORDER BY name;

-- name: CountSpacePermissionSets :one
SELECT COUNT(*) FROM app.permissions
WHERE space_id = $1;

-- name: ListSpacePermissionSetsPaginated :many
SELECT * FROM app.permissions
WHERE space_id = $1
ORDER BY updated_at DESC
LIMIT $2 OFFSET $3;

-- name: CreatePermissionSet :one
INSERT INTO app.permissions (
  id, organization_id, space_id, name, permissions, created_by, updated_by
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdatePermissionSet :one
UPDATE app.permissions
SET name = $2,
    permissions = $3,
    updated_by = $4,
    space_id = $5,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeletePermissionSet :exec
DELETE FROM app.permissions
WHERE id = $1;

-- name: GetApiKey :one
SELECT * FROM app."apiKeys"
WHERE id = $1;

-- name: ListApiKeysByOrganization :many
SELECT * FROM app."apiKeys"
WHERE organization_id = $1
ORDER BY created_at DESC;

-- name: CountOrganizationApiKeys :one
SELECT COUNT(*) FROM app."apiKeys"
WHERE organization_id = $1
  AND space_id IS NULL;

-- name: ListOrganizationApiKeysPaginated :many
SELECT * FROM app."apiKeys"
WHERE organization_id = $1
  AND space_id IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListApiKeysBySpace :many
SELECT * FROM app."apiKeys"
WHERE space_id = $1
ORDER BY created_at DESC;

-- name: CountSpaceApiKeys :one
SELECT COUNT(*) FROM app."apiKeys"
WHERE organization_id = $1
  AND space_id = $2;

-- name: ListSpaceApiKeysPaginated :many
SELECT * FROM app."apiKeys"
WHERE organization_id = $1
  AND space_id = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: CreateApiKey :one
INSERT INTO app."apiKeys" (
  id, organization_id, organization_permission_set_id,
  space_id, space_permission_set_id, name, created_by, expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: DeleteApiKey :exec
DELETE FROM app."apiKeys"
WHERE id = $1;

-- name: DeleteApiKeysBySpace :exec
DELETE FROM app."apiKeys"
WHERE space_id = $1;

-- name: DeleteExpiredApiKeys :exec
DELETE FROM app."apiKeys"
WHERE expires_at < NOW();

-- name: DeletePermissionSetsBySpace :exec
DELETE FROM app.permissions
WHERE space_id = $1;
