-- name: GetChannel :one
SELECT * FROM app.channels
WHERE id = $1;

-- name: ListChannels :many
SELECT * FROM app.channels
WHERE space_id = $1
ORDER BY name;

-- name: CountChannelsFiltered :one
SELECT COUNT(*) FROM app.channels c
WHERE c.space_id = $1
  AND ($2::text = '' OR c.id ILIKE '%' || $2 || '%' OR c.name ILIKE '%' || $2 || '%')
  AND (
    cardinality($3::text[]) = 0
    OR c.id IN (
      SELECT DISTINCT at.channel_id
      FROM app.applied_tags at
      WHERE at.tag_id = ANY($3::text[])
    )
  );

-- name: ListChannelsFiltered :many
SELECT * FROM app.channels c
WHERE c.space_id = $1
  AND ($2::text = '' OR c.id ILIKE '%' || $2 || '%' OR c.name ILIKE '%' || $2 || '%')
  AND (
    cardinality($3::text[]) = 0
    OR c.id IN (
      SELECT DISTINCT at.channel_id
      FROM app.applied_tags at
      WHERE at.tag_id = ANY($3::text[])
    )
  )
ORDER BY c.updated_at DESC
LIMIT $4 OFFSET $5;

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

-- name: UpdateChannelConfig :one
UPDATE app.channels
SET config = $2,
    updated_by = $3,
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

-- name: DeleteChannelsBySpace :exec
DELETE FROM app.channels
WHERE space_id = $1;

-- name: ListPublishedChannels :many
SELECT * FROM app.channels
WHERE space_id = $1 AND published_at IS NOT NULL
ORDER BY published_at DESC;
