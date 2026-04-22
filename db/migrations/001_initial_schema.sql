-- Set the schema
SET search_path TO app;

-- Create app schema if it doesn't exist
CREATE SCHEMA IF NOT EXISTS app;

-- Enums
CREATE TYPE app."AuthTokenType" AS ENUM (
  'EMAIL_VERIFICATION',
  'MAGIC_LINK',
  'REFRESH'
);

CREATE TYPE app."OrganizationRole" AS ENUM (
  'MEMBER',
  'ADMIN'
);

CREATE TYPE app."UserOnSpaceRole" AS ENUM (
  'MEMBER',
  'OWNER'
);

-- Organizations table
CREATE TABLE app.organizations (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_by TEXT NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_by TEXT NOT NULL
);

-- Users table
CREATE TABLE app.users (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  name TEXT,
  email TEXT UNIQUE,
  email_verified TIMESTAMPTZ,
  password TEXT,
  image TEXT
);

-- Accounts table (OAuth)
CREATE TABLE app.accounts (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  user_id TEXT NOT NULL REFERENCES app.users(id) ON DELETE CASCADE,
  type TEXT NOT NULL,
  provider TEXT NOT NULL,
  provider_account_id TEXT NOT NULL,
  refresh_token TEXT,
  access_token TEXT,
  expires_at INTEGER,
  token_type TEXT,
  scope TEXT,
  id_token TEXT,
  session_state TEXT,
  UNIQUE(provider, provider_account_id)
);

-- Sessions table
CREATE TABLE app.sessions (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  session_token TEXT NOT NULL UNIQUE,
  user_id TEXT NOT NULL REFERENCES app.users(id) ON DELETE CASCADE,
  expires TIMESTAMPTZ NOT NULL
);

-- Auth tokens table
CREATE TABLE app.auth_tokens (
  user_id TEXT NOT NULL REFERENCES app.users(id),
  token TEXT NOT NULL UNIQUE,
  type app."AuthTokenType" NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Spaces table
CREATE TABLE app.spaces (
  id TEXT PRIMARY KEY,
  organization_id TEXT NOT NULL REFERENCES app.organizations(id),
  name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_by TEXT NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_by TEXT NOT NULL
);

-- Permission sets table
CREATE TABLE app.permissions (
  id TEXT PRIMARY KEY,
  organization_id TEXT NOT NULL REFERENCES app.organizations(id),
  space_id TEXT REFERENCES app.spaces(id),
  name TEXT NOT NULL,
  permissions TEXT[] NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_by TEXT NOT NULL,
  updated_by TEXT NOT NULL
);

-- API Keys table
CREATE TABLE app."apiKeys" (
  id TEXT PRIMARY KEY,
  organization_id TEXT NOT NULL REFERENCES app.organizations(id),
  organization_permission_set_id TEXT REFERENCES app.permissions(id),
  space_id TEXT REFERENCES app.spaces(id),
  space_permission_set_id TEXT REFERENCES app.permissions(id),
  name TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_by TEXT NOT NULL,
  expires_at TIMESTAMPTZ
);

-- Organization members table
CREATE TABLE app.organization_members (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  organization_id TEXT NOT NULL REFERENCES app.organizations(id),
  user_id TEXT NOT NULL REFERENCES app.users(id),
  role app."OrganizationRole" NOT NULL DEFAULT 'MEMBER',
  organization_permission_set_id TEXT NOT NULL REFERENCES app.permissions(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(organization_id, user_id)
);

-- Member space permissions table
CREATE TABLE app.member_space_permissions (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  organization_member_id TEXT NOT NULL REFERENCES app.organization_members(id),
  permission_set_id TEXT NOT NULL REFERENCES app.permissions(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(organization_member_id, permission_set_id)
);

-- Channels table
CREATE TABLE app.channels (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  space_id TEXT NOT NULL REFERENCES app.spaces(id),
  config JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_by TEXT NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_by TEXT NOT NULL,
  published_at TIMESTAMPTZ,
  published_by TEXT
);

-- Models table
CREATE TABLE app.models (
  id TEXT PRIMARY KEY,
  space_id TEXT NOT NULL REFERENCES app.spaces(id),
  name TEXT NOT NULL,
  attributes JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_by TEXT NOT NULL,
  updated_by TEXT NOT NULL
);

-- Things table
CREATE TABLE app.things (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  space_id TEXT NOT NULL REFERENCES app.spaces(id),
  attributes JSONB NOT NULL,
  internal BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_by TEXT NOT NULL,
  updated_by TEXT NOT NULL,
  thing_category TEXT NOT NULL
);

-- Certificates table
CREATE TABLE app.certificates (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  space_id TEXT NOT NULL REFERENCES app.spaces(id),
  policy JSONB NOT NULL,
  revoke BOOLEAN NOT NULL DEFAULT false
);

-- Dev environments table
CREATE TABLE app.environments (
  id SERIAL PRIMARY KEY,
  environment_id TEXT NOT NULL UNIQUE,
  environment_name TEXT NOT NULL,
  region TEXT NOT NULL,
  state TEXT NOT NULL,
  public_dns TEXT NOT NULL,
  ssh_key_name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  space_id TEXT NOT NULL REFERENCES app.spaces(id)
);

-- Tags table
CREATE TABLE app.tags (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  name TEXT NOT NULL,
  space_id TEXT NOT NULL REFERENCES app.spaces(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_by TEXT NOT NULL,
  updated_by TEXT NOT NULL
);

-- Applied tags table
CREATE TABLE app.applied_tags (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  space_id TEXT NOT NULL REFERENCES app.spaces(id),
  tag_id TEXT NOT NULL REFERENCES app.tags(id),
  channel_id TEXT REFERENCES app.channels(id),
  model_id TEXT REFERENCES app.models(id),
  thing_id TEXT REFERENCES app.things(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_by TEXT NOT NULL,
  updated_by TEXT NOT NULL
);

-- Create indexes for foreign keys and common queries
CREATE INDEX idx_accounts_user_id ON app.accounts(user_id);
CREATE INDEX idx_sessions_user_id ON app.sessions(user_id);
CREATE INDEX idx_auth_tokens_user_id ON app.auth_tokens(user_id);
CREATE INDEX idx_spaces_organization_id ON app.spaces(organization_id);
CREATE INDEX idx_permissions_organization_id ON app.permissions(organization_id);
CREATE INDEX idx_permissions_space_id ON app.permissions(space_id);
CREATE INDEX idx_api_keys_organization_id ON app."apiKeys"(organization_id);
CREATE INDEX idx_api_keys_space_id ON app."apiKeys"(space_id);
CREATE INDEX idx_organization_members_organization_id ON app.organization_members(organization_id);
CREATE INDEX idx_organization_members_user_id ON app.organization_members(user_id);
CREATE INDEX idx_member_space_permissions_organization_member_id ON app.member_space_permissions(organization_member_id);
CREATE INDEX idx_channels_space_id ON app.channels(space_id);
CREATE INDEX idx_models_space_id ON app.models(space_id);
CREATE INDEX idx_things_space_id ON app.things(space_id);
CREATE INDEX idx_certificates_space_id ON app.certificates(space_id);
CREATE INDEX idx_environments_space_id ON app.environments(space_id);
CREATE INDEX idx_tags_space_id ON app.tags(space_id);
CREATE INDEX idx_applied_tags_space_id ON app.applied_tags(space_id);
CREATE INDEX idx_applied_tags_tag_id ON app.applied_tags(tag_id);
CREATE INDEX idx_applied_tags_channel_id ON app.applied_tags(channel_id);
CREATE INDEX idx_applied_tags_model_id ON app.applied_tags(model_id);
CREATE INDEX idx_applied_tags_thing_id ON app.applied_tags(thing_id);
