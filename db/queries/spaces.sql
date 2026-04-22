-- name: GetSpace :one
SELECT * FROM app.spaces
WHERE id = $1;

-- name: ListSpaces :many
SELECT * FROM app.spaces
WHERE organization_id = $1
ORDER BY name;

-- name: CreateSpace :one
INSERT INTO app.spaces (
  id, organization_id, name, created_by, updated_by
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateSpace :one
UPDATE app.spaces
SET name = $2,
    updated_by = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteSpace :exec
DELETE FROM app.spaces
WHERE id = $1;

-- name: GetSpaceWithOrganization :one
SELECT 
  s.*,
  o.name as organization_name
FROM app.spaces s
JOIN app.organizations o ON s.organization_id = o.id
WHERE s.id = $1;

-- name: CountSpacesByOrganization :one
SELECT COUNT(*) FROM app.spaces
WHERE organization_id = $1;
