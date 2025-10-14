export enum PermissionNamespace {
  "spaces",
  "api-keys",
  "channels",
  "channel-executions",
  "analytics",
  "things",
  "models",
  "users",
  "data",
  "permissions",
  "policies",
  "certificates",
  "environments",
}

export enum PermissionAction {
  "*",
  "list",
  "create",
  "get",
  "update",
  "delete",
}

export const defaultOrgPermissions = {
  organizations: {
    get: false,
    update: false,
    delete: false,
  },
  'organization-api-keys': {
    list: false,
    create: false,
    delete: false,
  },
  spaces: {
    list: false,
    create: false,
    get: false,
    update: false,
    delete: false,
  },
  members: {
    list: false,
    create: false,
    get: false,
    update: false,
    delete: false,
  },
  permissions: {
    list: false,
    create: false,
    get: false,
    update: false,
    delete: false,
  },
}

export const defaultSpacePermissions = {
  'space-api-keys': {
    list: false,
    create: false,
    delete: false,
  },
  channels: {
    list: false,
    create: false,
    get: false,
    update: false,
    delete: false,
  },
  'channel-executions': {
    list: false,
    get: false,
  },
  // analytics: {
  //   list: false,
  //   create: false,
  //   get: false,
  //   update: false,
  //   delete: false,
  // },
  things: {
    list: false,
    create: false,
    get: false,
    update: false,
    delete: false,
  },
  models: {
    list: false,
    create: false,
    get: false,
    update: false,
    delete: false,
  },
  // data: {
  //   list: false,
  //   create: false,
  //   get: false,
  //   update: false,
  //   delete: false,
  // },
  policies: {
    list: false,
    get: false,
  },
  certificates: {
    list: false,
    get: false,
  },
  // environments: {
  //   list: false,
  //   create: false,
  //   get: false,
  //   update: false,
  //   delete: false,
  // },
}