-- name: GetThing :one
SELECT * FROM app.things
WHERE id = $1;

-- name: CountThingsFiltered :one
SELECT COUNT(DISTINCT t.id)
FROM app.things t
WHERE t.space_id = $1
  AND t.internal = false
  AND ($2::text IS NULL OR t.thing_category = $2)
  AND (
    $3::text IS NULL
    OR t.id ILIKE $3
    OR t.name ILIKE $3
    OR t.thing_category ILIKE $3
  )
  AND (
    $4::text[] IS NULL
    OR EXISTS (
      SELECT 1
      FROM app.applied_tags at
      WHERE at.thing_id = t.id
        AND at.tag_id = ANY($4::text[])
    )
  );

-- name: ListThingsFiltered :many
SELECT t.*
FROM app.things t
WHERE t.space_id = $1
  AND t.internal = false
  AND ($2::text IS NULL OR t.thing_category = $2)
  AND (
    $3::text IS NULL
    OR t.id ILIKE $3
    OR t.name ILIKE $3
    OR t.thing_category ILIKE $3
  )
  AND (
    $4::text[] IS NULL
    OR EXISTS (
      SELECT 1
      FROM app.applied_tags at
      WHERE at.thing_id = t.id
        AND at.tag_id = ANY($4::text[])
    )
  )
ORDER BY t.thing_category DESC
LIMIT $5 OFFSET $6;

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

-- name: DeleteThingsBySpace :exec
DELETE FROM app.things
WHERE space_id = $1;

-- name: ListInternalThings :many
SELECT * FROM app.things
WHERE space_id = $1 AND internal = true
ORDER BY name;

-- name: CreateCertificate :one
INSERT INTO app.certificates (
  id, name, policy, space_id
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetCertificate :one
SELECT * FROM app.certificates
WHERE id = $1;

-- name: DeleteCertificatesBySpace :exec
DELETE FROM app.certificates
WHERE space_id = $1;
