-- name: GetTag :one
SELECT * FROM app.tags
WHERE id = $1;

-- name: ListTags :many
SELECT * FROM app.tags
WHERE space_id = $1
ORDER BY name;

-- name: CreateTag :one
INSERT INTO app.tags (
  name, space_id, created_by, updated_by
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateTag :one
UPDATE app.tags
SET name = $2,
    updated_by = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteTag :exec
DELETE FROM app.tags
WHERE id = $1;

-- name: ApplyTag :one
INSERT INTO app.applied_tags (
  space_id, tag_id, channel_id, model_id, thing_id, created_by, updated_by
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: RemoveAppliedTag :exec
DELETE FROM app.applied_tags
WHERE id = $1;

-- name: ListAppliedTagsByChannel :many
SELECT 
  at.*,
  t.name as tag_name
FROM app.applied_tags at
JOIN app.tags t ON at.tag_id = t.id
WHERE at.channel_id = $1;

-- name: ListAppliedTagsByModel :many
SELECT 
  at.*,
  t.name as tag_name
FROM app.applied_tags at
JOIN app.tags t ON at.tag_id = t.id
WHERE at.model_id = $1;

-- name: ListAppliedTagsByThing :many
SELECT 
  at.*,
  t.name as tag_name
FROM app.applied_tags at
JOIN app.tags t ON at.tag_id = t.id
WHERE at.thing_id = $1;
