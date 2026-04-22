-- name: GetModel :one
SELECT * FROM app.models
WHERE id = $1;

-- name: ListModels :many
SELECT * FROM app.models
WHERE space_id = $1
ORDER BY name;

-- name: CreateModel :one
INSERT INTO app.models (
  id, space_id, name, attributes, created_by, updated_by
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateModel :one
UPDATE app.models
SET name = $2,
    attributes = $3,
    updated_by = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteModel :exec
DELETE FROM app.models
WHERE id = $1;
