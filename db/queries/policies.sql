-- name: GetCertificateInSpace :one
SELECT *
FROM app.certificates
WHERE id = $1
  AND space_id = $2;

-- name: CountPolicies :one
SELECT COUNT(*)
FROM app.certificates
WHERE space_id = $1;

-- name: ListPoliciesPaginated :many
SELECT *
FROM app.certificates
WHERE space_id = $1
ORDER BY id
LIMIT $2 OFFSET $3;

-- name: UpdatePolicy :one
UPDATE app.certificates
SET policy = $3,
    revoke = $4,
    updated_at = NOW()
WHERE id = $1
  AND space_id = $2
RETURNING *;
