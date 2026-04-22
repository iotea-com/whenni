-- name: GetThing :one
SELECT * FROM app.things
WHERE id = $1;

-- name: ListThings :many
SELECT * FROM app.things
WHERE space_id = $1
ORDER BY name;

-- name: ListThingsByCategory :many
SELECT * FROM app.things
WHERE space_id = $1 AND thing_category = $2
ORDER BY name;

-- name: CreateThing :one
INSERT INTO app.things (
  id, name, space_id, attributes, internal, thing_category, created_by, updated_by
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: UpdateThing :one
UPDATE app.things
SET name = $2,
    attributes = $3,
    updated_by = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteThing :exec
DELETE FROM app.things
WHERE id = $1;

-- name: ListInternalThings :many
SELECT * FROM app.things
WHERE space_id = $1 AND internal = true
ORDER BY name;
