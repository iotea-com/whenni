-- name: GetModel :one
SELECT * FROM app.models
WHERE id = $1;

-- name: ListModels :many
SELECT * FROM app.models
WHERE space_id = $1
ORDER BY name;

-- name: CountModelsFiltered :one
SELECT COUNT(*) FROM app.models m
WHERE m.space_id = $1
  AND ($2::text = '' OR m.id ILIKE '%' || $2 || '%' OR m.name ILIKE '%' || $2 || '%')
  AND (
    cardinality($3::text[]) = 0
    OR m.id IN (
      SELECT DISTINCT at.model_id
      FROM app.applied_tags at
      WHERE at.tag_id = ANY($3::text[])
    )
  );

-- name: ListModelsFiltered :many
SELECT * FROM app.models m
WHERE m.space_id = $1
  AND ($2::text = '' OR m.id ILIKE '%' || $2 || '%' OR m.name ILIKE '%' || $2 || '%')
  AND (
    cardinality($3::text[]) = 0
    OR m.id IN (
      SELECT DISTINCT at.model_id
      FROM app.applied_tags at
      WHERE at.tag_id = ANY($3::text[])
    )
  )
ORDER BY m.updated_at DESC
LIMIT $4 OFFSET $5;

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

-- name: DeleteModelsBySpace :exec
DELETE FROM app.models
WHERE space_id = $1;
