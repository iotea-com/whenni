-- name: GetOrganization :one
SELECT * FROM app.organizations
WHERE id = $1;

-- name: ListOrganizations :many
SELECT * FROM app.organizations
ORDER BY name;

-- name: CreateOrganization :one
INSERT INTO app.organizations (
  id, name, created_by, updated_by
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateOrganization :one
UPDATE app.organizations
SET name = $2,
    updated_by = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteOrganization :exec
DELETE FROM app.organizations
WHERE id = $1;

-- name: DeleteOrganizationApiKeys :exec
DELETE FROM app."apiKeys"
WHERE organization_id = $1;

-- name: DeleteOrganizationPermissionSets :exec
DELETE FROM app.permissions
WHERE organization_id = $1;

-- name: ListOrganizationMembers :many
SELECT 
  om.*,
  u.name as user_name,
  u.email as user_email
FROM app.organization_members om
JOIN app.users u ON om.user_id = u.id
WHERE om.organization_id = $1
ORDER BY om.created_at DESC;

-- name: GetOrganizationMember :one
SELECT * FROM app.organization_members
WHERE organization_id = $1 AND user_id = $2;

-- name: CreateOrganizationMember :one
INSERT INTO app.organization_members (
  organization_id, user_id, role, organization_permission_set_id
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateOrganizationMemberRole :one
UPDATE app.organization_members
SET role = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateOrganizationMemberRoleByOrgAndUser :one
UPDATE app.organization_members
SET role = $3,
    updated_at = NOW()
WHERE organization_id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteOrganizationMember :exec
DELETE FROM app.organization_members
WHERE id = $1;

-- name: DeleteOrganizationMembers :exec
DELETE FROM app.organization_members
WHERE organization_id = $1;

-- name: DeleteOrganizationMemberByOrgAndUser :exec
DELETE FROM app.organization_members
WHERE organization_id = $1 AND user_id = $2;

-- name: CountOrganizationMembersWithFilter :one
SELECT COUNT(*)
FROM app.organization_members om
JOIN app.users u ON u.id = om.user_id
WHERE om.organization_id = $1
  AND ($2 = '' OR LOWER(u.email) LIKE LOWER('%' || $2 || '%'));

-- name: ListOrganizationMembersWithFilter :many
SELECT
  om.*,
  u.id as user_id,
  u.name as user_name,
  u.email as user_email
FROM app.organization_members om
JOIN app.users u ON u.id = om.user_id
WHERE om.organization_id = $1
  AND ($2 = '' OR LOWER(u.email) LIKE LOWER('%' || $2 || '%'))
ORDER BY om.created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountOrganizationAdmins :one
SELECT COUNT(*)
FROM app.organization_members
WHERE organization_id = $1 AND role = 'ADMIN';

-- name: GetDefaultOrganizationPermissionSet :one
SELECT *
FROM app.permissions
WHERE organization_id = $1
  AND space_id IS NULL
  AND name = 'Default'
LIMIT 1;
