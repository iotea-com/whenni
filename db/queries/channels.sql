-- name: GetChannel :one
SELECT * FROM app.channels
WHERE id = $1;

-- name: ListChannels :many
SELECT * FROM app.channels
WHERE space_id = $1
ORDER BY name;

-- name: CreateChannel :one
INSERT INTO app.channels (
  id, name, space_id, config, created_by, updated_by
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateChannel :one
UPDATE app.channels
SET name = $2,
    config = $3,
    updated_by = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: PublishChannel :one
UPDATE app.channels
SET published_at = NOW(),
    published_by = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UnpublishChannel :one
UPDATE app.channels
SET published_at = NULL,
    published_by = NULL,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteChannel :exec
DELETE FROM app.channels
WHERE id = $1;

-- name: ListPublishedChannels :many
SELECT * FROM app.channels
WHERE space_id = $1 AND published_at IS NOT NULL
ORDER BY published_at DESC;
